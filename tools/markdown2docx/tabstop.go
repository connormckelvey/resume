package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
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
