package repo_test

import (
	"os"
	"testing"

	"github.com/aol-nnov/debian/deb822"
	"github.com/aol-nnov/debian/repo"
)

func TestMissingSourceField(t *testing.T) {
	in, _ := os.Open("./testdata/binaryIndex")

	var ii repo.BinaryIndexItem
	if err := deb822.NewDecoder(in).Decode(&ii); err != nil {
		t.Fatal(err)
	}

	if ii.Name != ii.Source {
		t.Fail()
	}
	t.Log(ii)
}
