package topbar

import (
	"github.com/ambientsound/pms/api"
)

// Tag draws song tags from the currently playing song.
type Tag struct {
	api api.API
	tag string
}

// NewTag returns Tag.
func NewTag(a api.API, param string) Fragment {
	return &Tag{a, param}
}

// Text implements Fragment.
func (w *Tag) Text() (string, string) {
	song := w.api.Song()
	if song == nil {
		return `<none>`, `tagMissing`
	}
	if text, ok := song.StringTags[w.tag]; ok {
		return text, w.tag
	}
	return `<unknown>`, `tagMissing`
}
