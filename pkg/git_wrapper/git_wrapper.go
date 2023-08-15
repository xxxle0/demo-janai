package git_wrapper

import (
	"fmt"
	"os"
)

var commandBuilderFunc = NewCommandBuilder
var listWorktreeFunc = ListWorktree

type CloneOptions struct {
	Username   string
	AuthToken  string
	Branch     string
	FilterSpec string
	Progress   bool
}

func Clone(url string, dest string, option *CloneOptions) (*Repository, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.AddCommand("clone")
	if option.FilterSpec != "" {
		commandBuilder.AddArg("--filter=" + option.FilterSpec)
	}
	if option.Branch != "" {
		commandBuilder.AddArg("--branch=" + option.Branch)
	}
	if option.AuthToken != "" {
		urlWithAuthToken := fmt.Sprintf("https://%s:%s@%s", option.Username, option.AuthToken, url)
		commandBuilder.AddArg(urlWithAuthToken)
	}
	if dest != "" {
		commandBuilder.AddArg(dest)
	}
	if option.Progress {
		_, err := commandBuilder.ExecStreamStdout()
		if err != nil {
			return nil, err
		}
	} else {
		_, err := commandBuilder.Exec()
		if err != nil {
			return nil, err
		}
	}
	worktrees, _ := listWorktreeFunc(dest)
	return &Repository{
		Url:       url,
		Dest:      dest,
		Worktrees: worktrees,
	}, nil
}

func PlainClone(url string, dest string) (*Repository, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.AddCommand("clone")
	commandBuilder.AddArg(url)
	if dest != "" {
		commandBuilder.AddArg(dest)
	}
	_, err := commandBuilder.Exec()
	if err != nil {
		return nil, err
	}
	worktrees, _ := listWorktreeFunc(dest)
	return &Repository{
		Url:       url,
		Dest:      dest,
		Worktrees: worktrees,
	}, nil
}

func RemoveRepository(dest string) error {
	commandBuilder := commandBuilderFunc()
	commandBuilder.AddCommand("worktree")
	commandBuilder.AddArg("prune")
	if dest != "" {
		commandBuilder.AddArg(dest)
	}
	_, err := commandBuilder.Exec()
	if err != nil {
		return err
	}
	os.RemoveAll(dest)
	return nil
}
