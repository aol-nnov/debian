package fields

/*
	 *

	 	- https://www.debian.org/doc/debian-policy/ch-relationships.html
		- https://wiki.debian.org/BuildProfileSpec
		- https://www.debian.org/doc/debian-policy/ch-customized-programs.html#s-arch-spec
		- https://wiki.debian.org/CrossBuildPackagingGuidelines#Architecture_qualifiers
		- https://www.debian.org/doc/debian-policy/ch-controlfields.html#s-f-vcs-fields
		- man deb-src-control Build-Depends

	    Represents a list of groups of alternative packages.

	    1. Each group is a list of packages separated by vertical  bar  (or“pipe”)  symbols,  ‘|’.
	    2. The  groups  are separated by commas ‘,’, and can end with a trailing comma that will be eliminated when generating the fields for deb-control(5) (since dpkg 1.10.14). Commas are to be read as “AND”, and pipes as “OR”, with pipes binding more tightly.
	    3. Each package name is optionally followed by an architecture qualifier appended after a colon ‘:’, optionally followed by a version number specification in parentheses ‘(’ and ‘)’, an architecture specification in square brackets ‘[’ and ‘]’, and a restriction formula consisting of one or more lists of

profile names in angle brackets ‘<’ and ‘>’.

	dependencies := group, group
	group := item | item
	item := name:qualifier versionspec archspec profilespec
	name := string
	qualifier := string

	versionspec := (comparator debversion)
	comparator := << <= = >= >>
	debversion := epoch:upstream-debian

	negate := !

	archspec := [archconstraint,...]
	archconstraint: := negate debianarch
	debianarch := abispec-libcspec-osspec-cpuspec

full packed example: pkg-name:any (<< 1.2.3) [amd64 !hurd-any] <profile1 !profile2> <profile3> | another-name
*/
type Dependency struct {
	Name                    string
	ArchQualifier           string // "", "any", "native"
	VersionConstraint       *VersionConstraint
	ProfileConstraints      ProfileConstraints
	ArchitectureConstraints ArchitectureConstraints
	Alt                     *Dependency
	raw                     string // for templated dependency in control file like ${python3:Depends}
}

/*
Terminology:

  - BUILD is the machine we are building on
  - HOST is the machine we are building for
  - TARGET (is only relevant for compilers and is the architecture that a compiler outputs code for. Unless packaging binutils, gcc or hurd, the target architecture is irrelevant.)

This somewhat confusing terminology is GNU's fault. :clown_face:
*/
func (dep Dependency) Satisfies(buildArch Architecture, hostArch Architecture, profiles []string) bool {
	dc := true
	ac := dep.ArchitectureConstraints.SatisfiedBy(hostArch)
	pc := dep.ProfileConstraints.SatisfiedBy(profiles)

	// fmt.Printf("%s: ac %v, pc %v\n", dep.Name, ac, pc)
	return dc && ac && pc
}
