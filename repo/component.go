package repo

import (
	"fmt"

	"github.com/aol-nnov/debian/internal/universalreader"
	"golang.org/x/exp/slices"
)

type ComponentMap map[string]Component

type Component struct {
	repo *Repository
	name string
}

func (c *Component) SourceIndex() (*SourceIndex, error) {
	indexFile := fmt.Sprintf("%s/dists/%s/%s/source/Sources.gz",
		c.repo.baseUrl,
		c.repo.Codename,
		c.name)

	return NewSourceIndex(universalreader.New(indexFile))
}

func (c *Component) BinaryIndex(arch string) (*BinaryIndex, error) {
	if slices.Contains(c.repo.Architectures, arch) {
		indexFile := fmt.Sprintf("%s/dists/%s/%s/binary-%s/Packages.gz",
			c.repo.baseUrl,
			c.repo.Codename,
			c.name, arch)

		return NewBinaryIndex(universalreader.New(indexFile))
	}
	return nil, fmt.Errorf("%s %s no such architecture %s", c.repo.baseUrl, c.name, arch)
}
