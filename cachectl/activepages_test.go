package cachectl

import (
	"testing"
)

func TestActivePages(t *testing.T) {
	_, err := activePages("activepages.go")
	if err != nil {
		t.Fatal(err)
	}
}
