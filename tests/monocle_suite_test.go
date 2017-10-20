package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMonocle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Monocle Tests Suite")
}
