package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/term"
)

func Connect(conn config.Connection, record bool) error {
	// If password exists, we MUST use native to auto-fill it
	// If it's a key-only connection, syscall.Exec is better (more stable)
	if conn.Password != "" || record {
		return connectNative(conn, record)
	}

	return connectWithSystemSSH(conn)
}

func connectWithSystemSSH(conn config.Connection) error {
	binary, err := exec.LookPath("ssh")
	if err != nil {
		binary = "/usr/bin/ssh"
	}

	args := []string{"ssh", "-t"}
	if conn.IdentityFile != "" {
		args = append(args, "-i", conn.IdentityFile)
	}
	args = append(args, "-p", fmt.Sprintf("%d", conn.Port))
	if conn.JumpHost != "" {
		args = append(args, "-J", conn.JumpHost)
	}
	target := fmt.Sprintf("%s@%s", conn.User, conn.Host)
	args = append(args, target)

	return syscall.Exec(binary, args, os.Environ())
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

		// Handle window resize signals
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGWINCH)
		defer signal.Stop(sig)
		go func() {
			for range sig {
				nw, nh, _ := term.GetSize(fd)
				session.WindowChange(nh, nw)
			}
		}()
	}

	// PIPING - Donma sorununu çözen asıl kısım
	stdin, _ := session.StdinPipe()
	stdout, _ := session.StdoutPipe()
	stderr, _ := session.StderrPipe()

	go io.Copy(stdin, os.Stdin)

	var logout io.Writer = os.Stdout
	if recordingFile != nil {
		logout = io.MultiWriter(os.Stdout, recordingFile)
	}

	// Goroutines for output
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
