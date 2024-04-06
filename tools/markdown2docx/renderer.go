package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type docxRenderer struct {
	template *docxTemplate
	links    map[string]string
}

func newDocxRenderer(templateDir fs.FS) *docxRenderer {
	return &docxRenderer{
		template: newDocxTemplate(templateDir),
		links:    make(map[string]string),
	}
}

func (r *docxRenderer) Render(w io.Writer, source []byte, doc ast.Node) error {
	var document bytes.Buffer
	if err := r.renderDocument(&document, source, doc); err != nil {
		return err
	}

	var relationships bytes.Buffer
	if err := r.renderRelationships(&relationships); err != nil {
		return err
	}

	err := r.template.Render(w, map[string]io.Reader{
		"word/document.xml":            &document,
		"word/_rels/document.xml.rels": &relationships,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r docxRenderer) renderRelationships(w io.Writer) error {
	if _, err := fmt.Fprint(w, relsHeader); err != nil {
		return err
	}
	for url, rid := range r.links {
		if _, err := fmt.Fprintf(w, relTemplate, rid, url); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, relsFooter); err != nil {
		return err
	}

	return nil
}

func (r *docxRenderer) renderDocument(w io.Writer, source []byte, doc ast.Node) error {
	_, err := w.Write(documentHeader)
	if err != nil {
		return err
	}

	for c := doc.FirstChild(); c != nil; c = c.NextSibling() {
		err := r.renderBlock(w, source, c)
		if err != nil {
			return err
		}
	}
	_, err = w.Write(documentFooter)
	return err
}

func (r *docxRenderer) renderBlock(w io.Writer, source []byte, node ast.Node) error {
	switch n := node.(type) {
	case *ast.ThematicBreak:
		err := r.renderThematicBreak(w, source, n)
		if err != nil {
			return err
		}
	case *ast.Paragraph:
		err := r.renderParagraph(w, source, n)
		if err != nil {
			return err
		}
	case *ast.Heading:
		err := r.renderHeading(w, source, n)
		if err != nil {
			return err
		}
	case *ast.List:
		err := r.renderList(w, source, n)
		if err != nil {
			return err
		}
	case *ast.ListItem:
		err := r.renderListItem(w, source, n)
		if err != nil {
			return err
		}
	default:
		// Not supported
		// case *ast.TextBlock:
		// case *ast.ThematicBreak:
		// case *ast.CodeBlock:
		// case *ast.HTMLBlock:
		// case *ast.FencedCodeBlock:
		// case *ast.Blockquote:
		return fmt.Errorf("render block: no renderer implemented: %s|%v", n.Kind().String(), n.Type())
	}
	return nil
}

func (r *docxRenderer) renderTabStop(w io.Writer, wpProps io.Writer, _ io.Writer, _ []byte, tn *tabStop) error {
	align := "left"
	if tn.alignRight {
		align = "right"
	}
	fmt.Fprintf(wpProps, `<w:tabs><w:tab w:val="%s" w:pos="9026" /></w:tabs>`, align)
	fmt.Fprint(w, `<w:tab />`)
	return nil
}

func (r *docxRenderer) renderParagraphWithTabStop(w io.Writer, source []byte, n *ast.Paragraph) error {
	var text bytes.Buffer
	var wpProps bytes.Buffer
	var wrProps bytes.Buffer
	fmt.Fprint(w, `<w:p>`)

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *tabStop:
			err := r.renderTabStop(&text, &wpProps, &wrProps, source, c)
			if err != nil {
				return err
			}
		default:
			err := r.renderInline(&text, &wpProps, &wrProps, source, child)
			if err != nil {
				return err
			}
		}
	}
	if wpProps.Len() > 0 {
		fmt.Fprint(w, `<w:pPr>`)
		io.Copy(w, &wpProps)
		fmt.Fprint(w, `</w:pPr>`)
	}
	fmt.Fprint(w, `<w:r>`)
	if wrProps.Len() > 0 {
		fmt.Fprint(w, `<w:rPr>`)
		io.Copy(w, &wrProps)
		fmt.Fprint(w, `</w:rPr>`)
	}
	io.Copy(w, &text)
	fmt.Fprint(w, `</w:r>`)
	fmt.Fprint(w, `</w:p>`)

	return nil
}

