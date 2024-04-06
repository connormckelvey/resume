package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/yuin/goldmark"
)

var templateDir string
var outputPath string

func main() {
	flag.StringVar(&templateDir, "t", "", "path/to/template")
	flag.StringVar(&outputPath, "o", "", "path/to/output")
	flag.Parse()

	if templateDir == "" {
		log.Fatal("template dir required")
	}
	if outputPath == "" {
		log.Fatal("output path required")
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			newDocx(os.DirFS(templateDir)),
		),
	)
	src, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := output.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := md.Convert(src, output); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote to: %s", outputPath)
}
