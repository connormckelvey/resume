package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"

	"github.com/yuin/goldmark/renderer"
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
	case *east.DefinitionList:
		err := r.renderDefinitionList(w, source, n)
		if err != nil {
			return err
		}
	case *east.DefinitionTerm:
		err := r.renderDefinitionTerm(w, nil, nil, source, n)
		if err != nil {
			return err
		}
	case *east.DefinitionDescription:
		err := r.renderDefinitionDescription(w, nil, nil, source, n)
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
	fmt.Fprintf(wpProps, `<w:tabs><w:tab w:val="%s" w:pos="9746" /></w:tabs>`, align)
	fmt.Fprint(w, `<w:tab />`)
	return nil
}

func (r *docxRenderer) renderParagraphWithTabStop(w io.Writer, source []byte, n *ast.Paragraph) error {
	return wParagraph(w, func(wpProps, content io.Writer) error {
		return wRun(content, func(wrProps, content io.Writer) error {
			for child := n.FirstChild(); child != nil; child = child.NextSibling() {
				switch c := child.(type) {
				case *tabStop:
					err := r.renderTabStop(content, wpProps, wrProps, source, c)
					if err != nil {
						return err
					}
				default:
					err := r.renderInline(content, wpProps, wrProps, source, child)
					if err != nil {
						return err
					}
				}
			}
			return nil
		})
	})
}

const spacing = `<w:p><w:pPr><w:spacing w:line="%d"/></w:pPr></w:p>`

func (r *docxRenderer) renderThematicBreak(w io.Writer, _ []byte, _ *ast.ThematicBreak) error {
	_, err := fmt.Fprintf(w, spacing, 50)
	return err
}

