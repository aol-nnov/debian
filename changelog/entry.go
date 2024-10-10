package changelog

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/aol-nnov/debian/fields"
)

// https://manpages.debian.org/testing/dpkg-dev/deb-changelog.5.en.html

type Timestamp time.Time

func (ts Timestamp) String() string {
	return time.Time(ts).Format(time.RFC1123Z)
}

type ChangelogBody string

func (b ChangelogBody) String() (res string) {

	for _, line := range strings.Split(string(b), "\n") {
		res += "\n  " + line
	}

	res += "\n"
	return
}

// Debian changelog entry
type Entry struct {
	PackageName  string
	Version      fields.Version
	Distribution string
	Metadata     string
	body         ChangelogBody
	Tags         Tags
	Maintainer   fields.Maintainer
	Timestamp    Timestamp
}

func NewEntry() Entry {
	return Entry{
		Metadata: "urgency=medium",
		Maintainer: fields.Maintainer{
			Name:  os.Getenv("DEBFULLNAME"),
			Email: os.Getenv("DEBEMAIL"),
		},
		Timestamp: Timestamp(time.Now()),
	}
}

func NewEntryFromTemplate(e Entry) Entry {
	return Entry{
		PackageName:  e.PackageName,
		Version:      e.Version,
		Distribution: e.Distribution,
		Metadata:     "urgency=medium",
		Maintainer: fields.Maintainer{
			Name:  os.Getenv("DEBFULLNAME"),
			Email: os.Getenv("DEBEMAIL"),
		},
		Timestamp: Timestamp(time.Now()),
	}
}

func (e Entry) GetBody() string {
	return string(e.body)
}

func (e *Entry) SetBody(body string) {

	bodyWithoutTags := body

	// as the (#tagsOrder) is known, we search for first known tag. Everything before it is a changelog body
	// FIXME: only known tags are parsed
	if tagsStartIdx := strings.Index(body, SrcRefTag); tagsStartIdx != -1 {
		bodyWithoutTags = body[0:tagsStartIdx]

		for _, line := range strings.Split(body[tagsStartIdx:], "\n") {
			fields := strings.Fields(line)
			if len(fields) == 2 {
				tagName := strings.Trim(fields[0], " :")
				tagValue := fields[1]

				if slices.Contains(tagsOrder, tagName) {
					e.AddTag(tagName, tagValue)
				}

			}

		}
	}

	e.body = ChangelogBody(strings.Trim(bodyWithoutTags, "\n\t"))
}

func (e *Entry) AddTag(key, value string) {
	if e.Tags == nil {
		e.Tags = make(Tags)
	}

	e.Tags[key] = value
}

func (e *Entry) GetTag(key string) string {
	return e.Tags[key]
}

func (e Entry) String() (res string) {
	return fmt.Sprintf("%s (%s) %s; %s\n%s%s\n -- %s  %s\n\n",
		e.PackageName,
		e.Version,
		e.Distribution,
		e.Metadata,

		e.body,
		e.Tags,

		e.Maintainer,
		e.Timestamp,
	)
}
