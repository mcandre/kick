//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/mcandre/kick"
	mageextras "github.com/mcandre/mage-extras"
)

// artifactsPath describes where artifacts are produced.
var artifactsPath = "bin"

// Default references the default build task.
var Default = Test

// Govulncheck runs govulncheck.
func Govulncheck() error { return mageextras.Govulncheck("-scan", "package", "./...") }

// Snyk runs Snyk SCA.
func Snyk() error { return mageextras.SnykTest() }

// Audit runs a security audit.
func Audit() error {
	mg.Deps(Govulncheck)
	return Snyk()
}

// Test executes a test suite.
func Test() error { return IntegrationTest() }

// IntegrationTest executes kick operations.
func IntegrationTest() error {
	cmd := exec.Command("kick")
	cmd.Args = append(cmd.Args, "-help")
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Deadcode runs deadcode.
func Deadcode() error { return mageextras.Deadcode("./...") }

// Gofmt runs gofmt.
func GoFmt() error { return mageextras.GoFmt("-s", "-w") }

// GoImports runs goimports.
func GoImports() error { return mageextras.GoImports("-w") }

// GoVet runs default go vet analyzers.
func GoVet() error { return mageextras.GoVet() }

// Errcheck runs errcheck.
func Errcheck() error { return mageextras.Errcheck("-blank") }

// Nakedret runs nakedret.
func Nakedret() error { return mageextras.Nakedret("-l", "0") }

// Revive runs revive.
func Revive() error { return mageextras.Revive("-set_exit_status") }

// Shadow runs go vet with shadow checks enabled.
func Shadow() error { return mageextras.GoVetShadow() }

// Staticcheck runs staticcheck.
func Staticcheck() error { return mageextras.Staticcheck("./...") }

// Unmake runs unmake.
func Unmake() error {
	err := mageextras.Unmake(".")

	if err != nil {
		return err
	}

	return mageextras.Unmake("-n", ".")
}

// Lint runs the lint suite.
func Lint() error {
	mg.Deps(Deadcode)
	mg.Deps(GoFmt)
	mg.Deps(GoImports)
	mg.Deps(GoVet)
	mg.Deps(Errcheck)
	mg.Deps(Nakedret)
	mg.Deps(Revive)
	mg.Deps(Shadow)
	mg.Deps(Staticcheck)
	mg.Deps(Unmake)
	return nil
}

// portBasename labels the artifact basename.
var portBasename = fmt.Sprintf("kick-%s", kick.Version)

// repoNamespace identifies the Go namespace for this project.
var repoNamespace = "github.com/mcandre/kick"

// Factorio cross-compiles Go binaries for a multitude of platforms.
func Factorio() error { return mageextras.Factorio(portBasename) }

// Port builds and compresses artifacts.
func Port() error { mg.Deps(Factorio); return mageextras.Archive(portBasename, artifactsPath) }

// Install builds and installs Go applications.
func Install() error { return mageextras.Install() }

// Uninstall deletes installed Go applications.
func Uninstall() error { return mageextras.Uninstall("kick") }

// Clean deletes artifacts.
func Clean() error { return os.RemoveAll(artifactsPath) }
