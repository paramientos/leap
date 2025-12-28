package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
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
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
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