func (r *docxRenderer) renderThematicBreak(w io.Writer, _ []byte, _ *ast.ThematicBreak) error {
	_, err := fmt.Fprint(w, `<w:p><w:pPr><w:spacing w:line="50"/></w:pPr></w:p>`)
	return err
}

func (r *docxRenderer) renderParagraph(w io.Writer, source []byte, n *ast.Paragraph) error {
	var ts *tabStop
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if c, ok := child.(*tabStop); ok {
			ts = c
			break
		}
	}
	if ts != nil {
		return r.renderParagraphWithTabStop(w, source, n)
	}
	fmt.Fprint(w, `<w:p>`)
	var text bytes.Buffer
	var wpProps bytes.Buffer
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		err := r.renderInline(&text, &wpProps, nil, source, child)
		if err != nil {
			return err
		}
	}
	if wpProps.Len() > 0 {
		fmt.Fprint(w, `<w:pPr>`)
		io.Copy(w, &wpProps)
		fmt.Fprint(w, `</w:pPr>`)

	}
	io.Copy(w, &text)
	fmt.Fprint(w, `</w:p>`)
	return nil
}

func (r *docxRenderer) renderHeading(w io.Writer, source []byte, h *ast.Heading) error {
	var text bytes.Buffer
	var wpProps bytes.Buffer

	for _, attr := range h.Attributes() {
		fmt.Fprintf(&wpProps, `<w:%s w:val="%v" />`, string(attr.Name), string(attr.Value.([]byte)))
	}

	for child := h.FirstChild(); child != nil; child = child.NextSibling() {
		err := r.renderInline(&text, &wpProps, nil, source, child)
		if err != nil {
			return err
		}
	}
	fmt.Fprint(w, `<w:p>`)
	fmt.Fprintf(&wpProps, `<w:pStyle w:val="Heading%d" />`, h.Level)

	if wpProps.Len() > 0 {
		fmt.Fprint(w, `<w:pPr>`)
		io.Copy(w, &wpProps)
		fmt.Fprint(w, `</w:pPr>`)
	}
	io.Copy(w, &text)
	fmt.Fprint(w, `</w:p>`)

	return nil
}

func (r *docxRenderer) renderList(w io.Writer, source []byte, l *ast.List) error {
	for child := l.FirstChild(); child != nil; child = child.NextSibling() {
		err := r.renderBlock(w, source, child)
		if err != nil {
			return err
		}
	}
	return nil
}
func (r *docxRenderer) renderListItem(w io.Writer, source []byte, l *ast.ListItem) error {
	var text bytes.Buffer
	var wpProps bytes.Buffer
	fmt.Fprint(&wpProps, `<w:numPr>`)
	fmt.Fprint(&wpProps, `<w:pStyle w:val="ListParagraph" />`)
	fmt.Fprint(&wpProps, `<w:ilvl w:val="0" />`)
	fmt.Fprint(&wpProps, `<w:numId w:val="2" />`)
	fmt.Fprint(&wpProps, `</w:numPr>`)

	for child := l.FirstChild(); child != nil; child = child.NextSibling() {
		err := r.renderInline(&text, &wpProps, nil, source, child)
		if err != nil {
			return err
		}
	}

	fmt.Fprint(w, `<w:p>`)
	if wpProps.Len() > 0 {
		fmt.Fprint(w, `<w:pPr>`)
		io.Copy(w, &wpProps)
		fmt.Fprint(w, `</w:pPr>`)
	}
	io.Copy(w, &text)
	fmt.Fprint(w, `</w:p>`)
	return nil

}

func (r *docxRenderer) renderTextRun(w io.Writer, wpProps io.Writer, source []byte, n ast.Node, renderInside func(w io.Writer, wpProps io.Writer, wrProps io.Writer, source []byte, n ast.Node) error) error {
	var runPropsW bytes.Buffer
	var text bytes.Buffer
	fmt.Fprint(w, `<w:r>`)
	err := renderInside(&text, wpProps, &runPropsW, source, n)
	if err != nil {
		return err
	}
	if runPropsW.Len() > 0 {
		fmt.Fprint(w, `<w:rPr>`)
		io.Copy(w, &runPropsW)
		fmt.Fprint(w, `</w:rPr>`)
	}
	io.Copy(w, &text)
	fmt.Fprint(w, `</w:r>`)
	return nil
}

