// Package kick provides high level git sync operations.
package kick

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Config prepares high level git sync operations.
type Config struct {
	// Debug enables additional logging (default: false).
	Debug bool
}

// Kick automates:
//
// * Staging all file changes
// * Committing all changes
// * Pulling any remote changes
// * Pushing any local changes
func (o Config) Kick() error {
	if o.Debug {
		log.Printf("config: %v\n", o)
	}

	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "add", ".")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git")
	cmd.Args = append(cmd.Args, "commit", "-am", "up")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	cmd = exec.Command("git")
	cmd.Args = append(cmd.Args, "pull")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git")
	cmd.Args = append(cmd.Args, "push")
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}
