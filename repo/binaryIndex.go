package repo

import (
	"io"

	"github.com/aol-nnov/debian/deb822"
	"github.com/aol-nnov/debian/internal/stringspp"
)

type BinaryIndex struct {
	Items []BinaryIndexItem

	byName       map[string]int
	bySourceName map[string]int
}

type BinaryIndexItem struct {
	Name           string   `deb822:"Package" required:"true"`
	Provides       []string `delim:"," strip:" "`
	Source         string   `if_missing:"Package"`
	InstalledSize  int      `deb822:"Installed-Size"`
	DescriptionMd5 string   `deb822:"Description-md5"`
	Filename       string
	Size           int
	MD5sum         string
	SHA256         string
	Architecture   string
}

func (p BinaryIndexItem) String() string {
	return stringspp.UniversalStringer(p)
}

func NewBinaryIndex(reader io.Reader, inErr error) (*BinaryIndex, error) {
	if inErr != nil {
		return nil, inErr
	}
	var res BinaryIndex

	err := deb822.NewDecoder(reader).Decode(&res.Items)

	if s, ok := reader.(io.Closer); ok {
		s.Close()
	}

	if err != nil {
		return nil, err
	}

	res.byName = make(map[string]int, len(res.Items))
	res.bySourceName = make(map[string]int, len(res.Items))
	for idx, pkg := range res.Items {
		res.byName[pkg.Name] = idx
		// res.bySourceName[pkg.Source] = idx
	}

	return &res, nil
}

func (bi BinaryIndex) FindByName(name string) (*BinaryIndexItem, bool) {

	if idx, found := bi.byName[name]; found {
		return &bi.Items[idx], true
	}

	return nil, false
}

func (bi BinaryIndex) FindBySourceName(name string) (*BinaryIndexItem, bool) {
	if idx, found := bi.bySourceName[name]; found {
		return &bi.Items[idx], true
	}

	return nil, false
}
