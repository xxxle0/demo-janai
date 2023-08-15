package git_wrapper

import (
	"strings"

	"github.com/samber/lo"
)

type Worktree struct {
	CommitSHA string
	Path      string
	IsMain    bool
}

func NewWorkTree(path string, commitSHA string) Worktree {
	return Worktree{
		CommitSHA: commitSHA,
		Path:      path,
	}
}

func (w Worktree) Remove() error {
	commandBuilder := commandBuilderFunc()
	commandBuilder.AddCommand("worktree")
	commandBuilder.AddArgs([]string{"remove", w.Path})
	_, err := commandBuilder.Exec()
	if err != nil {
		return err
	}
	return nil
}

func ListWorktree(path string) ([]Worktree, error) {
	commandBuilder := commandBuilderFunc()
	commandBuilder.AddCommand("worktree")
	commandBuilder.AddArg("list")
	commandBuilder.SetDir(path)
	output, err := commandBuilder.Exec()
	if err != nil {
		return nil, err
	}
	formatedOutput := strings.Split(output, "\n")
	worktrees := lo.FilterMap(formatedOutput, func(s string, i int) (Worktree, bool) {
		trimmed := strings.Trim(s, " ")
		if len(trimmed) == 0 {
			return Worktree{}, false
		}
		w := GenerateWorktree(trimmed)
		if i == 0 {
			w.IsMain = true
		}
		return w, true
	})
	return worktrees, nil
}

func GenerateWorktree(row string) Worktree {
	splitWord := []string{}
	word := ""
	for i := 0; i < len(row); i++ {
		ch := string(row[i])
		if ch == " " {
			if len(word) > 0 {
				splitWord = append(splitWord, word)
				word = ""
			}
			continue
		}
		word += ch
		if i == len(row)-1 {
			if len(word) > 0 {
				splitWord = append(splitWord, word)
				word = ""
			}
		}
	}
	return Worktree{
		Path:      splitWord[0],
		CommitSHA: splitWord[1],
	}
}
