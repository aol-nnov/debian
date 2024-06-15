package pkg_test

import (
	"fmt"
	"testing"

	"github.com/aol-nnov/debian/internal/universalreader"
	"github.com/aol-nnov/debian/pkg"
)

func TestContolParser(t *testing.T) {
	srcInput, err := universalreader.New("./testdata/control")

	if err != nil {
		t.Fatal(err)
	}

	var control pkg.Control

	if err := control.Decode(srcInput); err != nil {
		t.Fatal(err)
	}

	t.Logf("deb-src name %s", control.DebSrc.Name)
	t.Logf("%d binary packages", len(control.Deb))

	fmt.Println(control.Deb[0].Depends)

}
