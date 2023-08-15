package git_wrapper

import (
	mock_git_wrapper "operarius/mock/pkg/git_wrapper"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worktree unit test", func() {
	var mockCtrl *gomock.Controller
	var mockCommandBuilder *mock_git_wrapper.MockICommandBuilder
	old := commandBuilderFunc
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockCommandBuilder = mock_git_wrapper.NewMockICommandBuilder(mockCtrl)
		commandBuilderFunc = func() ICommandBuilder {
			return mockCommandBuilder
		}
	})
	AfterEach(func() {
		defer func() { commandBuilderFunc = old }()
	})
	Context("GenerateWorktree(row string)", func() {
		It("Should return correct Worktree struct", func() {
			result := GenerateWorktree(`/tmp/operarius  acde210 [issue-SS192-golang-git-package]`)
			Expect(result).Should(Equal(Worktree{
				CommitSHA: "acde210",
				Path:      "/tmp/operarius",
			}))
		})

		It("Should return correct Worktree struct", func() {
			result := GenerateWorktree(`/tmp/operarius-2  123123 [kai-branch]`)
			Expect(result).Should(Equal(Worktree{
				CommitSHA: "123123",
				Path:      "/tmp/operarius-2",
			}))
		})
	})

	Context("ListWorktree(path string) ([]Worktree, error)", func() {
		It("Should return list worktrees", func() {
			mockCommandBuilder.EXPECT().AddCommand("worktree")
			mockCommandBuilder.EXPECT().AddArg("list")
			mockCommandBuilder.EXPECT().SetDir("./tmp/core-api")
			mockCommandBuilder.EXPECT().Exec().Return("/tmp/operarius  74580d7 [issue-SS192-golang-git-package]\n"+"/tmp/operarius/kai-test  acde210 [issue-SS192-golang-git-package]", nil)
			result, _ := ListWorktree("./tmp/core-api")
			Expect(result).Should(Equal([]Worktree{
				{
					CommitSHA: "74580d7",
					Path:      "/tmp/operarius",
					IsMain:    true,
				},
				{
					CommitSHA: "acde210",
					Path:      "/tmp/operarius/kai-test",
					IsMain:    false,
				},
			}))
		})

		It("Should return only main worktree if the list is one", func() {
			mockCommandBuilder.EXPECT().AddCommand("worktree")
			mockCommandBuilder.EXPECT().AddArg("list")
			mockCommandBuilder.EXPECT().SetDir("./tmp/core-api")
			mockCommandBuilder.EXPECT().Exec().Return("/tmp/operarius  74580d7 [issue-SS192-golang-git-package]", nil)
			result, _ := ListWorktree("./tmp/core-api")
			Expect(result).Should(Equal([]Worktree{
				{
					CommitSHA: "74580d7",
					Path:      "/tmp/operarius",
					IsMain:    true,
				},
			}))
		})
	})
})
