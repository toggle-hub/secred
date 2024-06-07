package utils_test

import (
	"os"
	"testing"

	"github.com/xsadia/secred/pkg/utils"
)

func TestOr(t *testing.T) {
	if actual := utils.Or(os.Getenv("something"), "something else"); actual != "something else" {
		t.Errorf("Expected something else, got %v", actual)
	}
}
