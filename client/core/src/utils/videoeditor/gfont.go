package videoeditor

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nicks344/moneytube/client/core/src/paths"
	"github.com/Nicks344/moneytube/client/core/src/utils"

	"github.com/imroc/req"
)

func GetGoogleFontPath(name string, bold, italic bool) (string, error) {
	fontPath := filepath.Join(paths.Fonts, name)

	if !utils.IsFileExists(fontPath) {
		resp, err := req.Get("https://fonts.google.com/download", req.QueryParam{"family": name})
		if err != nil {
			return "", err
		}

		tempFile := filepath.Join(paths.Temp, name+".zip")
		err = resp.ToFile(tempFile)
		if err != nil {
			return "", err
		}
		defer os.Remove(tempFile)

		err = unzip(tempFile, fontPath)
		if err != nil {
			return "", err
		}
	}

	fontFile := filepath.Join(fontPath, strings.ReplaceAll(name, " ", "")+"-")

	fontStyle := ""

	if bold && utils.IsFileExists(fontFile+"Bold.ttf") {
		fontStyle += "Bold"
	}

	if italic && utils.IsFileExists(fontFile+"Italic.ttf") {
		fontStyle += "Italic"
	}

	if fontStyle == "" {
		fontStyle += "Regular"
	}

	fontFile += fontStyle + ".ttf"

	if !utils.IsFileExists(fontFile) {
		fontFile = strings.Replace(fontFile, ".ttf", ".otf", 1)
	}

	if !utils.IsFileExists(fontFile) {
		return "", errors.New("font does not exist")
	}

	return filepath.Abs(fontFile)
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
