package git_wrapper

import (
	mock_git_wrapper "operarius/mock/pkg/git_wrapper"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Commit unit test", func() {
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
	Context("DiffListFileChange(targetCommit *Commit) ([]string, error)", func() {
		It("Should return the list file changed between 2 commits", func() {
			commit := NewCommit("ebc635acded8305a60fec5fad5b66d9d8c74d78f", "/tmp/scan")
			targetCommit := NewCommit("99cdb715ac9cdad0f90f6af6df2757661b117efb", "/tmp/scan")
			mockCommandBuilder.EXPECT().AddCommand("diff")
			mockCommandBuilder.EXPECT().SetDir("/tmp/scan")
			mockCommandBuilder.EXPECT().AddArgs([]string{"--name-only", "--diff-filter=ACMR"})
			mockCommandBuilder.EXPECT().AddArg("99cdb715ac9cdad0f90f6af6df2757661b117efb")
			mockCommandBuilder.EXPECT().AddArg("ebc635acded8305a60fec5fad5b66d9d8c74d78f")
			mockCommandBuilder.EXPECT().Exec().Return("pkg/git_wrapper/git.go", nil).Times(1)
			output, _ := commit.DiffListFileChanged(&targetCommit)
			Expect(output).To(Equal([]string{"pkg/git_wrapper/git.go"}))
		})

		It("Should return the current file changes of the last commit", func() {
			commit := NewCommit("ebc635acded8305a60fec5fad5b66d9d8c74d78f", "/tmp/scan")
			mockCommandBuilder.EXPECT().AddCommand("diff")
			mockCommandBuilder.EXPECT().SetDir("/tmp/scan")
			mockCommandBuilder.EXPECT().AddArgs([]string{"--name-only", "--diff-filter=ACMR"})
			mockCommandBuilder.EXPECT().AddArg("ebc635acded8305a60fec5fad5b66d9d8c74d78f")
			mockCommandBuilder.EXPECT().Exec().Return(
				"go.mod\n"+"go.sum\n"+"pkg/git_wrapper/branch.go",
				nil).Times(1)
			output, _ := commit.DiffListFileChanged(nil)
			Expect(output).To(Equal([]string{
				"go.mod",
				"go.sum",
				"pkg/git_wrapper/branch.go",
			}))
		})
	})
})

func TestCommit_DiffListFileChanged(t *testing.T) {
	var mockCtrl *gomock.Controller
	var mockCommandBuilder *mock_git_wrapper.MockICommandBuilder
	mockCtrl = gomock.NewController(t)
	mockCommandBuilder = mock_git_wrapper.NewMockICommandBuilder(mockCtrl)
	commandBuilderFunc = func() ICommandBuilder {
		return mockCommandBuilder
	}

	old := commandBuilderFunc

	defer func() { commandBuilderFunc = old }()

	commit := NewCommit("ebc635acded8305a60fec5fad5b66d9d8c74d78f", "/tmp/scan")
	targetCommit := NewCommit("99cdb715ac9cdad0f90f6af6df2757661b117efb", "/tmp/scan")
	mockCommandBuilder.EXPECT().AddCommand("diff")
	mockCommandBuilder.EXPECT().SetDir("/tmp/scan")
	mockCommandBuilder.EXPECT().AddArgs([]string{"--name-only", "--diff-filter=ACMR"})
	mockCommandBuilder.EXPECT().AddArg("99cdb715ac9cdad0f90f6af6df2757661b117efb")
	mockCommandBuilder.EXPECT().AddArg("ebc635acded8305a60fec5fad5b66d9d8c74d78f")
	mockCommandBuilder.EXPECT().Exec().Return("pkg/git_wrapper/git.go", nil).Times(1)
	output, _ := commit.DiffListFileChanged(&targetCommit)
	expected := []string{"pkg/git_wrapper/git.go"}
	if !reflect.DeepEqual(output, expected) {
		t.Fatalf("got %v, expect %v", output, expected)
	}
}
