package fields

import "strings"

type VersionModificators []string

type VersionMod int

const (
	VersionModNone VersionMod = iota
	VersionModNmuNative
	VersionModNmuQuilt
	VersionModNmuBinary
	VersionModSnapshot
)

func extractVersionModificators(in string, delimiters string) (
	versionWithoutModificators string,
	modificators []string) {

	modStartIdx := 0

	for currentRuneIdx, r := range in {
		if strings.ContainsRune(delimiters, r) {
			if versionWithoutModificators == "" {
				versionWithoutModificators = in[:currentRuneIdx]
				modStartIdx = currentRuneIdx
			}

			if modStartIdx != currentRuneIdx {
				modificators = append(modificators, in[modStartIdx:currentRuneIdx])
			}

			modStartIdx = currentRuneIdx
		}
	}

	if modStartIdx != 0 {
		modificators = append(modificators, in[modStartIdx:])
		return
	}

	return in, nil
}

func (v *Version) IsMod() VersionMod {
	if len(v.Modificators) == 0 {
		return VersionModNone
	}

	lastModificator := v.Mod()
	if strings.Contains(lastModificator, "+b") {
		// binary nmu
		return VersionModNmuBinary
	}

	if strings.Contains(lastModificator, "+nmu") {
		// native nmu
		return VersionModNmuNative
	}

	if v.IsQuilt() && lastModificator[0] == '.' {
		// quilt nmu
		return VersionModNmuQuilt
	}

	if strings.Contains(lastModificator, "gbp") {
		return VersionModSnapshot
	}

	return VersionModNone
}

func (v *Version) Mod() string {
	lastIdx := len(v.Modificators) - 1

	if lastIdx >= 0 {
		return v.Modificators[lastIdx]
	}

	return ""
}

func (v *Version) AddMod(mod string) {
	v.Modificators = append(v.Modificators, mod)
}

func (v *Version) RemoveMod() {
	lastIdx := len(v.Modificators) - 1

	if lastIdx >= 0 {
		v.Modificators = v.Modificators[:lastIdx]
	}
}

func (v *Version) UpdateMod(mod string) {
	lastIdx := len(v.Modificators) - 1

	if lastIdx >= 0 {
		v.Modificators[lastIdx] = mod
	}
}
