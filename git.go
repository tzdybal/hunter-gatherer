package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func git() {
	err := pushToGit()
	if err != nil {
		fmt.Println(err)
	}
}

func pushToGit() error {
	stat, err := os.Stat(filepath.Join(archiveDir, ".git"))
	if err != nil {
		if os.IsNotExist(err) {
			err := runGitCommand("init")
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !stat.IsDir() {
		return fmt.Errorf("'%s' is not a directory", stat.Name())
	}

	err = runGitCommand("add", ".")
	if err != nil {
		return err
	}
	err = runGitCommand("commit", "-m", "hunter-gatherer's auto commit")
	if err != nil {
		return err
	}
	return runGitCommand("push", "origin")
}

func runGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = archiveDir
	return cmd.Run()
}
