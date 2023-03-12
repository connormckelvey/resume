package writer

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/connormckelvey/resume/pkg/schema"
)

type ResumeWriter struct {
	Funcs    template.FuncMap
	Styles   string
	Template *template.Template
}

type ResumeWriterOption func(*ResumeWriter) error

func Funcs(funcs template.FuncMap) ResumeWriterOption {
	return func(r *ResumeWriter) error {
		r.Funcs = funcs
		return nil
	}
}

func LoadTemplate(path string, funcs template.FuncMap) (*template.Template, error) {
	name := filepath.Base(path)
	template, err := template.New(name).Funcs(funcs).ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return template, nil
}

func Styles(path string) ResumeWriterOption {
	return func(r *ResumeWriter) error {
		r.Styles = path
		return nil
	}
}

func Template(path string) ResumeWriterOption {
	return func(r *ResumeWriter) (err error) {
		r.Template, err = LoadTemplate(path, r.Funcs)
		return err
	}
}

func NewResumeWriter(options ...ResumeWriterOption) (*ResumeWriter, error) {
	resume := &ResumeWriter{}
	for _, apply := range options {
		err := apply(resume)
		if err != nil {
			return nil, err
		}
	}
	return resume, nil
}

type ResumeInput struct {
	Styles string
	Resume *schema.ResumeSpec
}

func (r *ResumeWriter) Execute(w io.Writer, resume *schema.ResumeSpec) error {
	return r.Template.Execute(w, &ResumeInput{
		Styles: r.Styles,
		Resume: resume,
	})
}
