package tests

import (
	"github.com/fudanchii/monocle/git"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Git", func() {
	var (
		cleanup = noop
		err     error
		repo    string
	)

	BeforeEach(func() {
		repo, cleanup, err = createFixtureRepo()
		err = seedSimpleCommit(repo, err)
	})

	AfterEach(func() {
		cleanup()
	})

	Context("test 1 file change", func() {
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should return 1 change", func() {
			changes, err := git.FilesChanged(repo, "")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(changes.WorkDir).Should(Equal(repo))
			Ω(len(changes.Entries)).Should(Equal(1))
		})
	})

	Context("test multiple files changes", func() {
		BeforeEach(func() {
			err = seedAnotherCommit(repo, err)
		})

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should have multiple changed entries", func() {
			changes, err := git.FilesChanged(repo, "")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(changes.WorkDir).Should(Equal(repo))
			Ω(len(changes.Entries)).Should(Equal(2))
		})
	})

})
