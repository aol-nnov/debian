package deb822_test

import (
	"fmt"
	"strings"

	"github.com/aol-nnov/debian/deb822"
)

func ExampleDecoder() {
	const deb822Stream = `Source: debusine
Section: devel
Priority: optional
`
	type Result struct {
		Source, Section, Priority string
	}
	var m Result

	deb822.NewDecoder(strings.NewReader(deb822Stream)).Decode(&m)
	fmt.Printf("%v", m)

	// Output: {debusine devel optional}
}

func ExampleDecoder_Decode_slice() {
	const deb822Stream = `Source: first
Section: devel
Priority: optional

Source: second
Section: devel
Priority: optional
`
	type Result struct {
		Source, Section, Priority string
	}
	var m []Result
	deb822.NewDecoder(strings.NewReader(deb822Stream)).Decode(&m)

	for _, r := range m {
		fmt.Println(r)
	}

	// Output:
	// {first devel optional}
	// {second devel optional}
}

func ExampleDecoder_Decode_stream() {
	const deb822Stream = `Source: first
Section: devel
Priority: optional

Source: second
Section: devel
Priority: optional

`
	type Result struct {
		Source, Section, Priority string
	}
	dec := deb822.NewDecoder(strings.NewReader(deb822Stream))
	var m Result
	for {
		if err := dec.Decode(&m); err != nil {
			break
		}

		fmt.Println(m)
	}

	// Output:
	// {first devel optional}
	// {second devel optional}
}
