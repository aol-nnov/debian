package fields

import (
	"bytes"
	"encoding"
	"fmt"
	"strings"

	"github.com/aol-nnov/debian/internal/stringspp"
)

var altSeparator = []byte{'|'}
var space = []byte{' '}

func (d *Dependency) UnmarshalText(text []byte) (err error) {
	// fmt.Printf("unmarshaling dependency '%s'\n", string(text))

	text = bytes.TrimSpace(text)

	// save templates as is
	if bytes.HasPrefix(text, []byte{'$'}) {
		d.raw = string(text)
		return
	}

	var primary, alt []byte
	var found bool

	if primary, alt, found = bytes.Cut(text, altSeparator); found {
		d.Alt = &Dependency{}
		err = d.Alt.UnmarshalText(alt)
	}

	if pkgNameMayBeQual, rest, found := bytes.Cut(primary, space); found {
		d.Name = string(pkgNameMayBeQual)

		/*
			parsed so far:
				- "pkgname:qual"
			 rest = "(>= 1.2.3) [arch] <!profile>"
		*/
		if versionConstraint, found, tail := stringspp.Between(rest, '(', ')', false); found {
			d.VersionConstraint = &VersionConstraint{}
			if err = d.VersionConstraint.UnmarshalText(versionConstraint); err != nil {
				return err
			}

			// rest = rest[bytes.IndexByte(rest, ')'):]
			rest = tail
			/*
				parsed so far:
					- "pkgname:qual"
					- "(>= 1.2.3)"
				 rest = " [arch] <!profile>"
			*/
		}

		if archConstraints, found, tail := stringspp.Between(rest, '[', ']', false); found {
			d.ArchitectureConstraints.UnmarshalText(archConstraints)
			rest = tail
		}

		if profileConstraints, found, _ := stringspp.Between(rest, '<', '>', true); found {
			d.ProfileConstraints.UnmarshalText(profileConstraints)
		}
	} else { // pkg name without constraints
		d.Name = string(text)
		d.ArchitectureConstraints = make([]architectureConstraint, 0)
		// d.ProfileConstraints = make([]ProfileConstraint, 0)
	}

	if pkgName, archQual, found := strings.Cut(d.Name, ":"); found {
		d.ArchQualifier = archQual
		d.Name = pkgName
	}

	return
	// return fmt.Errorf("Dependency Unmarshal malformed string '%s'", text)
}

func (d Dependency) String() (res string) {
	if d.raw != "" {
		return d.raw
	}

	res = d.Name

	if d.ArchQualifier != "" {
		res += fmt.Sprintf(":%s", d.ArchQualifier)
	}

	if d.VersionConstraint != nil {
		res += fmt.Sprintf(" %s", d.VersionConstraint)
	}

	if len(d.ArchitectureConstraints) != 0 {
		res += fmt.Sprintf(" [%s]", d.ArchitectureConstraints)
	}

	if len(d.ProfileConstraints) != 0 {
		res += fmt.Sprintf(" %s", d.ProfileConstraints)
	}

	return
}

func (d *Dependency) MarshalText() (text []byte, err error) {
	text = []byte(d.String())

	if d.Alt != nil {
		text = fmt.Appendf(text, " | %s", d.Alt.String())
	}

	return
}

var _ encoding.TextMarshaler = (*Dependency)(nil)
var _ encoding.TextUnmarshaler = (*Dependency)(nil)
