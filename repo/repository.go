package repo

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/aol-nnov/debian/deb822"
	"github.com/aol-nnov/debian/internal/stringspp"
	"github.com/aol-nnov/debian/internal/universalreader"
)

type Repository struct {
	baseUrl                     string
	Origin                      string
	Label                       string
	Suite                       string
	Version                     string
	Codename                    string
	Changelogs                  string
	Date                        string
	AcquireByHash               bool     `deb822:"Acquire-By-Hash"`
	NoSupportforArchitectureall string   `deb822:"No-Support-for-Architecture-all"`
	Architectures               []string `delim:" "`
	Components                  []string `delim:" "`
	Description                 string
}

func New(baseUrl, codename string) (*Repository, error) {
	var err error

	repoReader, err := universalreader.New(
		fmt.Sprintf("%s/dists/%s/Release", baseUrl, codename))
	defer universalreader.MaybeClose(repoReader)

	if err != nil {
		return nil, err
	}

	var repository Repository
	err = deb822.NewDecoder(repoReader).Decode(&repository)
	if err != nil {
		return nil, err
	}
	repository.baseUrl = baseUrl

	return &repository, nil
}

func (r Repository) String() string {
	return stringspp.UniversalStringer(r)
}

func (r Repository) Component(name string) *Component {
	if slices.Contains(r.Components, name) {
		return &Component{
			name: name,
			repo: &r,
		}
	}
	return nil
	//, fmt.Errorf("component %s was not found", name)
}
