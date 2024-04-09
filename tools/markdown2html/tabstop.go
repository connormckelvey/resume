package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var tabStopRegexp = regexp.MustCompile(`^\[(?:(\s\.)|(\.\s))\]\s*`)

type tabStopParser struct{}

func newTabStopParser() parser.InlineParser {
	return &tabStopParser{}
}

func (s *tabStopParser) Trigger() []byte {
	return []byte{'['}
}

func (s *tabStopParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	if parent.Parent() == nil {
		return nil
	}

	line, _ := block.PeekLine()
	if _, ok := parent.(*ast.Paragraph); !ok {
		return nil
	}

	m := tabStopRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}

	value := line[m[2]:m[3]]
	alignRight := bytes.HasSuffix(value, []byte("."))

	block.Advance(m[1])

	return newTabStop(alignRight)
}

func (s *tabStopParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

// A tabStop struct represents a tabstop.
type tabStop struct {
	ast.BaseInline
	alignRight bool
}

// Dump implements Node.Dump.
func (n *tabStop) Dump(source []byte, level int) {
	m := map[string]string{
		"AlignRight": fmt.Sprintf("%v", n.alignRight),
	}
	ast.DumpHelper(n, source, level, m, nil)
}

var kindTabStop = ast.NewNodeKind("TabStop")

// Kind implements Node.Kind.
func (n *tabStop) Kind() ast.NodeKind {
	return kindTabStop
}

// newTabStop returns a new TabStop node.
func newTabStop(alignRight bool) *tabStop {
	return &tabStop{
		alignRight: alignRight,
	}
}

// TabStopHTMLRenderer is a renderer.NodeRenderer implementation that
// renders checkboxes in list items.
type TabStopHTMLRenderer struct {
	html.Config
}

// NewTabStopHTMLRenderer returns a new TabStopHTMLRenderer.
func NewTabStopHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &TabStopHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *TabStopHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(kindTabStop, r.renderTabStop)
}

func (r *TabStopHTMLRenderer) renderTabStop(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*tabStop)

	var err error
	if n.alignRight {
		_, err = w.WriteString(`<span class="tabStop tabStop-alignRight"></span>`)

	} else {
		_, err = w.WriteString(`<span class="tabStop"></span>`)
	}
	if err != nil {
		return ast.WalkStop, err
	}
	return ast.WalkContinue, nil
}
