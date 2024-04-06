package main

import (
	"archive/zip"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/spf13/afero"
)

type docxTemplate struct {
	dirFS fs.FS
}

func newDocxTemplate(dirFS fs.FS) *docxTemplate {
	return &docxTemplate{
		dirFS: dirFS,
	}
}

func (dfs *docxTemplate) Render(w io.Writer, overrides map[string]io.Reader) error {
	template, err := newFsTemplate(dfs.dirFS)
	if err != nil {
		return err
	}

	for path, f := range overrides {
		path = filepath.Clean(path)
		dir := filepath.Dir(path)
		err := template.MkdirAll(dir, 0777)
		if err != nil {
			return err
		}
		df, err := template.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(df, f)
		if err != nil {
			return err
		}
		err = df.Close()
		if err != nil {
			return err
		}
	}

	zw := zip.NewWriter(w)
	defer zw.Close()

	structure := afero.NewIOFS(template)
	fs.WalkDir(structure, ".", func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := structure.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := zw.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	})

	return nil
}

func newFsTemplate(dir fs.FS) (afero.Fs, error) {
	mfs := afero.NewMemMapFs()

	err := fs.WalkDir(dir, ".", func(path string, info fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}
		if err != nil {
			return err
		}
		if info.IsDir() {
			if err := mfs.MkdirAll(path, 0755); err != nil {
				return err
			}
			return nil
		}
		sf, err := dir.Open(path)
		if err != nil {
			return err
		}
		defer sf.Close()

		f, err := mfs.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, sf)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return mfs, nil
}
