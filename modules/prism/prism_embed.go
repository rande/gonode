package prism

import (
	"embed"
)

//go:embed all:templates
var content embed.FS

func GetEmbedFS() embed.FS {
	return content
}
