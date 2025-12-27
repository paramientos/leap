package ssh

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

func Connect(conn config.Connection, record bool) error {
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

	// If password is provided, use Go SSH client
	if conn.Password != "" && conn.IdentityFile == "" {
		return connectWithPassword(conn, recordingFile)
	}

	// Fallback to system SSH
	return connectWithSystemSSH(conn, recordingFile)
}

func connectWithSystemSSH(conn config.Connection, recording io.Writer) error {
	args := buildSSHArgs(conn)
	cmd := exec.Command("ssh", args...)

	if recording != nil {
		return connectWithSystemSSHRecording(cmd, recording)
	}

	return connectWithSystemSSHNormal(conn)
}

func buildSSHArgs(conn config.Connection) []string {
	args := []string{}

	if conn.IdentityFile != "" {
		args = append(args, "-i", conn.IdentityFile)
	}
	args = append(args, "-p", fmt.Sprintf("%d", conn.Port))
	if conn.JumpHost != "" {
		args = append(args, "-J", conn.JumpHost)
	}
	target := fmt.Sprintf("%s@%s", conn.User, conn.Host)
	args = append(args, target)

	return args
}

func connectWithPassword(conn config.Connection, recording io.Writer) error {
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

	if recording != nil {
		session.Stdout = io.MultiWriter(os.Stdout, recording)
	} else {
		session.Stdout = os.Stdout
	}
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return err
	}

	return session.Wait()
}
