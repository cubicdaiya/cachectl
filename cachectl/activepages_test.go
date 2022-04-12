package cachectl

import (
	"testing"
)

func TestActivePages(t *testing.T) {
	res, err := activePages("activepagtes.go")
	if err != nil {
		t.Fatal(err)
	}
}
