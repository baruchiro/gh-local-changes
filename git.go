package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type Repo interface {
	getUnpushedBranches() (map[string]string, error)
	getUnpushedChanges() (int, error)
}

type GitRepo struct {
	folder string
}

var _ Repo = (*GitRepo)(nil)

func runGit(folder string, commandAndArgs ...string) (string, error) {
	cmd := exec.Command("git", commandAndArgs...)
	cmd.Dir = folder
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running git: %w", err)
	}
	return string(out), nil
}

func (r *GitRepo) getUnpushedBranches() (map[string]string, error) {
	out, err := runGit(r.folder, "log", "--branches", "--not", "--remotes", `--pretty=format:"%h %d"`)
	if err != nil {
		return nil, fmt.Errorf("error running git log: %w", err)
	}

	branches := make(map[string]string)
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		commit := words[0]
		branch := words[1]
		if !strings.HasPrefix(branch, "(") {
			continue
		}
		branch = strings.Trim(branch, "()")
		branches[branch] = commit
	}
	return branches, nil
}

func (r *GitRepo) getUnpushedChanges() (int, error) {
	out, err := runGit(r.folder, "status", "--porcelain")
	if err != nil {
		return -1, fmt.Errorf("error running git status: %w", err)
	}
	return strings.Count(out, "\n"), nil
}
