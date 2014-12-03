package extra

import (
	"regexp"
)

var (
	rexMeta = regexp.MustCompile(`meta\.([a-zA-Z]*)`)
	rexData = regexp.MustCompile(`data\.([a-zA-Z]*)`)
)