func (r *docxRenderer) renderParagraph(w io.Writer, source []byte, n *ast.Paragraph) error {
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		if _, ok := child.(*tabStop); ok {
			return r.renderParagraphWithTabStop(w, source, n)
		}
	}
	return wParagraph(w, func(wpProps, content io.Writer) error {
		for child := n.FirstChild(); child != nil; child = child.NextSibling() {
			err := r.renderInline(content, wpProps, nil, source, child)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func wParagraph(w io.Writer, handler func(wpProps, content io.Writer) error) error {
	var (
		wpProps bytes.Buffer
		content bytes.Buffer
	)
	if err := handler(&wpProps, &content); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, `<w:p>`); err != nil {
		return err
	}
	if wpProps.Len() > 0 {
		if _, err := fmt.Fprintf(w, `<w:pPr>%s</w:pPr>`, wpProps.String()); err != nil {
			return err
		}
	}
	if _, err := io.Copy(w, &content); err != nil {
		return err
	}

	if _, err := fmt.Fprint(w, `</w:p>`); err != nil {
		return err
	}
	return nil
}

const headingStyle = `<w:pStyle w:val="Heading%d" />`

func (r *docxRenderer) renderHeading(w io.Writer, source []byte, h *ast.Heading) error {
	return wParagraph(w, func(wpProps, content io.Writer) error {
		if _, err := fmt.Fprintf(wpProps, headingStyle, h.Level); err != nil {
			return err
		}
		for _, attr := range h.Attributes() {
			if _, err := fmt.Fprintf(wpProps, `<w:%s w:val="%v" />`, string(attr.Name), string(attr.Value.([]byte))); err != nil {
				return err
			}
		}
		for child := h.FirstChild(); child != nil; child = child.NextSibling() {
			err := r.renderInline(content, wpProps, nil, source, child)
			if err != nil {
				return err
			}
		}
		return nil
	})
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

func (r *docxRenderer) renderDefinitionList(w io.Writer, source []byte, dl *east.DefinitionList) error {
	for child := dl.FirstChild(); child != nil && child.NextSibling() != nil; child = child.NextSibling().NextSibling() {
		term := child
		def := child.NextSibling()

		err := wParagraph(w, func(wpProps, content io.Writer) error {
			err := r.renderDefinitionTerm(content, wpProps, nil, source, term.(*east.DefinitionTerm))
			if err != nil {
				return err
			}

			err = wRun(content, func(wrProps, content io.Writer) error {
				return r.renderText(content, wrProps, []byte(": "), &ast.Text{
					Segment: text.Segment{
						Start:   0,
						Stop:    2,
						Padding: 0,
					},
				})
			})
			if err != nil {
				return err
			}
			err = r.renderDefinitionDescription(content, wpProps, nil, source, def.(*east.DefinitionDescription))
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

	}
	return nil
}

func (r *docxRenderer) renderDefinitionTerm(content, wpProps, _ io.Writer, source []byte, dt *east.DefinitionTerm) error {
	for child := dt.FirstChild(); child != nil; child = child.NextSibling() {
		err := r.renderInline(content, wpProps, nil, source, child)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *docxRenderer) renderDefinitionDescription(content, wpProps, _ io.Writer, source []byte, dt *east.DefinitionDescription) error {
	return wRun(content, func(wrProps, content io.Writer) error {
		for child := dt.FirstChild(); child != nil; child = child.NextSibling() {
			err := r.renderInline(content, wpProps, wrProps, source, child)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

const listItemStyle = `<w:numPr><w:pStyle w:val="ListParagraph" /><w:ilvl w:val="0" /><w:numId w:val="2" /></w:numPr>`

func (r *docxRenderer) renderListItem(w io.Writer, source []byte, l *ast.ListItem) error {
	return wParagraph(w, func(wpProps, content io.Writer) error {
		if _, err := fmt.Fprint(wpProps, listItemStyle); err != nil {
			return err
		}
		for child := l.FirstChild(); child != nil; child = child.NextSibling() {
			err := r.renderInline(content, wpProps, nil, source, child)
			if err != nil {
				return err
			}
		}
		return nil

	})
}

func wRun(w io.Writer, handler func(wrProps io.Writer, content io.Writer) error) error {
	if _, err := fmt.Fprint(w, `<w:r>`); err != nil {
		return err
	}
	var wrProps bytes.Buffer
	var content bytes.Buffer
	if err := handler(&wrProps, &content); err != nil {
		return err
	}
	if wrProps.Len() > 0 {
		if _, err := fmt.Fprintf(w, `<w:rPr>%s</w:rPr>`, wrProps.String()); err != nil {
			return err
		}
	}
	if _, err := io.Copy(w, &content); err != nil {
		return err
	}
	if _, err := fmt.Fprint(w, `</w:r>`); err != nil {
		return err
	}
	return nil
}

const hyperlinkStyle = `<w:rStyle w:val="Hyperlink" />`

func (r *docxRenderer) renderLink(w io.Writer, wpProps, _ io.Writer, source []byte, ln *ast.Link) error {
	rid, ok := r.links[string(ln.Destination)]
	if !ok {
		rid = uniqueRID()
		r.links[string(ln.Destination)] = rid
	}
	fmt.Fprintf(w, `<w:hyperlink w:history="1" r:id="%s">`, rid)
	for child := ln.FirstChild(); child != nil; child = child.NextSibling() {
		err := wRun(w, func(runProps io.Writer, content io.Writer) error {
			fmt.Fprint(runProps, hyperlinkStyle)
			return r.renderInline(content, wpProps, runProps, source, child)
		})
		if err != nil {
			return err
		}
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
		return wRun(w, func(wrProps, content io.Writer) error {
			return r.renderInline(content, wpProps, wrProps, source, n)
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
	case *east.DefinitionTerm:
		err := r.renderDefinitionTerm(w, wpProps, wrProps, source, nn)
		if err != nil {
			return err
		}
	case *east.DefinitionDescription:
		err := r.renderDefinitionDescription(w, wpProps, wrProps, source, nn)
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
	_, err := fmt.Fprintf(w, `<w:t xml:space="preserve">%s</w:t>`, string(text))
	return err
}

func (r *docxRenderer) renderTextBlock(w io.Writer, _ io.Writer, source []byte, n *ast.TextBlock) error {
	text := escapeXMLBytes(n.Text(source))
	_, err := fmt.Fprintf(w, `<w:t xml:space="preserve">%s</w:t>`, string(text))
	return err
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
