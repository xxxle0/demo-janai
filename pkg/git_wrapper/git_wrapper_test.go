package git_wrapper

import (
	"fmt"
	mock_git_wrapper "operarius/mock/pkg/git_wrapper"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Git wrapper unit test", func() {
	var mockCtrl *gomock.Controller
	var mockCommandBuilder *mock_git_wrapper.MockICommandBuilder
	old := commandBuilderFunc
	oldListWorktree := ListWorktree
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCommandBuilder = mock_git_wrapper.NewMockICommandBuilder(mockCtrl)
		commandBuilderFunc = func() ICommandBuilder {
			return mockCommandBuilder
		}
		listWorktreeFunc = func(path string) ([]Worktree, error) {
			return []Worktree{}, nil
		}
	})
	AfterEach(func() {
		defer func() { commandBuilderFunc = old }()
		defer func() { listWorktreeFunc = oldListWorktree }()
	})
	Context("Clone(url string, dest string, option *CloneOptions) (*Repository, error)", func() {
		It("Should call the clone command with auth token in clone url", func() {
			urlWithAuthToken := fmt.Sprintf("https://guardrails:%s@%s", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c", "git@github.com:guardrailsio/core-api.git")
			cloneOption := &CloneOptions{
				AuthToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				Username:  "guardrails",
			}
			mockCommandBuilder.EXPECT().AddCommand("clone")
			mockCommandBuilder.EXPECT().AddArg(urlWithAuthToken)
			mockCommandBuilder.EXPECT().Exec().Return("", nil)
			repository, _ := Clone("git@github.com:guardrailsio/core-api.git", "", cloneOption)
			Expect(repository).To(Equal(&Repository{
				Url:       "git@github.com:guardrailsio/core-api.git",
				Dest:      "",
				Worktrees: []Worktree{},
			}))
		})
		It("Should call the clone command with blobless filter", func() {
			urlWithAuthToken := fmt.Sprintf("https://guardrails:%s@%s", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ", "git@github.com:guardrailsio/core-api.git")
			cloneOption := &CloneOptions{
				AuthToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ",
				FilterSpec: "blobless",
				Username:   "guardrails",
			}
			mockCommandBuilder.EXPECT().AddCommand("clone")
			mockCommandBuilder.EXPECT().AddArg(urlWithAuthToken)
			mockCommandBuilder.EXPECT().AddArg("--filter=blobless")
			mockCommandBuilder.EXPECT().Exec()
			repository, _ := Clone("git@github.com:guardrailsio/core-api.git", "", cloneOption)
			Expect(repository).To(Equal(&Repository{
				Url:       "git@github.com:guardrailsio/core-api.git",
				Dest:      "",
				Worktrees: []Worktree{},
			}))
		})
		It("Should call the clone command with treeless filter", func() {
			urlWithAuthToken := fmt.Sprintf("https://guardrails:%s@%s", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ", "git@github.com:guardrailsio/core-api.git")
			cloneOption := &CloneOptions{
				AuthToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ",
				FilterSpec: "treeless",
				Username:   "guardrails",
			}
			mockCommandBuilder.EXPECT().AddCommand("clone")
			mockCommandBuilder.EXPECT().AddArg(urlWithAuthToken)
			mockCommandBuilder.EXPECT().AddArg("--filter=treeless")
			mockCommandBuilder.EXPECT().Exec()
			repository, _ := Clone("git@github.com:guardrailsio/core-api.git", "", cloneOption)
			Expect(repository).To(Equal(&Repository{
				Url:       "git@github.com:guardrailsio/core-api.git",
				Dest:      "",
				Worktrees: []Worktree{},
			}))
		})
		It("Should call the clone command with dest", func() {
			urlWithAuthToken := fmt.Sprintf("https://guardrails:%s@%s", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ", "git@github.com:guardrailsio/core-api.git")
			cloneOption := &CloneOptions{
				AuthToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.NHVaYe26MbtOYhSKkoKYdFVomg4i8ZJd8_-RU8VNbftc4TSMb4bXP3l3YlNWACwyXPGffz5aXHc6lty1Y2t4SWRqGteragsVdZufDn5BlnJl9pdR_kdVFUsra2rWKEofkZeIC4yWytE58sMIihvo9H1ScmmVwBcQP6XETqYd0aSHp1gOa9RdUPDvoXQ5oqygTqVtxaDr6wUFKrKItgBMzWIdNZ6y7O9E0DhEPTbE9rfBo6KTFsHAZnMg4k68CDp2woYIaXbmYTWcvbzIuHO7_37GT79XdIwkm95QJ7hYC9RiwrV7mesbY4PAahERJawntho0my942XheVLmGwLMBkQ",
				Username:  "guardrails",
			}
			mockCommandBuilder.EXPECT().AddCommand("clone")
			mockCommandBuilder.EXPECT().AddArg(urlWithAuthToken)
			mockCommandBuilder.EXPECT().AddArg("./kaka")
			mockCommandBuilder.EXPECT().Exec()
			repository, _ := Clone("git@github.com:guardrailsio/core-api.git", "./kaka", cloneOption)
			Expect(repository).To(Equal(&Repository{
				Url:       "git@github.com:guardrailsio/core-api.git",
				Dest:      "./kaka",
				Worktrees: []Worktree{},
			}))
		})
	})
	Context("PlainClone(url string, dest string) (*Repository, error)", func() {
		It("Should call plan clone with provided url", func() {
			mockCommandBuilder.EXPECT().AddCommand("clone")
			mockCommandBuilder.EXPECT().AddArg("git@github.com:guardrailsio/core-api.git")
			mockCommandBuilder.EXPECT().Exec()
			repository, _ := PlainClone("git@github.com:guardrailsio/core-api.git", "")
			Expect(repository).To(Equal(&Repository{
				Url:       "git@github.com:guardrailsio/core-api.git",
				Worktrees: []Worktree{},
			}))
		})
		It("Should call plan clone with provided url and dest", func() {
			mockCommandBuilder.EXPECT().AddCommand("clone")
			mockCommandBuilder.EXPECT().AddArg("git@github.com:guardrailsio/core-api.git")
			mockCommandBuilder.EXPECT().AddArg("./tmp")
			mockCommandBuilder.EXPECT().Exec()
			repository, _ := PlainClone("git@github.com:guardrailsio/core-api.git", "./tmp")
			Expect(repository).To(Equal(&Repository{
				Url:       "git@github.com:guardrailsio/core-api.git",
				Dest:      "./tmp",
				Worktrees: []Worktree{},
			}))
		})
	})
})
