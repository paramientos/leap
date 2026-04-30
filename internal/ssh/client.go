package ssh

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/term"
)

func Connect(conn config.Connection, record bool) error {
	// If password exists, we MUST use native to auto-fill it
	// If it's a key-only connection, system SSH via syscall.Exec (on Unix) is better
	if conn.Password != "" || record {
		return connectNative(conn, record)
	}

	return connectWithSystemSSH(conn)
}

// knownHostsPath returns the path to the user's known_hosts file.
func knownHostsPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssh", "known_hosts")
}

// hostKeyCallback returns an ssh.HostKeyCallback that verifies against
// ~/.ssh/known_hosts and prompts the user to trust new hosts (like OpenSSH).
func hostKeyCallback() ssh.HostKeyCallback {
	khPath := knownHostsPath()

	// Ensure the file exists so knownhosts.New doesn't fail.
	if _, err := os.Stat(khPath); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(filepath.Dir(khPath), 0700)
		os.WriteFile(khPath, nil, 0600)
	}

	checker, err := knownhosts.New(khPath)
	if err != nil {
		// Fallback: if we can't read known_hosts, refuse all connections rather
		// than silently ignoring host keys.
		return func(_ string, _ net.Addr, _ ssh.PublicKey) error {
			return fmt.Errorf("failed to load known_hosts: %v", err)
		}
	}

	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		err := checker(hostname, remote, key)
		if err == nil {
			// Host is already known and key matches.
			return nil
		}

		// knownhosts returns a *knownhosts.KeyError when the host is unknown
		// (Wants is empty) or when the key has changed (Wants is non-empty).
		var keyErr *knownhosts.KeyError
		if !errors.As(err, &keyErr) {
			return err
		}

		if len(keyErr.Want) > 0 {
			// Key mismatch — possible MITM, refuse immediately.
			fmt.Fprintf(os.Stderr, "\n\033[31m@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
			fmt.Fprintf(os.Stderr, "@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @\n")
			fmt.Fprintf(os.Stderr, "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\033[0m\n")
			fmt.Fprintf(os.Stderr, "IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!\n")
			fmt.Fprintf(os.Stderr, "The host key for %s has changed.\n", hostname)
			fmt.Fprintf(os.Stderr, "Remove the old key from %s if you are sure this is legitimate.\n\n", khPath)
			return err
		}

		// Host not yet in known_hosts — ask the user.
		fingerprint := ssh.FingerprintSHA256(key)
		fmt.Fprintf(os.Stderr, "\nThe authenticity of host '%s' can't be established.\n", hostname)
		fmt.Fprintf(os.Stderr, "%s key fingerprint is %s.\n", key.Type(), fingerprint)
		fmt.Fprintf(os.Stderr, "Are you sure you want to continue connecting (yes/no)? ")

		reader := bufio.NewReader(os.Stdin)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "yes" && answer != "y" {
			return fmt.Errorf("host key verification failed: user declined")
		}

		// Append the new host key to known_hosts.
		f, err := os.OpenFile(khPath, os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("could not update known_hosts: %v", err)
		}
		defer f.Close()

		line := knownhosts.Line([]string{knownhosts.Normalize(hostname)}, key)
		if _, err := fmt.Fprintln(f, line); err != nil {
			return fmt.Errorf("could not write to known_hosts: %v", err)
		}

		fmt.Fprintf(os.Stderr, "Warning: Permanently added '%s' (%s) to the list of known hosts.\n", hostname, key.Type())
		return nil
	}
}

func connectNative(conn config.Connection, record bool) error {
	var recordingFile *os.File
	if record {
		home, _ := os.UserHomeDir()
		historyDir := filepath.Join(home, ".leap", "history")
		os.MkdirAll(historyDir, 0700)
		fileName := fmt.Sprintf("%s_%s.cast", conn.Name, time.Now().Format("20060102_150405"))
		path := filepath.Join(historyDir, fileName)
		f, err := os.Create(path)
		if err == nil {
			recordingFile = f
			defer f.Close()
			fmt.Printf("\n⏺️  \033[90mRecording session to %s\033[0m\n", path)
		}
	}

	sshConfig := &ssh.ClientConfig{
		User:            conn.User,
		HostKeyCallback: hostKeyCallback(),
		Timeout:         15 * time.Second,
	}

	var auth []ssh.AuthMethod
	if conn.IdentityFile != "" {
		if key, err := os.ReadFile(conn.IdentityFile); err == nil {
			if signer, err := ssh.ParsePrivateKey(key); err == nil {
				auth = append(auth, ssh.PublicKeys(signer))
			}
		}
	}
	if conn.Password != "" {
		auth = append(auth, ssh.Password(conn.Password))
	}
	if socket := os.Getenv("SSH_AUTH_SOCK"); socket != "" {
		if netConn, err := net.Dial("unix", socket); err == nil {
			agentClient := agent.NewClient(netConn)
			auth = append(auth, ssh.PublicKeysCallback(agentClient.Signers))
		}
	}
	sshConfig.Auth = auth

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", conn.Host, conn.Port), sshConfig)
	if err != nil {
		return fmt.Errorf("dial failed: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("session connection failed: %v", err)
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			return err
		}
		defer term.Restore(fd, oldState)

		w, h, _ := term.GetSize(fd)
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		if err := session.RequestPty("xterm", h, w, modes); err != nil {
			return err
		}

		// Handle window resize signals (Platform specific)
		setupWindowResizeHandler(fd, session)
	}

	// PIPING
	stdin, _ := session.StdinPipe()
	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()

	go io.Copy(stdin, os.Stdin)

	var logout io.Writer = os.Stdout
	if recordingFile != nil {
		logout = io.MultiWriter(os.Stdout, recordingFile)
	}

	go io.Copy(logout, stdout)
	go io.Copy(os.Stderr, stderr)

	if err := session.Shell(); err != nil {
		return err
	}

	return session.Wait()
}

func RunCommand(conn config.Connection, command string) (string, error) {
	sshConfig := &ssh.ClientConfig{
		User:            conn.User,
		HostKeyCallback: hostKeyCallback(),
		Timeout:         10 * time.Second,
	}
	var auth []ssh.AuthMethod
	if conn.IdentityFile != "" {
		if key, err := os.ReadFile(conn.IdentityFile); err == nil {
			if signer, err := ssh.ParsePrivateKey(key); err == nil {
				auth = append(auth, ssh.PublicKeys(signer))
			}
		}
	}
	if conn.Password != "" {
		auth = append(auth, ssh.Password(conn.Password))
	}
	sshConfig.Auth = auth

	addr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return "", err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}
