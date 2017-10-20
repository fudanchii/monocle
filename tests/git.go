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

	AfterEach(func() {
		cleanup()
	})

	Context("test 1 folder change", func() {
		BeforeEach(func() {
			repo, cleanup, err = createFixtureRepo()
			err = seedSimpleCommit(repo, err)
		})
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should return 1 change", func() {
			changes := git.FilesChanged(repo, "")
			Ω(changes.WorkDir).Should(Equal(repo))
			Ω(len(changes.Entries)).Should(Equal(1))
		})
	})
})
