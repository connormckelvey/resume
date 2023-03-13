package main

import (
	"flag"
	"log"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func main() {

	var (
		htmlFilePath  string
		pdfOutputPath string
	)

	flag.StringVar(&htmlFilePath, "html", "", "path/to/file.html")
	flag.StringVar(&pdfOutputPath, "output", "", "path/to/output.pdf")
	flag.Parse()

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatalf("error creating PDFGenerator: %v", err)
	}

	pdfg.Dpi.Set(300)
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	htmlFile, err := os.Open(htmlFilePath)
	if err != nil {
		log.Fatalf("error opening html file '%v': %v", htmlFilePath, err)
	}
	defer htmlFile.Close()

	htmlPage := wkhtmltopdf.NewPageReader(htmlFile)

	pdfg.AddPage(htmlPage)
	err = pdfg.Create()
	if err != nil {
		log.Fatalf("error generating pdf: %v", err)
	}

	err = pdfg.WriteFile(pdfOutputPath)
	if err != nil {
		log.Fatalf("error saving pdf: %v", err)
	}

	log.Printf("successfully wrote pdf to '%v'", pdfOutputPath)
}
