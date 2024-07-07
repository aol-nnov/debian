package fields

import (
	"fmt"
	"strconv"
)

// https://www.debian.org/doc/manuals/developers-reference/pkgs.html#nmu
// https://www.debian.org/doc/manuals/developers-reference/pkgs.html#nmus-and-debian-changelog
// https://www.debian.org/doc/manuals/developers-reference/pkgs.html#recompilation-or-binary-only-nmu
// https://github.com/canonical/ubuntu-maintainers-handbook/blob/main/VersionStrings.md

/*
Call this method to denote a NMU

# Non-maintainer upload versioning scheme

If the package is a native package, the version must be the version of the last maintainer upload, plus +nmuX, where X is a counter starting at 1. If the last upload was also an NMU, the counter should be increased.

Examples:
  - 1.2.3 -> 1.2.3+nmu1
  - 2.4.1+nmu2 -> 2.4.1+nmu3

If the package is not a native package, a minor version number should be added to the Debian version. This extra number must start at 1.
If a new upstream version is packaged in the NMU, the Debian version is set to 0.

Examples:
  - 1.5-2 -> 1.5-2.1
  - new upstream ver 1.6 nmu	->	1.6-0.1

NOTE: NMU must be reflected in the changelog (as source is being uploaded)!

Subsequent call to [DebianVersion.Bump] will increase version and clear the nmu part.
*/
func (v *Version) Nmu() {
	lastModificator := v.Mod()

	switch v.IsMod() {
	case VersionModNmuNative:
		lastBuildNum, _ := strconv.Atoi(lastModificator[len("+nmu"):])
		v.UpdateMod(fmt.Sprintf("+nmu%d", lastBuildNum+1))
	case VersionModNmuQuilt:
		lastBuildNum, _ := strconv.Atoi(lastModificator[1:])
		v.UpdateMod(fmt.Sprintf(".%d", lastBuildNum+1))
	default:
		if v.IsNative() {
			v.AddMod("+nmu1")
		} else {
			v.AddMod(".1")
		}

	}
}

/*
Call this method to denote a binary-only NMU or rebuild

# Binary rebuild versioning scheme

A suffix appended to the package version number, following the form +b<number>. The rule is true for both, native and quilt packages.

Examples:
  - 1.2 -> 1.2+b1
  - 4.5-3 -> 4.5-3+b5

NOTE: make sure that your binary-only NMU doesn't render the package uninstallable. This could happen when a source package generates arch-dependent and arch-independent packages that have inter-dependencies generated using dpkg's substitution variable $(Source-Version).

Source package after the rebuild is not uploaded to the archive.
*/
func (v *Version) BinaryNmu() {
	if v.IsMod() == VersionModNmuBinary {
		lastModificatorIdx := len(v.Modificators) - 1
		lastModificator := v.Modificators[lastModificatorIdx]

		lastBuildNum, _ := strconv.Atoi(lastModificator[len("+b"):])

		v.UpdateMod(fmt.Sprintf("+b%d", lastBuildNum+1))

	} else {
		v.AddMod("+b1")
	}
}

func (v *Version) IsNmu() bool {
	switch v.IsMod() {
	case VersionModNmuBinary, VersionModNmuNative, VersionModNmuQuilt:
		return true
	}
	return false
}
