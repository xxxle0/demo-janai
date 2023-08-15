package git_wrapper

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	worktreeScheduler           *WorktreeScheduler
	createWorktreeSchedulerOnce sync.Once
)

type WorktreeRequest struct {
	ResultChan     chan *Worktree
	ErrorChan      chan string
	RootDest       string
	CommitSHA      string
	FullSourcePath string
}

type WorktreeScheduler struct {
	worktreeRequestChan chan *WorktreeRequest
}

func NewWorktreeScheduler() *WorktreeScheduler {
	if worktreeScheduler == nil {
		createWorktreeSchedulerOnce.Do(func() {
			worktreeScheduler = &WorktreeScheduler{
				worktreeRequestChan: make(chan *WorktreeRequest),
			}
			go worktreeScheduler.WorktreeConsumer()
		})
	}
	return worktreeScheduler
}

func (w *WorktreeScheduler) CreateWorktree(ctx context.Context, fullSourcePath string, commitSHA string, rootDest string) (*Worktree, error) {
	resultChan := make(chan *Worktree)
	errorChan := make(chan string)
	createWorktreeRequest := WorktreeRequest{
		FullSourcePath: fullSourcePath,
		CommitSHA:      commitSHA,
		ResultChan:     resultChan,
		ErrorChan:      errorChan,
		RootDest:       rootDest,
	}
	w.worktreeRequestChan <- &createWorktreeRequest
	select {
	case worktree := <-resultChan:
		return worktree, nil
	case <-time.After(20 * time.Minute):
		return nil, errors.New("Create worktree timeout")
	case err := <-errorChan:
		return nil, errors.New(err)
	}
}

func (w *WorktreeScheduler) WorktreeConsumer() {
	for {
		worktreeRequest := <-w.worktreeRequestChan
		rootRepository, err := Load(worktreeRequest.RootDest)
		if err != nil {
			worktreeRequest.ErrorChan <- fmt.Sprintf("Load root worktree fail %s", err.Error())
			continue
		}
		worktree, err := rootRepository.AddWorktree(worktreeRequest.FullSourcePath, worktreeRequest.CommitSHA)
		if err != nil {
			worktreeRequest.ErrorChan <- fmt.Sprintf("Create worktree fail %s", err.Error())
			continue
		}
		worktreeRequest.ResultChan <- worktree
	}
}
