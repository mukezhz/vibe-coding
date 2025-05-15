package todo_test

import (
	"clean-architecture/pkg/utils"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTodo(t *testing.T) {
	utils.ChDir()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Todo Suite")
}
