package bugreport

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nicks344/moneytube/client/core/src/modules/tools"
	"github.com/Nicks344/moneytube/client/core/src/paths"
)

type AEErrorData struct {
	ErrorData
	tools.VideoGeneratorData

	DataFile string `json:"dataFile"`
}

func getAEReport(desc string, dataJSON string) (archive []byte, err error) {
	var data AEErrorData
	err = json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		err = errors.New("Ошибка при парсинге данных")
	}

	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	addFileToZip(data.DataFile, "", w)

	for _, row := range data.Data {
		for _, layerData := range row {
			switch layerData.Type {
			case "data", "background-color":
				break

			default:
				addFileToZip(layerData.Data, "data", w)
			}
		}
	}

	temp := findLastPathWithProject(data.AepFile)

	if temp != "" {
		filepath.Walk(temp, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if temp == p {
				return nil
			}

			sizeMb := float32(info.Size()) / 1e+6
			if sizeMb > 100 {
				return nil
			}

			addFileToZip(p, "", w)

			return nil
		})

		logFile := filepath.Join(paths.Temp, fmt.Sprintf("aerender-%s.log", filepath.Base(temp)))
		addFileToZip(logFile, "", w)
	} else {
		addFileToZip(data.AepFile, "", w)
	}

	info := fmt.Sprintf(`Composition: %s

Get from queue: %t
Resolution: %s
Extention AE: %s

Compress: %t
Extention: %s
Quality: %d

Error text: %s
Description: %s`,
		data.CompositionName,
		data.FromRenderQueue,
		data.Resolution,
		data.Ext,
		data.IsCompress,
		data.OutputExt,
		data.Quality,
		data.Error,
		desc)

	f, err := w.Create("info.txt")
	if err != nil {
		return nil, err
	}
	_, err = f.Write([]byte(info))
	if err != nil {
		return
	}

	err = w.Close()
	if err != nil {
		return
	}

	archive = buf.Bytes()

	return
}

func findLastPathWithProject(aepFile string) string {
	type fileInfo struct {
		path string
		info os.FileInfo
	}

	path := paths.Temp
	name := filepath.Base(aepFile)

	var aepFiles []fileInfo

	filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(p, name) {
			aepFiles = append(aepFiles, fileInfo{
				path: p,
				info: info,
			})
		}

		return nil
	})

	if len(aepFiles) == 0 {
		return ""
	}

	sort.Slice(aepFiles, func(i, j int) bool {
		return aepFiles[i].info.ModTime().After(aepFiles[j].info.ModTime())
	})

	return strings.ReplaceAll(aepFiles[0].path, `\`+name, "")
}

func addFileToZip(path string, zipPath string, w *zip.Writer) error {
	if zipPath != "" {
		zipPath = zipPath + `\` + filepath.Base(path)
	} else {
		zipPath = filepath.Base(path)
	}

	f, err := w.Create(zipPath)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}
