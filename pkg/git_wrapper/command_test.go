package git_wrapper

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Command unit test", func() {
	Context("Build() string", func() {
		It("Should return command git checkout with commitSHA", func() {
			commandBuilder := NewCommandBuilder()
			commandBuilder.AddCommand("checkout")
			commandBuilder.AddArg("ebc635acded8305a60fec5fad5b66d9d8c74d78f")
			command := commandBuilder.Build()
			Expect(command).To(Equal("git checkout ebc635acded8305a60fec5fad5b66d9d8c74d78f"))
		})

		It("Should return command git pull", func() {
			commandBuilder := NewCommandBuilder()
			commandBuilder.AddCommand("pull")
			command := commandBuilder.Build()
			Expect(command).To(Equal("git pull"))
		})

		It("Should return git clone with ssh url git@github.com:nodegit/nodegit.git", func() {
			commandBuilder := NewCommandBuilder()
			commandBuilder.AddCommand("clone")
			commandBuilder.AddArg("git@github.com:nodegit/nodegit.git")
			command := commandBuilder.Build()
			Expect(command).To(Equal("git clone git@github.com:nodegit/nodegit.git"))
		})

		It("Should return git commit with message Init commit", func() {
			commandBuilder := NewCommandBuilder()
			commandBuilder.AddCommand("commit")
			commandBuilder.AddArgs([]string{"-m", "'Init commit'"})
			command := commandBuilder.Build()
			Expect(command).To(Equal("git commit -m 'Init commit'"))
		})
	})
})
