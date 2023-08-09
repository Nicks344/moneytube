package tools

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/modules/macros"
	"github.com/Nicks344/moneytube/client/core/src/paths"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/client/core/src/utils/videoeditor"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/google/uuid"

	"github.com/gabriel-vasile/mimetype"
)

type CopiesGeneratorWorker struct {
	*ToolService

	Input            string
	ResultPath       string
	Count            int
	ImageOverlay     utils.ImageOverlayData
	TextOverlay      TextOverlayData
	Intro            string
	IntroDurationSec int
	Outro            string
	OutroDurationSec int
	CutMethod        int
	CutSecondsFrom   float64
	CutSecondsTo     float64
	CutParts         int
}

type TextOverlayData struct {
	Text            moneytubemodel.UploadOptions
	Enabled         bool
	X               int
	Y               int
	From            int
	To              int
	Background      bool
	Color           string
	BackgroundColor string
	Font            string
	Size            int
	Bold            bool
	Italic          bool
}

func (w *CopiesGeneratorWorker) Start() {
	w.maxProgress = w.Count
	w.handleProgressChanged("")

	rand.Seed(time.Now().UnixNano())

	isDir, err := utils.IsDir(w.Input)
	if err != nil {
		w.handleError(err)
		return
	}

	if isDir {
		files, err := utils.GetFilesFromPath(w.Input, []string{"video"})
		if err != nil {
			w.handleError(err)
			return
		}

		w.maxProgress = w.Count * len(files)
		w.handleProgressChanged("")

		for _, filename := range files {
			err = w.generate(filename)

			if w.isCancelled(true) {
				return
			}

			if err != nil {
				w.handleError(err)
				return
			}
		}
	} else {
		err = w.generate(w.Input)
		if w.isCancelled(true) {
			return
		}
		if err != nil {
			w.handleError(err)
			return
		}
	}

}

func (w *CopiesGeneratorWorker) generate(filename string) error {
	for i := 0; i < w.Count; i++ {
		if w.isCancelled(true) {
			return errors.New("cancelled")
		}
		name := filepath.Base(filename)
		err := w.generateFile(filename, filepath.Join(w.ResultPath, fmt.Sprintf("%d_%s", i, name)))
		if err != nil {
			return err
		}
		w.progress++
		w.handleProgressChanged("")
	}
	return nil
}

func getOverlayImage(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	isDir, err := utils.IsDir(input)
	if err != nil {
		return "", err
	}

	if isDir {
		files, err := utils.GetFilesFromPath(input, []string{"video", "image"})
		if err != nil {
			return "", err
		}
		if len(files) == 0 {
			return "", nil
		}
		input = files[utils.RandRange(0, len(files)-1)]
	}

	return input, nil
}

