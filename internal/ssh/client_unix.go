//go:build !windows

package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/paramientos/leap/internal/config"
	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

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

func setupWindowResizeHandler(fd int, session *ssh.Session) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)
	go func() {
		for range sig {
			nw, nh, _ := term.GetSize(fd)
			session.WindowChange(nh, nw)
		}
	}()
}
