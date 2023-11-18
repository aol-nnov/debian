package stringspp_test

import (
	"bytes"
	"testing"

	"github.com/aol-nnov/debian/internal/stringspp"
)

func TestBetween(t *testing.T) {
	in := []byte("<abc>_def")

	if match, found, rest := stringspp.Between(in, '<', '>', false); found {
		t.Logf("match '%s'", match)
		t.Logf("rest '%s'", rest)
		if !bytes.Equal(match, []byte("<abc>")) {
			t.Fail()
		}

		if !bytes.Equal(rest, []byte("_def")) {
			t.Fail()
		}
	}
}

func TestBetweenGreedy(t *testing.T) {
	in := []byte("moo <abc> <qwe> def")

	if match, found, rest := stringspp.Between(in, '<', '>', true); found {
		t.Logf("match '%s'", match)
		t.Logf("rest '%s'", rest)
		if !bytes.Equal(match, []byte("<abc> <qwe>")) {
			t.Fail()
		}

		if !bytes.Equal(rest, []byte(" def")) {
			t.Fail()
		}
	}
}
