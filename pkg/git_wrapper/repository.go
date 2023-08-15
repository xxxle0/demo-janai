package git_wrapper

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/samber/lo"

	"github.com/guardrailsio/go-scan-helper/logger"
)

type Repository struct {
	Url             string     `json:"url"`
	Dest            string     `json:"dest"`
	Branch          Branch     `json:"branch"`
	BasicAuthHeader string     `json:"basic_auth_header"`
	Worktrees       []Worktree `json:"worktrees"`
}

// SetBasicAuthHeader implements IRepository
func (r *Repository) SetBasicAuthHeader(basicAuthHeader string) {
	r.BasicAuthHeader = basicAuthHeader
}

// GetDestination implements IRepository
func (r *Repository) GetDestination() string {
	return r.Dest
}

// Load implements IRepository
func (*Repository) Load(url string, dest string) *Repository {
	panic("unimplemented")
}

// SetPrivate implements IRepository
func (*Repository) SetPrivate(isPrivate bool) {
	panic("unimplemented")
}

// SetProtol implements IRepository
func (*Repository) SetProtol(protocol string) {
	panic("unimplemented")
}

type IRepository interface {
	Load(url string, dest string) *Repository
	SetPrivate(isPrivate bool)
	SetProtol(protocol string)
	Branches() ([]string, error)
	CheckoutBranch(branch string) (*Branch, error)
	CheckoutCommit(commit string) (*Commit, error)
	AddWorktree(path string, commitSHA string) (*Worktree, error)
	UpdateRemoteOrigin(remoteUrl string, logger logger.ILogger) error
	FlushWorktree() error
	Fetch() error
	Pull() error
	GetDestination() string
	RemoveRepository() error
	SetBasicAuthHeader(string)
	GetDiffContentBetweenCommits(commit, target string) (string, error)
}

func NewRepository(url string, dest string) *Repository {
	return &Repository{
		Url:  url,
		Dest: dest,
	}
}

func Load(dest string) (IRepository, error) {
	worktrees, err := ListWorktree(dest)
	if err != nil {
		return nil, err
	}
	return &Repository{
		Dest:      dest,
		Worktrees: worktrees,
	}, nil
}

func (r *Repository) SetBranch(branch Branch) {
	r.Branch = branch
}

func (r *Repository) Branches() ([]string, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(r.Dest)
	commandBuilder.AddCommand("branch")
	output, err := commandBuilder.Exec()
	formatedOutput := strings.Split(output, "\n")
	formatedOutput = lo.FilterMap(formatedOutput, func(s string, i int) (string, bool) {
		formated := strings.Trim(strings.Trim(s, " "), "* ")
		if len(formated) > 0 {
			return formated, true
		}
		return "", false
	})
	return formatedOutput, err
}

func (r *Repository) CheckoutBranch(branch string) (*Branch, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(r.Dest)
	commandBuilder.AddCommand("checkout")
	commandBuilder.AddArg(branch)
	_, err := commandBuilder.Exec()
	if err != nil {
		return nil, err
	}
	b := NewBranch(branch)
	return &b, nil
}

func (r *Repository) CheckoutCommit(commit string) (*Commit, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(r.Dest)
	commandBuilder.AddCommand("checkout")
	commandBuilder.AddArg(commit)
	_, err := commandBuilder.Exec()
	if err != nil {
		return nil, err
	}
	c := NewCommit(commit, r.Dest)
	return &c, nil
}

func (r *Repository) AddWorktree(path string, commitSHA string) (*Worktree, error) {
	for _, worktree := range r.Worktrees {
		if worktree.Path == path && strings.HasPrefix(commitSHA, worktree.CommitSHA) { // worktree is short sha, but scan request is a full sha, so it's okay to check just prefix
			return &worktree, nil
		}
	}
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(r.Dest)
	commandBuilder.AddCommand("worktree")
	commandBuilder.AddArgs([]string{"add", path})
	if commitSHA != "" {
		commandBuilder.AddArg(commitSHA)
	}
	_, err := commandBuilder.Exec()
	if err != nil {
		return nil, err
	}
	w := NewWorkTree(path, commitSHA)
	r.Worktrees = append(r.Worktrees, w)
	return &w, nil
}

func (r *Repository) Fetch() error {
	commandBuilder := commandBuilderFunc()
	addBasicAuthHeader(commandBuilder, r.BasicAuthHeader)
	commandBuilder.SetDir(r.Dest)
	commandBuilder.AddCommand("fetch")
	_, err := commandBuilder.Exec()
	return err
}

func (r *Repository) Pull() error {
	commandBuilder := commandBuilderFunc()
	addBasicAuthHeader(commandBuilder, r.BasicAuthHeader)
	commandBuilder.SetDir(r.Dest)
	commandBuilder.AddCommand("pull")
	commandBuilder.Exec()
	_, err := commandBuilder.Exec()
	return err
}

func (r *Repository) FlushWorktree() error {
	for _, w := range r.Worktrees {
		if w.IsMain {
			continue
		}
		commandBuilder := commandBuilderFunc()
		commandBuilder.SetDir(r.Dest)
		commandBuilder.AddCommand("worktree")
		commandBuilder.AddArgs([]string{"remove", w.Path})
		_, err := commandBuilder.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func addBasicAuthHeader(cmd ICommandBuilder, token string) {
	if token != "" {
		basicAuthHeaderArgs := []string{"-c", fmt.Sprintf("http.extraHeader=Authorization: Basic %s", token)}
		cmd.AddBaseCommandArgs(basicAuthHeaderArgs)
	}
}

func (r *Repository) UpdateRemoteOrigin(remoteUrl string, logger logger.ILogger) error {
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(r.Dest)
	commandBuilder.SetLogger(logger)
	commandBuilder.AddCommand("remote")
	commandBuilder.AddArgs([]string{"set-url", "origin", remoteUrl})
	_, err := commandBuilder.Exec()
	return err
}

func (r *Repository) RemoveWorktree(worktreeDest string) error {
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(r.Dest)
	commandBuilder.AddCommand("worktree")
	commandBuilder.AddArgs([]string{"remove", worktreeDest})
	_, err := commandBuilder.Exec()
	return err
}

func (r *Repository) RemoveRepository() error {
	for _, w := range r.Worktrees {
		if w.IsMain {
			continue
		}
		err := r.RemoveWorktree(w.Path)
		if err != nil {
			log.Println("remove worktree fail", err)
			continue
		}
	}
	err := os.RemoveAll(r.Dest)
	return err
}

func (r *Repository) GetDiffContentBetweenCommits(commit, target string) (string, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(r.Dest)
	if commit == target {
		commandBuilder.AddCommand("show")
		commandBuilder.AddArg(commit)
		output, err := commandBuilder.Exec()
		return output, err
	}
	commandBuilder.AddCommand("diff")
	commandBuilder.AddArg(fmt.Sprintf("%s..%s", target, commit))
	output, err := commandBuilder.Exec()
	return output, err
}
