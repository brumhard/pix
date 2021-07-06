package assets

import "embed"

//go:embed *.gohtml
var HTMLTemplates embed.FS
