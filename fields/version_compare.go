package fields

func (v Version) Less(another Version) bool {
	return v.Compare(another) == VersionCompareResultLessThan
}

/*
Compares two debian `Version`s

The `upstream_version` and `debian_revision` parts are compared by the package management system using the same algorithm:

1. The strings are compared from left to right.

2. First the initial part of each string consisting entirely of non-digit characters is determined. These two parts (one of
which may be empty) are compared lexically. If a difference is found it is returned. The lexical comparison is a
comparison of ASCII values modified so that all the letters sort earlier than all the non-letters and so that a tilde
sorts before anything, even the end of a part. For example, the following parts are in sorted order: '~~', '~~a', '~',
the empty part, 'a'.

3. Then the initial part of the remainder of each string which consists entirely of digit characters is determined. The
numerical values of these two parts are compared, and any difference found is returned as the result of the comparison.
For these purposes an empty string (which can only occur at the end of one or both version strings being compared)
counts as zero.

4. These two steps (comparing and removing initial non-digit strings and initial digit strings) are repeated until a
difference is found or both strings are exhausted.

*Note* that the purpose of epochs is to allow us to leave behind mistakes in version numbering, and to cope with
situations where the version numbering scheme changes. It is not intended to cope with version numbers containing
strings of letters which the package management system cannot interpret (such as 'ALPHA' or 'pre-'), or with silly
orderings.
*/
func (v Version) Compare(another Version) VersionCompareResult {
	if v.raw != "" || another.raw != "" {
		return VersionCompareResultNonComparable
	}

	if v.Epoch > another.Epoch {
		return VersionCompareResultGreaterThan
	}

	return verrevcmp(v.String(), another.String())
}
