package repo

import (
	"fmt"
	"io"
	"slices"

	"github.com/aol-nnov/debian/deb822"
	"github.com/aol-nnov/debian/fields"
	"github.com/aol-nnov/debian/internal/universalreader"
	"golang.org/x/exp/maps"
)

type SourceIndex struct {
	Packages []SourceIndexItem

	// []int for storing indices to multiple versions of single package
	byName       map[string][]int
	byBinaryName map[string]int
}

type SourceIndexItem struct {
	Name    string         `control:"Package"`
	Binary  []string       `delim:"," strip:" " required:"true"`
	Version fields.Version `required:"true"`

	BuildDepends      fields.Dependencies `control:"Build-Depends" delim:"," strip:" "`
	BuildDependsArch  fields.Dependencies `control:"Build-Depends-Arch" delim:"," strip:" "`
	BuildDependsIndep fields.Dependencies `control:"Build-Depends-Indep" delim:"," strip:" "`

	Architecture []fields.Architecture `required:"true" delim:" " strip:" "`
	// StandardsVersion string `control:"Standards-Version"`
	// Format  string
}

func (p SourceIndexItem) String() string {
	// return strings.UniversalStringer(p)
	return fmt.Sprintf("%s %s", p.Name, p.Version)
}

func NewSourceIndex(reader io.Reader, inErr error) (*SourceIndex, error) {
	if inErr != nil {
		return nil, inErr
	}
	var res SourceIndex

	err := deb822.NewDecoder(reader).Decode(&res.Packages)
	defer universalreader.MaybeClose(reader)

	if err != nil {
		return nil, err
	}

	res.byName = make(map[string][]int, len(res.Packages))
	res.byBinaryName = make(map[string]int, len(res.Packages))

	// for idx, pkg := range res.Packages {
	// 	// let's cache index of deb-src with max version
	// 	if cachedIdx, found := res.byName[pkg.Name]; found {
	// 		if res.Packages[cachedIdx].Version.Less(pkg.Version) {
	// 			// fmt.Println(res.Packages[cachedIdx], pkg)
	// 			res.byName[pkg.Name] = idx
	// 		}
	// 	} else {
	// 		res.byName[pkg.Name] = idx
	// 	}

	// 	// for _, bin := range pkg.Binary {
	// 	// 	res.byBinaryName[bin] = idx
	// 	// }
	// }

	for idx, pkg := range res.Packages {
		res.byName[pkg.Name] = append(res.byName[pkg.Name], idx)

		for _, bin := range pkg.Binary {
			res.byBinaryName[bin] = idx
		}
	}

	// for name, i := range res.byName {
	// 	fmt.Println(name, len(i))
	// }

	return &res, nil
}

func (si SourceIndex) FindByName(name string) ([]SourceIndexItem, bool) {

	if packages, found := si.byName[name]; found {
		res := make([]SourceIndexItem, 0, len(packages))
		for _, idx := range packages {
			res = append(res, si.Packages[idx])
		}
		// return si.Packages[idx], true
		return res, true
	}

	return nil, false
}

func (si SourceIndex) FindByBinaryName(name string) (*SourceIndexItem, bool) {
	if idx, found := si.byBinaryName[name]; found {
		return &si.Packages[idx], true
	}

	return nil, false
}

func (si SourceIndex) FindByConstraint(dep fields.Dependency) (*SourceIndexItem, bool) {

	return nil, false
}

// Kahn toposort algorithm
// https://github.com/amwolff/gorder/blob/master/gorder.go
//
// alternatives:
// depth first traversal
// https://github.com/ninedraft/tsort/blob/master/sort.go (colour sort)
// https://github.com/paultag/go-topsort/blob/master/topsort.go
// TODO: найти алгоритм, разделяющий граф на сортированные поддеревья - возможность распараллелить задачу сборки
func (si SourceIndex) BuildOrder() (topologicalOrder []string, missing []string, err error) {
	indegree := make(map[string]int)
	pkgDeps := make(map[string][]string)
	missingQueue := make(map[string]bool, len(si.Packages))
	var queue []string

	for pkgName := range si.byName {

		pkg := si.Packages[si.byName[pkgName][0]]

		for _, binDependency := range pkg.BuildDepends {
			if srcDep, found := si.FindByBinaryName(binDependency.Name); found {
				pkgDeps[pkgName] = append(pkgDeps[pkgName], srcDep.Name)
				indegree[srcDep.Name]++
			} else {
				missingQueue[binDependency.Name] = true
			}
		}
	}

	for pkgName := range si.byName {
		if _, ok := indegree[pkgName]; !ok {
			queue = append(queue, pkgName)
		}
	}

	cnt := 0
	for len(queue) > 0 {
		pkgName := queue[len(queue)-1]
		queue = queue[:(len(queue) - 1)] // pops last element from queue

		topologicalOrder = append(topologicalOrder, pkgName)

		for _, dep := range pkgDeps[pkgName] {
			indegree[dep]--
			if indegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
		cnt++
	}

	missing = maps.Keys(missingQueue)

	if cnt != len(si.byName) {
		err = fmt.Errorf("cyclic graph: visited %d, total %d", cnt, len(si.byName))
	}

	slices.Reverse(topologicalOrder)
	return
}

// func (si SourceIndex) BuildOrderIndices() (topologicalOrder []int, missing []string, err error) {
// 	indegree := make(map[int]int)                           // map[pkgIndex]IndegreeCount
// 	pkgDeps := make(map[int][]int)                          // map[pkgIndex][]depIndex
// 	missingQueue := make(map[string]bool, len(si.Packages)) //map[binaryPkgName]true
// 	var queue []string

// 	for pkgName := range si.byName {

// 		pkg := si.Packages[si.byName[pkgName][0]]

// 		for _, binDependency := range pkg.BuildDepends {
// 			if srcDep, found := si.FindByBinaryName(binDependency.Name); found {
// 				pkgDeps[pkgName] = append(pkgDeps[pkgName], srcDep.Name)
// 				indegree[srcDep.Name]++
// 			} else {
// 				missingQueue[binDependency.Name] = true
// 			}
// 		}
// 	}

// 	for pkgName := range si.byName {
// 		if _, ok := indegree[pkgName]; !ok {
// 			queue = append(queue, pkgName)
// 		}
// 	}

// 	cnt := 0
// 	for len(queue) > 0 {
// 		pkgName := queue[len(queue)-1]
// 		queue = queue[:(len(queue) - 1)]

// 		topologicalOrder = append(topologicalOrder, pkgName)

// 		for _, dep := range pkgDeps[pkgName] {
// 			indegree[dep]--
// 			if indegree[dep] == 0 {
// 				queue = append(queue, dep)
// 			}
// 		}
// 		cnt++
// 	}

// 	missing = maps.Keys(missingQueue)

// 	if cnt != len(si.byName) {
// 		err = fmt.Errorf("cyclic graph: visited %d, total %d", cnt, len(si.byName))
// 	}

// 	slices.Reverse(topologicalOrder)
// 	return
// }
