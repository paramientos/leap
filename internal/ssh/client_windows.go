//go:build windows

package ssh

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/paramientos/leap/internal/config"
)

func connectWithSystemSSHRecording(cmd *exec.Cmd, recording io.Writer) error {
	// Windows doesn't support PTY, so we can't capture output while staying interactive
	// Fall back to normal connection without recording
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
