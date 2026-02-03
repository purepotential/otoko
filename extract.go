package main

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extractAlbum(src *os.File, dir string) error {
	s, err := src.Stat()
	if err != nil {
		return err
	}
	r, err := zip.NewReader(src, s.Size())
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			return errors.New("unexpected directory")
		}
		track := strings.Split(f.Name, " - ")
		filename := track[len(track)-1]
		
		// Sanityzuj Unicode - zostaw tylko ASCII, cyfry, spacje, myślniki, kropki
		filename = sanitizeUnicode(filename)
		
		name := filepath.Join(dir, filename)
		if err := unzipFile(f, name); err != nil {
			return err
		}
	}
	return nil
}

func sanitizeUnicode(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || 
		   (r >= '0' && r <= '9') || r == ' ' || r == '-' || r == '.' || r == '_' {
			return r
		}
		return '_'
	}, s)
}

func unzipFile(src *zip.File, name string) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, src.Mode())
	if err != nil {
		return err
	}
	defer f.Close()

	z, err := src.Open()
	if err != nil {
		return err
	}
	defer z.Close()

	if _, err := io.Copy(f, z); err != nil {
		return err
	}

	return nil
}
