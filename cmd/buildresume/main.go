package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/connormckelvey/resume/pkg/schema"
	"github.com/connormckelvey/resume/pkg/writer"
)

func main() {
	//
	var (
		resumePath   string
		templatePath string
		stylePath    string
		outputPath   string
	)

	flag.StringVar(&resumePath, "resume", "", "/path/to/resume.yml")
	flag.StringVar(&templatePath, "template", "", "/path/to/template.html")
	flag.StringVar(&stylePath, "styles", "", "/path/to/styles.css")
	flag.StringVar(&outputPath, "output", "", "/path/to/output.html")
	flag.Parse()

	resume, err := schema.LoadResume(resumePath)
	if err != nil {
		log.Fatalf("error loading resume: %v", err)
	}

	resumeWriter, err := writer.NewResumeWriter(
		writer.Funcs(writer.Utils),
		writer.Template(templatePath),
		writer.Styles(stylePath),
	)
	if err != nil {
		log.Fatalf("error creating resumeWriter '%v': %v", templatePath, err)
	}

	outputDir := filepath.Dir(outputPath)
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("error creating output directory '%v': %v", outputDir, err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("error creating output file '%v': %v", outputPath, err)
	}
	defer f.Close()

	err = resumeWriter.Execute(f, resume)
	if err != nil {
		log.Fatalf("error writing to output file '%v': %v", outputPath, err)
	}

	log.Printf("successfully wrote resume to '%v'", outputPath)
}
