package fields

import "fmt"

func (v *Version) NewUpstream(upstreamVersion string) error {

	if verrevcmp(upstreamVersion, v.UpstreamVersion) != VersionCompareResultGreaterThan {
		return fmt.Errorf("new upstream version must be greater than %s", v.UpstreamVersion)
	}

	v.UpstreamVersion = upstreamVersion
	v.DebianRevision = "1"

	return nil
}

func (v *Version) RollBackTo(goodVersion string) error {
	if v.UpstreamVersion == "" {
		return fmt.Errorf("RollBackTo is not supported for Debian native packages")
	}

	badVersion := v.UpstreamVersion
	v.UpstreamVersion = fmt.Sprintf("%s+really%s", badVersion, goodVersion)
	v.DebianRevision = "1"

	return nil
}
