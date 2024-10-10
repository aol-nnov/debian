package changelog

import "fmt"

const (
	BuildRefTag = "BuildRef"
	SrcRefTag   = "SrcRef"
)

// known tags in canonical order
var tagsOrder = []string{
	SrcRefTag,
	BuildRefTag,
}

type Tags map[string]string

func (tags Tags) String() (res string) {
	if len(tags) == 0 {
		return ""
	}

	for _, name := range tagsOrder {
		if tags[name] != "" {
			res += fmt.Sprintf("\n  %s: %s", name, tags[name])
		}
	}

	res += "\n"
	return
}
