package repo_test

import (
	"fmt"
	"testing"

	"github.com/aol-nnov/debian/internal/universalreader"
	"github.com/aol-nnov/debian/repo"
)

func TestTopo(t *testing.T) {
	src := "./testdata/Sources"

	si, _ := repo.NewSourceIndex(universalreader.New(src))
	// fmt.Println(si.Packages)
	if pkg, found := si.FindByName("gateconfig"); found {
		fmt.Println(pkg)
	}
}