func (w *CopiesGeneratorWorker) generateFile(filename string, newFilename string) error {
	name := strings.TrimSuffix(filepath.Base(filename), path.Ext(filename))
	tempfile1 := filepath.Join(paths.Temp, fmt.Sprintf("%s_%d_1%s", name, time.Now().Nanosecond(), path.Ext(newFilename)))
	defer os.Remove(tempfile1)
	tempfile2 := filepath.Join(paths.Temp, fmt.Sprintf("%s_%d_2%s", name, time.Now().Nanosecond(), path.Ext(newFilename)))
	defer os.Remove(tempfile2)

	if err := w.executeOverlay(filename, tempfile1); err != nil {
		if err.Error() == "no overlays" {
			tempfile1 = filename
		} else {
			return err
		}
	}

	cutSeconds := float64(utils.RandRange(int(w.CutSecondsFrom*1000), int(w.CutSecondsTo*1000))) / 1000

	if err := videoeditor.CutFragments(tempfile1, videoeditor.Cutmethod(w.CutMethod), cutSeconds, w.CutParts, tempfile2); err != nil {
		if err.Error() == "no cutting" {
			tempfile2 = tempfile1
		} else {
			return err
		}
	}

	if err := w.executeConcat(tempfile2, newFilename); err != nil {
		if err.Error() == "no intro or outro" {
			if err := copyFile(tempfile2, newFilename); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	resFile, err := os.OpenFile(newFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer resFile.Close()

	r := make([]byte, 64)
	rand.Read(r)
	_, err = resFile.Write(r)
	return err
}

func (w *CopiesGeneratorWorker) executeConcat(filename string, result string) error {
	if w.Intro == "" && w.Outro == "" {
		return errors.New("no intro or outro")
	}

	info, err := videoeditor.GetVideoInfo(filename)
	if err != nil {
		return err
	}

	if info.AudioCodec == "" {
		info.AudioCodec = "libmp3lame"
	}

	var ready int
	var errs []error
	var l sync.Mutex
	var intro, outro, video string

	go func() {
		defer func() {
			ready++
		}()

		var err error
		intro, err = getIntroOrOutro(w.Intro, info, w.IntroDurationSec*1000, path.Ext(result))
		if err != nil {
			l.Lock()
			defer l.Unlock()
			errs = append(errs, err)
		}
	}()

	go func() {
		defer func() {
			ready++
		}()

		if info.VideoCodec != "h264" {
			ext := path.Ext(filename)
			name := strings.TrimSuffix(filename, ext)
			video = fmt.Sprintf("%s/%s_%d%s%s", paths.Temp, name, time.Now().Nanosecond(), uuid.New().String(), ext)

			err := videoeditor.ConvertVideo(filename, video)
			if err != nil {
				l.Lock()
				defer l.Unlock()
				errs = append(errs, err)
			}
		} else {
			video = filename
		}
	}()

	go func() {
		defer func() {
			ready++
		}()

		var err error
		outro, err = getIntroOrOutro(w.Outro, info, w.OutroDurationSec*1000, path.Ext(result))
		if err != nil {
			l.Lock()
			defer l.Unlock()
			errs = append(errs, err)
		}
	}()

	for ready != 3 {
		time.Sleep(10 * time.Millisecond)
	}

	if len(errs) != 0 {
		errStr := ""
		for _, e := range errs {
			errStr += e.Error() + "\r\n"
		}

		return errors.New(errStr)
	}

	var files []string
	defer func() {
		for _, f := range files {
			if f != filename {
				os.Remove(f)
			}
		}
	}()

	if intro != "" {
		defer os.Remove(intro)
		files = append(files, intro)
	}
	files = append(files, video)
	if outro != "" {
		defer os.Remove(outro)
		files = append(files, outro)
	}

	return videoeditor.ConcatVideo(files, result)
}

func (w *CopiesGeneratorWorker) executeOverlay(filename string, result string) error {
	conf := videoeditor.OverlayConfig{
		Name:      "Copies",
		VideoFile: filename,
	}

	if w.TextOverlay.Enabled {
		w.TextOverlay.Text.ClearFromEmpty()

		text, err := w.TextOverlay.Text.GetOne()
		if err != nil {
			return err
		}
		text = macros.Execute(text, macros.StaticMacroses{})
		conf.Overlays = append(conf.Overlays, videoeditor.TextOverlayOpts{
			Bold:   w.TextOverlay.Bold,
			Italic: w.TextOverlay.Italic,
			Color:  strings.Replace(w.TextOverlay.Color, "#", "", 1),
			Position: videoeditor.Point{
				X: w.TextOverlay.X,
				Y: w.TextOverlay.Y,
			},
			Time: videoeditor.TimeRange{
				From: w.TextOverlay.From,
				To:   w.TextOverlay.To,
			},
			Size:            w.TextOverlay.Size,
			Font:            w.TextOverlay.Font,
			Text:            text,
			Background:      w.TextOverlay.Background,
			BackgroundColor: strings.Replace(w.TextOverlay.BackgroundColor, "#", "", 1),
		})
	}

	if w.ImageOverlay.Enabled {
		overlay, err := getOverlayImage(w.ImageOverlay.OverlaySrc)
		if err != nil {
			return err
		}

		conf.Overlays = append(conf.Overlays, videoeditor.MediaOverlayOpts{
			File: overlay,
			Time: videoeditor.TimeRange{
				From: w.ImageOverlay.From,
				To:   w.ImageOverlay.To,
			},
			Position: videoeditor.PositionOpts{
				Type: "coordinate",
				Coordinate: videoeditor.Point{
					X: w.ImageOverlay.X,
					Y: w.ImageOverlay.Y,
				},
			},
			Color: videoeditor.ColorOpts{
				Brightness: 100,
			},
		})
	}

	if len(conf.Overlays) == 0 {
		return errors.New("no overlays")
	}

	return videoeditor.DoOverlay(conf, result)

}

func getIntroOrOutro(input string, info videoeditor.VideoInfo, duration int, ext string) (string, error) {
	if input == "" {
		return "", nil
	}

	isDir, err := utils.IsDir(input)
	if err != nil {
		return "", err
	}

	if isDir {
		files, err := utils.GetFilesFromPath(input, []string{"video", "image"})
		if err != nil {
			return "", err
		}
		if len(files) == 0 {
			return "", nil
		}
		input = files[utils.RandRange(0, len(files)-1)]
	}

	filename := filepath.Base(input)
	name := strings.TrimSuffix(filename, path.Ext(filename))
	tempfile := filepath.Join(paths.Temp, fmt.Sprintf("%s_%d%s%s", name, time.Now().Nanosecond(), uuid.New().String(), ext))

	mime, err := mimetype.DetectFile(input)
	if err != nil {
		return "", err
	}

	if strings.Contains(mime.String(), "video") {
		if err := videoeditor.ScaleVideo(input, info, tempfile); err != nil {
			return "", err
		}

		return tempfile, nil

	} else if strings.Contains(mime.String(), "image") {
		if err := videoeditor.ConvertImageToVideo(input, info, duration, tempfile); err != nil {
			return "", err
		}

		return tempfile, nil
	}

	return "", nil
}

func copyFile(from, to string) error {
	toFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFile.Close()

	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	_, err = io.Copy(toFile, fromFile)
	return err
}
