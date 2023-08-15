package git_wrapper

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type Commit struct {
	hash string
	dest string
}

func NewCommit(commitHash string, dest string) Commit {
	return Commit{
		hash: commitHash,
		dest: dest,
	}
}

func (c Commit) DiffListFileChanged(targetCommit *Commit) ([]string, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.SetDir(c.dest)
	commandBuilder.AddCommand("diff")
	// only file added, changed, modified
	commandBuilder.AddArgs([]string{"--name-only", "--diff-filter=ACMR"})
	if targetCommit != nil && targetCommit.hash != "" {
		commandBuilder.AddArg(targetCommit.hash)
	}
	commandBuilder.AddArg(c.hash)
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

func (c Commit) DiffShortStat(targetCommit *Commit) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.AddCommand("diff")
	commandBuilder.AddArgs([]string{"--shortstat", c.hash})
	if targetCommit != nil && targetCommit.hash != "" {
		commandBuilder.AddArg(targetCommit.hash)
	}
	output, err := commandBuilder.Exec()
	fmt.Print(string(output), err)
}
