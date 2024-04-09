package main

import (
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// DefinitionListHTMLRenderer is a renderer.NodeRenderer implementation that
// renders DefinitionList nodes.
type DefinitionListHTMLRenderer struct {
	html.Config
}

// NewDefinitionListHTMLRenderer returns a new DefinitionListHTMLRenderer.
func NewDefinitionListHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &DefinitionListHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *DefinitionListHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindDefinitionList, r.renderDefinitionList)
	reg.Register(ast.KindDefinitionTerm, r.renderDefinitionTerm)
	reg.Register(ast.KindDefinitionDescription, r.renderDefinitionDescription)
}

func (r *DefinitionListHTMLRenderer) renderDefinitionList(
	w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<dl")
			html.RenderAttributes(w, n, extension.DefinitionListAttributeFilter)
			_, _ = w.WriteString(">\n")
		} else {
			_, _ = w.WriteString("<dl>\n")
		}
	} else {
		_, _ = w.WriteString("</dl>\n")
	}
	return gast.WalkContinue, nil
}

func (r *DefinitionListHTMLRenderer) renderDefinitionTerm(
	w util.BufWriter, source []byte, n gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString(`<div>`)
		if n.Attributes() != nil {
			_, _ = w.WriteString("<dt")
			html.RenderAttributes(w, n, extension.DefinitionTermAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<dt>")
		}
	} else {
		_, _ = w.WriteString(":</dt>\n")
	}
	return gast.WalkContinue, nil
}

func (r *DefinitionListHTMLRenderer) renderDefinitionDescription(
	w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if entering {
		n := node.(*ast.DefinitionDescription)
		_, _ = w.WriteString("<dd")
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, extension.DefinitionDescriptionAttributeFilter)
		}
		if n.IsTight {
			_, _ = w.WriteString(">")
		} else {
			_, _ = w.WriteString(">\n")
		}
	} else {
		_, _ = w.WriteString("</dd></div>\n")
	}
	return gast.WalkContinue, nil
}
