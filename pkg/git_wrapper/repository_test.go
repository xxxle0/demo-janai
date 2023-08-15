package git_wrapper

import (
	"errors"
	mock_git_wrapper "operarius/mock/pkg/git_wrapper"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository unit test", func() {
	var mockCtrl *gomock.Controller
	var mockCommandBuilder *mock_git_wrapper.MockICommandBuilder
	old := commandBuilderFunc
	var repository *Repository
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCommandBuilder = mock_git_wrapper.NewMockICommandBuilder(mockCtrl)
		commandBuilderFunc = func() ICommandBuilder {
			return mockCommandBuilder
		}
		repository = NewRepository("kai-repo", "./tmp/kai-test")
	})
	AfterEach(func() {
		defer func() { commandBuilderFunc = old }()
	})
	Context("Branches() ([]string, error)", func() {
		It("Should return list branch", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("branch")
			mockCommandBuilder.EXPECT().Exec().Return("master\ndevelop\n", nil).Times(1)
			branches, _ := repository.Branches()
			Expect(branches).To(Equal([]string{"master", "develop"}))
		})

		It("Should return error if commandBuilder.Exec return error", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("branch")
			mockCommandBuilder.EXPECT().Exec().Return("", errors.New("Exec Error")).Times(1)
			_, err := repository.Branches()
			Expect(err).To(Equal(errors.New("Exec Error")))
		})
	})

	Context("CheckoutBranch(branch string) (*Branch, error)", func() {
		It("Should trigger checkout branch command and return new branch", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("checkout")
			mockCommandBuilder.EXPECT().AddArg("master")
			mockCommandBuilder.EXPECT().Exec().Return("", nil).Times(1)
			branch, _ := repository.CheckoutBranch("master")
			Expect(branch).Should(Equal(&Branch{
				name: "master",
			}))
		})

		It("Should return error if commandBuilder.Exec return error", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("checkout")
			mockCommandBuilder.EXPECT().AddArg("master")
			mockCommandBuilder.EXPECT().Exec().Return("", errors.New("Exec Error")).AnyTimes()
			_, err := repository.CheckoutBranch("master")
			Expect(err).To(Equal(errors.New("Exec Error")))
		})
	})

	Context("CheckoutCommit(commit string) (*Commit, error)", func() {
		It("Should trigger checkout commit command and return new commit", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("checkout")
			mockCommandBuilder.EXPECT().AddArg("commitSHA")
			mockCommandBuilder.EXPECT().Exec().Return("", nil).AnyTimes()
			commit, _ := repository.CheckoutCommit("commitSHA")
			Expect(commit).Should(Equal(&Commit{
				hash: "commitSHA",
				dest: "./tmp/kai-test",
			}))
		})

		It("Should return error if commandBuilder.Exec return error", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("checkout")
			mockCommandBuilder.EXPECT().AddArg("commitSHA")
			mockCommandBuilder.EXPECT().Exec().Return("", errors.New("Exec Error")).AnyTimes()
			_, err := repository.CheckoutCommit("commitSHA")
			Expect(err).To(Equal(errors.New("Exec Error")))
		})
	})

	Context("Fetch()", func() {
		It("Should trigger git fetch command", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("fetch")
			mockCommandBuilder.EXPECT().Exec().Return("", nil).AnyTimes()
			repository.Fetch()
		})
	})

	Context("Pull()", func() {
		It("Should trigger git pull command", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("pull")
			mockCommandBuilder.EXPECT().Exec().Return("", nil).AnyTimes()
			repository.Pull()
		})
	})

	Context("AddWorktree(path string, commitSHA string) (*Worktree, error)", func() {
		It("Should trigger checkout commit command and return new commit", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("worktree")
			mockCommandBuilder.EXPECT().AddArgs([]string{"add", "./kai-clone-repo"})
			mockCommandBuilder.EXPECT().Exec().Return("", nil)
			worktree, _ := repository.AddWorktree("./kai-clone-repo", "")
			Expect(worktree).Should(Equal(&Worktree{
				Path:   "./kai-clone-repo",
				IsMain: false,
			}))
		})

		It("Should return error if commandBuilder.Exec return error", func() {
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("worktree")
			mockCommandBuilder.EXPECT().AddArgs([]string{"add", "./kai-clone-repo"})
			mockCommandBuilder.EXPECT().Exec().Return("", errors.New("Exec Error")).AnyTimes()
			_, err := repository.AddWorktree("./kai-clone-repo", "")
			Expect(err).To(Equal(errors.New("Exec Error")))
		})
	})

	Context("FlushWorktree() error ", func() {
		It("Should trigger multiple worktree remove command", func() {
			repository.Worktrees = []Worktree{
				{
					Path:   "./kai-clone-repo",
					IsMain: true,
				},
				{
					Path:   "./kai-clone-repo-2",
					IsMain: false,
				},
				{
					Path:   "./kai-clone-repo-3",
					IsMain: false,
				},
			}
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("worktree")
			mockCommandBuilder.EXPECT().AddArgs([]string{"remove", "./kai-clone-repo-2"})
			mockCommandBuilder.EXPECT().Exec().Return("", nil)
			mockCommandBuilder.EXPECT().SetDir("./tmp/kai-test")
			mockCommandBuilder.EXPECT().AddCommand("worktree")
			mockCommandBuilder.EXPECT().AddArgs([]string{"remove", "./kai-clone-repo-3"})
			mockCommandBuilder.EXPECT().Exec().Return("", nil)
			err := repository.FlushWorktree()
			Expect(err).Should(BeNil())
		})
	})
})
