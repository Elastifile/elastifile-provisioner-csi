package size_test

import (
	. "tester"

	. "github.com/onsi/gomega"

	"testing"
)

func TestSize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Size Suite")
}
