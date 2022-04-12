package cachectl

import (
	"testing"
)

func TestActivePages(t *testing.T) {
	_, err := activePages("activepagtes.go")
	if err != nil {
		t.Fatal(err)
	}
}
