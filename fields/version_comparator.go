package fields

type VersionCompareResult int

const (
	VersionCompareResultLessThan VersionCompareResult = iota - 1
	VersionCompareResultEquals
	VersionCompareResultGreaterThan
	VersionCompareResultNonComparable
)

func (c VersionCompareResult) String() string {
	return [...]string{"<<", "=", ">>", "!!"}[c+1]
}

func cisdigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func cisalpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func order(r rune) int {
	if cisdigit(r) {
		return 0
	}
	if cisalpha(r) {
		return int(r)
	}
	if r == '~' {
		return -1
	}
	if int(r) != 0 {
		return int(r) + 256
	}
	return 0
}

func intToVersionCompareResult(in int) VersionCompareResult {
	if in < 0 {
		return VersionCompareResultLessThan
	} else {
		return VersionCompareResultGreaterThan
	}
}

func verrevcmp(a string, b string) VersionCompareResult {
	i := 0
	j := 0
	for i < len(a) || j < len(b) {
		var first_diff int
		for (i < len(a) && !cisdigit(rune(a[i]))) ||
			(j < len(b) && !cisdigit(rune(b[j]))) {
			ac := 0
			if i < len(a) {
				ac = order(rune(a[i]))
			}
			bc := 0
			if j < len(b) {
				bc = order(rune(b[j]))
			}
			if ac != bc {
				return intToVersionCompareResult(ac - bc)
			}
			i++
			j++
		}

		for i < len(a) && a[i] == '0' {
			i++
		}
		for j < len(b) && b[j] == '0' {
			j++
		}

		for i < len(a) && cisdigit(rune(a[i])) && j < len(b) && cisdigit(rune(b[j])) {
			if first_diff == 0 {
				first_diff = int(rune(a[i]) - rune(b[j]))
			}
			i++
			j++
		}

		if i < len(a) && cisdigit(rune(a[i])) {
			return VersionCompareResultGreaterThan
		}
		if j < len(b) && cisdigit(rune(b[j])) {
			return VersionCompareResultLessThan
		}
		if first_diff != 0 {
			return intToVersionCompareResult(first_diff)
		}
	}
	return VersionCompareResultEquals
}
