package fields_test

import (
	"fmt"
	"testing"

	"github.com/aol-nnov/debian/fields"
)

type testcase struct {
	in             []byte
	clear          []byte
	activeProfiles []string
	expected       bool
}

var cases = []testcase{
	{
		in:             []byte("  <!p1>    "),
		clear:          []byte("<!p1>"),
		activeProfiles: []string{},
		expected:       true,
	},
	{
		in:             []byte("<p1>    "),
		clear:          []byte("<p1>"),
		activeProfiles: []string{},
		expected:       false,
	},
	{
		in:             []byte("<p1>"),
		clear:          []byte("<p1>"),
		activeProfiles: []string{"p1"},
		expected:       true,
	},
	{
		in:             []byte("<!p1>"),
		clear:          []byte("<!p1>"),
		activeProfiles: []string{"p1"},
		expected:       false,
	},
	{
		in:             []byte("<p1>"),
		clear:          []byte("<p1>"),
		activeProfiles: []string{"p1", "p2"},
		expected:       true,
	},
	{
		in:             []byte("<p1>"),
		clear:          []byte("<p1>"),
		activeProfiles: []string{"p2", "p1"},
		expected:       true,
	},
	{
		in:             []byte("<!p1>"),
		clear:          []byte("<!p1>"),
		activeProfiles: []string{"p1", "p2"},
		expected:       false,
	},
	{
		in:             []byte("<!p1>"),
		clear:          []byte("<!p1>"),
		activeProfiles: []string{"p2", "p3"},
		expected:       true,
	},
	{
		in:             []byte("<p1> <p2>"),
		clear:          []byte("<p1> <p2>"),
		activeProfiles: []string{"p2", "p1"},
		expected:       true,
	},
	{
		in:             []byte("<p1 p2> <p3>"),
		clear:          []byte("<p1 p2> <p3>"),
		activeProfiles: []string{"p2", "p1"},
		expected:       true,
	},
	{
		in:             []byte("   < p1   p2  >    < p3 >"),
		clear:          []byte("<p1 p2> <p3>"),
		activeProfiles: []string{"p3", "p1"},
		expected:       true,
	},
	{
		in:             []byte("<p1 !p2> <p3>"),
		clear:          []byte("<p1 !p2> <p3>"),
		activeProfiles: []string{"p1", "p3"},
		expected:       true,
	},
	{
		in:             []byte("<p1 !p2> <p3>"),
		clear:          []byte("<p1 !p2> <p3>"),
		activeProfiles: []string{"p1"},
		expected:       true,
	},
	{
		in:             []byte("<p1 !p2> <p3>"),
		clear:          []byte("<p1 !p2> <p3>"),
		activeProfiles: []string{"p1", "p2"},
		expected:       false,
	},
	{
		in:             []byte("<p1 p2> <!p3>"),
		clear:          []byte("<p1 p2> <!p3>"),
		activeProfiles: []string{"p1", "p3"},
		expected:       false,
	},
	{
		in:             []byte("<p1 !p2> <!p3>"),
		clear:          []byte("<p1 !p2> <!p3>"),
		activeProfiles: []string{"p1", "p2"},
		expected:       true,
	},
	{
		in:             []byte("<p1 !p2> <!p3>"),
		clear:          []byte("<p1 !p2> <!p3>"),
		activeProfiles: []string{"p1"},
		expected:       true,
	},
}

func TestProfileContraints(t *testing.T) {

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s vs %v", tc.in, tc.activeProfiles), func(t *testing.T) {
			var pc fields.ProfileConstraints

			if err := pc.UnmarshalText(tc.in); err != nil {
				t.Fatalf("failed to unmarshal '%s'", tc.in)
			}

			actual := pc.SatisfiedBy(tc.activeProfiles)
			if actual != tc.expected {
				t.Fatalf("%v vs %v is %v. Expected: %v", pc, tc.activeProfiles, actual, tc.expected)
			}
		})
	}
}

func TestProfileContraintsStringer(t *testing.T) {

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s vs %v", tc.in, tc.activeProfiles), func(t *testing.T) {
			var pc fields.ProfileConstraints

			if err := pc.UnmarshalText(tc.in); err != nil {
				t.Fatalf("failed to unmarshal '%s'", tc.in)
			}

			if string(tc.clear) != pc.String() {
				t.Fatalf("stringer failed on '%s'", tc.in)
			}
			// t.Logf("'%s' -> '%v'", tc.in, pc)
		})
	}
}
