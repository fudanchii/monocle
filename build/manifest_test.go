package build

import (
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Build Manifest", func() {
	Context("when parse build manifest with simple shell runner", func() {
		var (
			bm  *Build
			err error
		)

		BeforeEach(func() {
			bm, err = ParseManifest("_fixtures/simple_shell.yml")
		})

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should parse Shell manifest correctly", func() {
			Ω(bm.Shell).ShouldNot(BeNil())
		})

		It("should parse Shell.Steps correctly", func() {
			Ω(strings.TrimSpace(bm.Shell.Steps)).Should(Equal("ls"))
		})
	})

	Context("when parse build manifest with simple docker run", func() {
		var (
			bm  *Build
			err error
		)

		BeforeEach(func() {
			bm, err = ParseManifest("_fixtures/simple_docker.yml")
		})

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should parse Docker manifest correctly", func() {
			Ω(bm.Docker).ShouldNot(BeNil())
			Ω(bm.Docker.Run).ShouldNot(BeNil())
		})

		It("should parse Docker image, steps, workdir, env, and volumes correctly", func() {
			run := bm.Docker.Run
			Ω(run.Image).Should(Equal("dummy/image"))

			Ω(strings.TrimSpace(run.Steps)).Should(Equal("ls"))

			Ω(run.Workdir).Should(Equal("test_path/yay"))

			Ω(run.Env).Should(HaveLen(1))
			Ω(run.Env[0]).Should(Equal("PATH=/nowhere/bin"))

			Ω(run.Volumes).Should(HaveLen(1))
			Ω(run.Volumes[0]).Should(Equal(".:/tmp"))
		})
	})

	Context("when parse build manifest with simple docker build", func() {
		var (
			bm  *Build
			err error
		)

		BeforeEach(func() {
			bm, err = ParseManifest("_fixtures/simple_docker_build.yml")
		})

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should parse Docker build manifest correctly", func() {
			Ω(bm.Docker).ShouldNot(BeNil())
			Ω(bm.Docker.Build).ShouldNot(BeNil())
		})

		It("should parse Docker build file, root, tags correctly", func() {
			build := bm.Docker.Build
			Ω(build.File).Should(Equal(".build/Dockerfile"))
			Ω(build.Root).Should(Equal(".build"))
			Ω(build.Tags).Should(HaveLen(1))
			Ω(build.Tags[0]).Should(Equal("test/image:latest"))
		})
	})

	Context("when parse build manifest with variables interpolation", func() {
		var (
			bm  *Build
			err error
		)

		BeforeEach(func() {
			os.Setenv("TEST_ENV", "hello")
			bm, err = ParseManifest("_fixtures/withvars_docker_build.yml")
		})

		It("should not error", func() {
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should parse Docker build manifest correctly", func() {
			Ω(bm.Docker).ShouldNot(BeNil())
			Ω(bm.Docker.Build).ShouldNot(BeNil())
		})

		It("should evaluate variables correctly", func() {
			Ω(bm.Variables).ShouldNot(BeNil())
		})

		It("should evaluate variables from env correctly", func() {
			env := bm.Variables.Env
			Ω(env).Should(HaveLen(1))
			Ω(env["test"]).Should(Equal("hello"))
		})

		It("should evaluate variables from eval correctly", func() {
			eval := bm.Variables.Eval
			Ω(eval).Should(HaveLen(2))
			Ω(eval["test"]).Should(Equal("ayyyyyy, hola!"))
			Ω(eval["forTag"]).Should(Equal("yay"))
		})

		It("should extrapolate tags with variables correctly", func() {
			Ω(bm.Docker.Build.Tags).Should(HaveLen(1))
			Ω(bm.Docker.Build.Tags[0]).Should(Equal("test/image:yay"))
		})
	})
})
