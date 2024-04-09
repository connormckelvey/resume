package main

import (
	"io/fs"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

type docxExtension struct {
	templateDir fs.FS
}

func newDocx(templateDir fs.FS) *docxExtension {
	return &docxExtension{
		templateDir: templateDir,
	}
}

func (e *docxExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(newTabStopParser(), 0),
		),
		parser.WithAttribute(),
	)
	m.SetRenderer(newDocxRenderer(e.templateDir))
}
