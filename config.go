// Package kick provides high level git sync operations.
package kick

import (
	"bufio"
	"bytes"
	"io"
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

// FetchAllEnvironmentVariable denotes the name of the environment variable controlling whether fetches process all remotes.
const FetchAllEnvironmentVariable = "KICK_FETCH_ALL"

// PullAllEnvironmentVariable denotes the name of the environment variable controlling whether pulls process all remotes.
const PullAllEnvironmentVariable = "KICK_PULL_ALL"

// PushAllEnvironmentVariable denotes the name of the environment variable controlling whether pushes process all remotes.
const PushAllEnvironmentVariable = "KICK_PUSH_ALL"

// SyncTagsEnvironmentVariable denotes the name of the environment variable controlling whether to push and pull tags.
const SyncTagsEnvironmentVariable = "KICK_SYNC_TAGS"

// Config prepares high level git sync operations.
type Config struct {
	// Debug enables additional logging (default: false).
	Debug bool

	// Nonce enables altering NoncePath to generate commits when repositories are otherwise unchanged (default: false).
	Nonce bool

	// FetchAll enables fetching from all remotes (default: true).
	FetchAll bool

	// PullAll enables pulling from all remotes (default: true).
	PullAll bool

	// PushAll enables pushing to all remotes (default: true).
	PushAll bool

	// SyncTags enables pushing and pulling tags (default: true).
	SyncTags bool

	// CommitMessage denotes a git commit message (default: DefaultCommitMessage).
	CommitMessage string

	// remotes tracks the repository's configured remote names.
	remotes []string
}

// NewConfig constructs a Config.
func NewConfig() Config {
	return Config{
		FetchAll:      true,
		PullAll:       true,
		PushAll:       true,
		SyncTags:      true,
		CommitMessage: DefaultCommitMessage,
	}
}

// QueryRemotes populates metadata for remotes.
func (o *Config) QueryRemotes() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "remote")
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	remotesBytes, err := cmd.Output()

	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(remotesBytes))
	o.remotes = o.remotes[:0]

	for scanner.Scan() {
		line := scanner.Text()
		o.remotes = append(o.remotes, line)
	}

	return nil
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

	if o.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

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

	if o.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}

// Pull pulls any remote changes.
func (o Config) Pull() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "pull")

	if o.PullAll {
		cmd.Args = append(cmd.Args, "--all")
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	if o.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}

// Push pushes any local changes.
func (o Config) Push() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "push")

	if o.PushAll {
		cmd.Args = append(cmd.Args, "--all")
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	if o.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}

// FetchTags fetches any remote tags.
func (o Config) FetchTags() error {
	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "fetch", "--tags")

	if o.FetchAll {
		cmd.Args = append(cmd.Args, "--all")
	}

	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	if o.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

	if o.Debug {
		log.Printf("cmd: %v\n", cmd)
	}

	return cmd.Run()
}

// PushTags pushes any local tags.
func (o Config) PushTags() error {
	if o.PushAll {
		for _, remote := range o.remotes {
			cmd := exec.Command("git")
			cmd.Args = append(cmd.Args, "push", remote, "--tags")
			cmd.Env = os.Environ()
			cmd.Stdin = os.Stdin

			if o.Debug {
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			} else {
				cmd.Stdout = io.Discard
				cmd.Stderr = io.Discard
			}

			if o.Debug {
				log.Printf("cmd: %v\n", cmd)
			}

			if err := cmd.Run(); err != nil {
				return err
			}
		}

		return nil
	}

	cmd := exec.Command("git")
	cmd.Args = append(cmd.Args, "push", "--tags")
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin

	if o.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}

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
// * Pulling and pushing tags
func (o Config) Kick() error {
	if o.Debug {
		log.Printf("config: %v\n", o)
	}

	if err := o.QueryRemotes(); err != nil {
		return err
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

	if err := o.Push(); err != nil {
		return err
	}

	if o.SyncTags {
		if err := o.FetchTags(); err != nil {
			return err
		}

		if err := o.PushTags(); err != nil {
			return err
		}
	}

	return nil
}
