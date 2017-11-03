package git

import (
	h "github.com/fudanchii/monocle/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Monocle Git Tests Suite")
}

var _ = Describe("Git", func() {
	var (
		cleanup = h.Noop
		err     error
		repo    string
	)

	BeforeEach(func() {
		repo, cleanup, err = h.CreateFixtureRepo()
		err = h.SeedSimpleCommit(repo, err)
	})

	AfterEach(func() {
		cleanup()
	})

	Context("test 1 file change", func() {
		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should return 1 change", func() {
			changes, err := FilesChanged(repo, "")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(changes.WorkDir).Should(Equal(repo))
			Ω(len(changes.Entries)).Should(Equal(1))
		})
	})

	Context("test multiple files changes", func() {
		BeforeEach(func() {
			err = h.SeedAnotherCommit(repo, err)
		})

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should have multiple changed entries", func() {
			changes, err := FilesChanged(repo, "")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(changes.WorkDir).Should(Equal(repo))
			Ω(len(changes.Entries)).Should(Equal(2))
		})
	})

})
