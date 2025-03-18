// Package kick provides high level git sync operations.
package kick

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// DefaultCommitMessage denotes a standard commit message.
const DefaultCommitMessage = "up"

// CommitMessageEnvironmentVariable denotes the name of the environment variable controlling commit messages.
const CommitMessageEnvironmentVariable = "KICK_MESSAGE"

// NoncePath denotes the file path to a nonce file, relative to the current working directory.
const NoncePath = ".kick"

// NonceEnvironmentVariable denotes the name of the environment variable controlling nonces.
const NonceEnvironmentVariable = "KICK_NONCE"

// Config prepares high level git sync operations.
type Config struct {
	// Debug enables additional logging (default: false).
	Debug bool

	// Nonce enables altering NoncePath to generate commits when repositories are otherwise unchanged (default: false).
	Nonce bool

	// CommitMessage denotes a git commit message (default: DefaultCommitMessage).
	CommitMessage string
}

// NewConfig constructs a Config.
func NewConfig() Config {
	return Config{
		CommitMessage: DefaultCommitMessage,
	}
}

// EnsureNonce updates NoncePath with the current timestamp.
func (o Config) EnsureNonce() error {
	tRFC3339Bytes := []byte(time.Now().UTC().Format(time.RFC3339))
	return os.WriteFile(NoncePath, tRFC3339Bytes, 0644)
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

	if o.Nonce {
		if err := o.EnsureNonce(); err != nil {
			return err
		}
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
	cmd.Args = append(cmd.Args, "commit", "-am", o.CommitMessage)
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
