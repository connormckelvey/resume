package main

import (
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func uniqueRID() string {
	rid, _ := gonanoid.New()
	rid = "rId" + strings.ToLower(rid)
	return rid
}

var relsHeader = `<?xml version="1.0" encoding="UTF-8"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
	<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/numbering" Target="numbering.xml"/>
	<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footnotes" Target="footnotes.xml"/>
	<Relationship Id="rId4" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/settings" Target="settings.xml"/>
	<Relationship Id="rId5" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/comments" Target="comments.xml"/>`

var relsFooter = `</Relationships>`

const relTemplate = `<Relationship Id="%s" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="%s" TargetMode="External"/>`

var documentHeader = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
	<w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas"
		xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006"
		xmlns:o="urn:schemas-microsoft-com:office:office"
		xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
		xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math"
		xmlns:v="urn:schemas-microsoft-com:vml"
		xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing"
		xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
		xmlns:w10="urn:schemas-microsoft-com:office:word"
		xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
		xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml"
		xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml"
		xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup"
		xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk"
		xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml"
		xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape"
		xmlns:cx="http://schemas.microsoft.com/office/drawing/2014/chartex"
		xmlns:cx1="http://schemas.microsoft.com/office/drawing/2015/9/8/chartex"
		xmlns:cx2="http://schemas.microsoft.com/office/drawing/2015/10/21/chartex"
		xmlns:cx3="http://schemas.microsoft.com/office/drawing/2016/5/9/chartex"
		xmlns:cx4="http://schemas.microsoft.com/office/drawing/2016/5/10/chartex"
		xmlns:cx5="http://schemas.microsoft.com/office/drawing/2016/5/11/chartex"
		xmlns:cx6="http://schemas.microsoft.com/office/drawing/2016/5/12/chartex"
		xmlns:cx7="http://schemas.microsoft.com/office/drawing/2016/5/13/chartex"
		xmlns:cx8="http://schemas.microsoft.com/office/drawing/2016/5/14/chartex"
		xmlns:aink="http://schemas.microsoft.com/office/drawing/2016/ink"
		xmlns:am3d="http://schemas.microsoft.com/office/drawing/2017/model3d"
		xmlns:w16cex="http://schemas.microsoft.com/office/word/2018/wordml/cex"
		xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid"
		xmlns:w16="http://schemas.microsoft.com/office/word/2018/wordml"
		xmlns:w16sdtdh="http://schemas.microsoft.com/office/word/2020/wordml/sdtdatahash"
		xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex"
		mc:Ignorable="w14 w15 wp14">
		<w:body>
`)

var documentFooter = []byte(`<w:sectPr>
<w:pgSz w:w="11906" w:h="16838" w:orient="portrait" />
<w:pgMar w:top="0.75in" w:right="1440" w:bottom="0.75in" w:left="1440" w:header="708"
	w:footer="708" w:gutter="0" />
<w:pgNumType />
<w:docGrid w:linePitch="360" />
</w:sectPr></w:body></w:document>`)