func (r *docxRenderer) renderLink(w io.Writer, wpProps, _ io.Writer, source []byte, ln *ast.Link) error {
	rid, ok := r.links[string(ln.Destination)]
	if !ok {
		rid = uniqueRID()
		r.links[string(ln.Destination)] = rid
	}
	fmt.Fprintf(w, `<w:hyperlink w:history="1" r:id="%s">`, rid)
	for child := ln.FirstChild(); child != nil; child = child.NextSibling() {
		fmt.Fprint(w, `<w:r>`)
		var wrProps bytes.Buffer

		fmt.Fprint(&wrProps, `<w:rStyle w:val="Hyperlink" />`)

		var text bytes.Buffer
		err := r.renderInline(&text, wpProps, &wrProps, source, child)
		if err != nil {
			return err
		}

		if wrProps.Len() > 0 {
			fmt.Fprint(w, `<w:rPr>`)
			io.Copy(w, &wrProps)
			fmt.Fprint(w, `</w:rPr>`)
		}
		io.Copy(w, &text)
		fmt.Fprint(w, `</w:r>`)
	}
	fmt.Fprint(w, `</w:hyperlink>`)
	return nil
}

func (r *docxRenderer) renderInline(w io.Writer, wpProps io.Writer, wrProps io.Writer, source []byte, n ast.Node) error {
	switch nn := n.(type) {
	case *ast.Link:
		return r.renderLink(w, wpProps, wrProps, source, nn)
	default:
	}

	// called from block level, not nested inline
	if wrProps == nil {
		return r.renderTextRun(w, wpProps, source, n, func(w, wpProps, wrProps io.Writer, source []byte, n ast.Node) error {
			return r.renderInline(w, wpProps, wrProps, source, n)
		})
	}

	// child from above, writing to existing rPr
	switch nn := n.(type) {
	case *tabStop:
		return r.renderTabStop(w, wpProps, wrProps, source, nn)
	case *ast.Text:
		return r.renderText(w, wrProps, source, nn)
	case *ast.TextBlock:
		return r.renderTextBlock(w, wrProps, source, nn)
	case *ast.Emphasis:
		err := r.renderEmphasis(w, wpProps, wrProps, source, nn)
		if err != nil {
			return err
		}
	default:
		// Not supported
		// case *ast.CodeSpan:
		// case *ast.Image:
		return fmt.Errorf("render inline: no renderer implemented: %s|%s|%v", nn.Parent().Kind().String(), nn.Kind().String(), nn.Type())
	}
	return nil
}

func (r *docxRenderer) renderText(w io.Writer, _ io.Writer, source []byte, n *ast.Text) error {
	text := escapeXMLBytes(n.Text(source))
	w.Write([]byte(`<w:t xml:space="preserve">`))
	w.Write(text)
	w.Write([]byte("</w:t>"))
	return nil
}

func (r *docxRenderer) renderTextBlock(w io.Writer, _ io.Writer, source []byte, n *ast.TextBlock) error {
	text := escapeXMLBytes(n.Text(source))
	w.Write([]byte(`<w:t xml:space="preserve">`))
	w.Write(text)
	w.Write([]byte("</w:t>"))
	return nil
}

func (r *docxRenderer) renderEmphasis(w, wpProps, wrProps io.Writer, source []byte, n *ast.Emphasis) error {
	tag := `<w:i /><w:iCs />`
	if n.Level == 2 {
		tag = `<w:b /><w:bCs />`
	}
	fmt.Fprint(wrProps, tag)

	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		err := r.renderInline(w, wpProps, wrProps, source, child)
		if err != nil {
			return err
		}
	}
	return nil
}

func escapeXMLBytes(src []byte) []byte {
	// Implement XML escaping for characters like <, >, &, etc.
	// This is a placeholder. Proper implementation needed.
	return src
}
func (r *docxRenderer) AddOptions(opts ...renderer.Option) {
}

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
