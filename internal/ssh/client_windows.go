//go:build windows

package ssh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
)

func connectWithSystemSSH(conn config.Connection) error {
	args := []string{"-t"}
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

func setupWindowResizeHandler(fd int, session *ssh.Session) {
	// SIGWINCH doesn't exist on Windows
}
