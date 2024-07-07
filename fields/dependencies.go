package fields

type Dependencies []Dependency

func (d Dependencies) For(buildArch Architecture, hostArch Architecture, profiles []string) Dependencies {
	var res Dependencies

	for _, dep := range d {
		if dep.Satisfies(buildArch, hostArch, profiles) {
			res = append(res, dep)
		}
	}
	return res
}

func (d Dependencies) NamesFor(buildArch Architecture, hostArch Architecture, profiles []string) []string {
	var res []string

	for _, dep := range d {
		if dep.Satisfies(buildArch, hostArch, profiles) {
			res = append(res, dep.Name)
		}
	}
	return res
}
