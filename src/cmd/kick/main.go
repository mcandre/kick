// Package main implements a git sync utility.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mcandre/kick"
)

var flagDebug = flag.Bool("debug", false, "Enable additional logging")
var flagVersion = flag.Bool("version", false, "Show version banner")
var flagHelp = flag.Bool("help", false, "Show usage menu")

func usage() {
	program, err := os.Executable()

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Usage: %v [OPTIONS]\n", program)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	switch {
	case *flagVersion:
		fmt.Printf("%s", kick.Version)
		os.Exit(0)
	case *flagHelp:
		usage()
		os.Exit(0)
	}

	config := kick.NewConfig()
	config.Debug = *flagDebug

	if nonce, ok := os.LookupEnv(kick.NonceEnvironmentVariable); ok && nonce == "1" {
		config.Nonce = true
	}

	if pullAll, ok := os.LookupEnv(kick.PullAllEnvironmentVariable); ok && pullAll != "1" {
		config.PullAll = false
	}

	if pushAll, ok := os.LookupEnv(kick.PushAllEnvironmentVariable); ok && pushAll != "1" {
		config.PushAll = false
	}

	if syncTags, ok := os.LookupEnv(kick.SyncTagsEnvironmentVariable); ok && syncTags != "1" {
		config.SyncTags = false
	}

	if commitMessage, ok := os.LookupEnv(kick.CommitMessageEnvironmentVariable); ok {
		config.CommitMessage = commitMessage
	}

	if err := config.Kick(); err != nil {
		log.Fatal(err)
	}
}
