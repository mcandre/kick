// Package kick provides high level git sync operations.
package kick

import (
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

// PullAllEnvironmentVariable denotes the name of the environment variable controlling whether pulls process all remotes.
const PullAllEnvironmentVariable = "KICK_PULL_ALL"

// PushAllEnvironmentVariable denotes the name of the environment variable controlling whether pushes process all remotes.
const PushAllEnvironmentVariable = "KICK_PUSH_ALL"

// Config prepares high level git sync operations.
type Config struct {
	// Debug enables additional logging (default: false).
	Debug bool

	// Nonce enables altering NoncePath to generate commits when repositories are otherwise unchanged (default: false).
	Nonce bool

	// PullAll enables pulling from all remotes (default: false).
	PullAll bool

	// PushAll enables pushing to all remotes (default: false).
	PushAll bool

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

// Stage stages any local file changes.
func (o Config) Stage() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "add", ".")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}

// Commit commits any staged changes.
func (o Config) Commit() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "commit", "-a")

	if o.CommitMessage != "" {
		cmd.Args = append(cmd.Args, "-m", o.CommitMessage)
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}

// Pull pulls any remote changes from the current branch's default remote.
func (o Config) Pull() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "pull")

	if o.PullAll {
		cmd.Args = append(cmd.Args, "--all")
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}

// Push pushes any local changes to the current branch's default remote.
func (o Config) Push() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "push")

	if o.PushAll {
		cmd.Args = append(cmd.Args, "--all")
	}

	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
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

	if err := o.Stage(); err != nil {
		return err
	}

	if err := o.Commit(); err != nil {
		if o.Debug {
			log.Println(err)
		}
	}

	if err := o.Pull(); err != nil {
		return err
	}

	return o.Push()
}
