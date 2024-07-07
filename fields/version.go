package fields

/*
Represents Debian package version number format

  - [Man page](https://linux.die.net/man/5/deb-version)
  - [Reference implementation](https://salsa.debian.org/dpkg-team/dpkg/-/blob/main/lib/dpkg/version.c)
  - Other good stuff: https://github.com/sdumetz/node-deb-version-compare/blob/master/lib/Version.js

BNF descriptor:

	<Version> ::= (<Epoch> ":")? (<UpstreamVersion> "-")? <DebianRevision>
	<Epoch> ::= <Num>+
	<UpstreamVersion> ::= <VersionPart> | <UpstreamVersion>  <Separator> (<Num> | <Alpha>)+
	<DebianRevision> ::= <VersionPart> | <DebianRevision>  <Separator> (<Num> | <Alpha>)+
	<VersionPart> ::= <Num> (<Num> | <Alpha>)*
	<Num> ::= [0-9]
	<Alpha> ::= ([A-Z] | [a-z])
	<Separator> ::= "." | "+" | "~"
*/
type Version struct {

	/*
	   This is a single (generally small) unsigned integer. It may be omitted in which case zero is assumed. If it is
	   omitted then the upstream_version may not contain any colons. It is provided to allow mistakes in the version
	   numbers of older versions of a package, and also a package's previous version numbering schemes, to be left
	   behind.
	*/
	Epoch int

	/*
		This is the main part of the version number. It is usually the version number of the original ("upstream")
		package from which the .deb file has been made, if this is applicable. Usually this will be in the same
		format as that specified by the upstream author(s); however, it may need to be reformatted to fit into the
		package management system's format and comparison scheme. The comparison behavior of the package management
		system with respect to the upstream_version is described below. The upstream_version portion of the version
		number is mandatory.

		The upstream_version may contain only alphanumerics ("A-Za-z0-9") and the characters . + - : ~ (full stop,
		plus, hyphen, colon, tilde) and should start with a digit. If there is no debian_revision then hyphens are
		not allowed; if there is no epoch then colons are not allowed.

		Из вышенаписанного делаю вывод, что : используется для отделения epoch от версии, а - для отделения UpstreamVersion от DebianRevision. Следовательно, для различных "модифиваторов" остаются ~ и +.
		Например ~N.gpp<commit_sha>, +bN или +nmuN
	*/
	UpstreamVersion string

	/*
		This part of the version number specifies the version of the Debian package based on the upstream version. It
		may contain only alphanumerics and the characters + . ~ (plus, full stop, tilde) and is compared in the same way
		as the upstream_version is. It is optional; if it isn't present then the upstream_version may not contain a
		hyphen. This format represents the case where a piece of software was written specifically to be turned into a
		Debian package, and so there is only one "debianisation" of it and therefore no revision indication is required.

		It is conventional to restart the debian_revision at '1' each time time the upstream_version is increased.

		Dpkg will break the version number apart at the last hyphen in the string (if there is one) to determine the
		upstream_version and debian_revision. The absence of a debian_revision compares earlier than the presence of one
		(but note that the debian_revision is the least significant part of the version number).
	*/
	DebianRevision string

	Modificators VersionModificators

	// for templated version in control file like ${source:Version}
	raw string
}

func MakeVersion(v string) Version {
	var res Version
	res.UnmarshalText([]byte(v))

	return res
}

func (v *Version) IsNative() bool {
	return v.DebianRevision == ""
}

func (v *Version) IsQuilt() bool {
	return v.DebianRevision != ""
}

func (v *Version) NewEpoch() {
	v.Epoch++
}
