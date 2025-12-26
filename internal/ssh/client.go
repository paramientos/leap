package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Connect(conn config.Connection) error {
	// If password is provided, use Go SSH client
	if conn.Password != "" && conn.IdentityFile == "" {
		return connectWithPassword(conn)
	}

	// Fallback to system SSH
	return connectWithSystemSSH(conn)
}

func connectWithSystemSSH(conn config.Connection) error {
	args := []string{}

	// Add Identity File
	if conn.IdentityFile != "" {
		args = append(args, "-i", conn.IdentityFile)
	}

	// Add Port
	args = append(args, "-p", fmt.Sprintf("%d", conn.Port))

	// Add Jump Host
	if conn.JumpHost != "" {
		args = append(args, "-J", conn.JumpHost)
	}

	// Targethost
	target := fmt.Sprintf("%s@%s", conn.User, conn.Host)
	args = append(args, target)

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func connectWithPassword(conn config.Connection) error {
	sshConfig := &ssh.ClientConfig{
		User: conn.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(conn.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", conn.Host, conn.Port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	fileDescriptor := int(os.Stdin.Fd())
	if term.IsTerminal(fileDescriptor) {
		originalState, err := term.MakeRaw(fileDescriptor)
		if err != nil {
			return err
		}
		defer term.Restore(fileDescriptor, originalState)

		width, height, err := term.GetSize(fileDescriptor)
		if err != nil {
			return err
		}

		if err := session.RequestPty("xterm-256color", height, width, modes); err != nil {
			return err
		}
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return err
	}

	return session.Wait()
}
