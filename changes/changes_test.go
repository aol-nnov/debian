package changes_test

import (
	"os"
	"testing"

	"github.com/aol-nnov/debian/changes"
)

func TestReadClearSign(t *testing.T) {
	in, _ := os.Open("./testdata/notebook_3.2.9_amd64.changes")

	if c, err := changes.FromStream(in); err == nil {
		t.Log(c.Files)
	} else {
		t.Error(err)
	}

}

func TestReadNotSigned(t *testing.T) {
	in, _ := os.Open("./testdata/notebook_3.2.9_amd64-unsigned.changes")

	if c, err := changes.FromStream(in); err == nil {
		t.Log(c.Changes)
	} else {
		t.Error(err)
	}

}
