//go:build unix

package ssh

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/creack/pty"
	"github.com/paramientos/leap/internal/config"
	"golang.org/x/term"
)

func connectWithSystemSSHRecording(cmd *exec.Cmd, recording io.Writer) error {
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer ptmx.Close()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				fmt.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	multi := io.MultiWriter(os.Stdout, recording)
	go func() { io.Copy(ptmx, os.Stdin) }()
	io.Copy(multi, ptmx)
	return nil
}

func connectWithSystemSSHNormal(conn config.Connection) error {
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

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
