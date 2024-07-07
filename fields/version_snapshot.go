package fields

import (
	"fmt"
	"strconv"
	"strings"
)

/*
Call this method wile building a git snapshot version (i.e. merge request)

# Snapshot versioning scheme

To denote a snapshot version add ~<buildNum>.gbp<short-commit-sha> to the debian version. The rule is true for both, native and quilt packages.

	Examples:
		- 1.4.2-7 -> 1.4.2-7~1.gbp876ad, 1.4.2-7~2.gbpdeadbe,...
		- 1.4.2 -> 1.4.2~1.gbp876ad, 1.4.2~2.gbpdeadbe,...

Subsequent call to [Version.Bump] will increase version and clear the snapshot part.

NOTE: set distribution to UNRELEASED in changelog for snapshot entry
*/
func (v *Version) Snapshot(commitSha string) {

	if len(commitSha) > 8 {
		commitSha = commitSha[:8]
	}

	for v.IsNmu() {
		v.RemoveMod()
	}

	if v.IsMod() == VersionModSnapshot {
		if snapshotNumStr, snapshotSuffix, dotFound := strings.Cut(v.Mod(), "."); dotFound {
			if strings.HasPrefix(snapshotSuffix, "gbp") {
				snapshotNum, _ := strconv.Atoi(snapshotNumStr[1:])
				v.UpdateMod(fmt.Sprintf("~%d.gbp%s", snapshotNum+1, commitSha))
			}
		}

	} else {
		v.AddMod(fmt.Sprintf("~1.gbp%s", commitSha))
	}
}
