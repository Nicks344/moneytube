package tools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/paths"

	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/utils"

	"github.com/gabriel-vasile/mimetype"
)

const (
	ByIntervalDuration   = 10
	ByAudioDuration      = 20
	BySlideCountDuration = 30
)

type VideoFFmpegGeneratorWorker struct {
	*ToolService

	Input             string
	ResultPath        string
	Ext               string
	DeleteImages      bool
	AudioFile         string
	IfVideoLonger     int
	IfAudioLonger     int
	Intro             string
	Outro             string
	SlideDurationFrom int
	SlideDurationTo   int
	VideoDurationType int
	VideoDurationFrom int
	VideoDurationTo   int
	SlideCountFrom    int
	SlideCountTo      int
	VideoCount        int
	Threads           int
	FPS               int
	Overlay           utils.ImageOverlayData

	introFrame Frame
	outroFrame Frame

	wg sync.WaitGroup
}

type Frame struct {
	File      string
	Duration  int
	VideoFile string
}

func (w *VideoFFmpegGeneratorWorker) getIntroOrOutro(file string) (Frame, error) {
	frame := Frame{}

	if file == "" {
		return frame, nil
	}

	isDir, err := utils.IsDir(file)
	if err != nil {
		return frame, err
	}

	if isDir {
		files, err := utils.GetFilesFromPath(file, []string{"video", "image"})
		if err != nil {
			return frame, err
		}
		if len(files) == 0 {
			return frame, nil
		}
		file = files[utils.RandRange(0, len(files)-1)]
	}

	mime, err := mimetype.DetectFile(file)
	if err != nil {
		return frame, err
	}

	if strings.Contains(mime.String(), "video") {
		filename := filepath.Base(file)
		name := strings.TrimSuffix(filename, path.Ext(filename))
		tempfile := filepath.Join(paths.Temp, fmt.Sprintf("%s_%d%s", name, time.Now().Nanosecond(), w.Ext))
		dur, err := utils.GetVideoDuration(file)
		if err != nil {
			return frame, err
		}
		cmd := exec.Command(config.GetFFmpegBin(), "-y",
			"-i", file,
			"-c:v", "libx264",
			"-preset", "ultrafast",
			"-crf", "18",
			"-pix_fmt", "yuv420p",
			"-filter:v", fmt.Sprintf("fps=fps=%d,scale=1920:1080", w.FPS),
			tempfile,
		)
		var out []byte
		out, err = cmd.CombinedOutput()
		if err != nil {
			logger.Error(errors.New("error on generate frame, out: " + string(out)))
			return frame, err
		}
		frame.Duration = int(dur.Milliseconds())
		frame.VideoFile = tempfile
	} else if strings.Contains(mime.String(), "image") {
		frame.File = file
		frame.Duration = utils.RandRange(w.SlideDurationFrom, w.SlideDurationTo)
	}

	return frame, nil
}

func (w *VideoFFmpegGeneratorWorker) Start() {
	log.Println("start generate video")
	w.wg = sync.WaitGroup{}
	w.Threads = 10
	w.maxProgress = w.VideoCount
	w.handleProgressChanged("")

	var err error
	w.introFrame, err = w.getIntroOrOutro(w.Intro)
	if err != nil {
		w.handleError(err)
		return
	}

	w.outroFrame, err = w.getIntroOrOutro(w.Outro)
	if err != nil {
		w.handleError(err)
		return
	}

	defer func() {
		if w.introFrame.VideoFile != "" {
			os.Remove(w.introFrame.VideoFile)
		}
		if w.outroFrame.VideoFile != "" {
			os.Remove(w.outroFrame.VideoFile)
		}
	}()

	for i := 0; i < w.VideoCount; i++ {
		err := w.generateVideo(i)
		if err != nil {
			w.handleError(err)
			return
		}
		if w.isCancelled(true) {
			return
		}
		w.progress++
		w.handleProgressChanged("")
		if w.isCancelled(true) {
			return
		}
	}

	log.Println("video generated")
}

func (w *VideoFFmpegGeneratorWorker) generateVideo(i int) error {
	files, err := utils.GetFilesFromPath(w.Input, []string{"image"})
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("files list is empty")
	}

	frames := []*Frame{}
	usedFiles := map[string]byte{}

	getFrame := func() (*Frame, error) {
		file := ""
		if len(usedFiles) == len(files) {
			return nil, errors.New("files ended")
		}
		for {
			i := utils.RandRange(0, len(files)-1)
			file = files[i]
			if _, ok := usedFiles[file]; !ok {
				usedFiles[file] = 0
				break
			}
		}
		frameDur := utils.RandRange(w.SlideDurationFrom, w.SlideDurationTo) * 1000
		return &Frame{
			File:     file,
			Duration: frameDur,
		}, nil
	}

	if w.Intro != "" {
		frames = append(frames, &w.introFrame)
	}

	audioFile := w.getAudio()

	var duration int

	switch w.VideoDurationType {
	case ByIntervalDuration:
		duration = utils.RandRange(w.VideoDurationFrom, w.VideoDurationTo) * 1000
		for dur := w.introFrame.Duration + w.outroFrame.Duration; dur < duration; {
			frame, err := getFrame()
			if err != nil {
				return err
			}
			frames = append(frames, frame)
			dur += frame.Duration
		}

	case BySlideCountDuration:
		count := utils.RandRange(w.SlideCountFrom, w.SlideCountTo)
		for i := 0; i < count; i++ {
			frame, err := getFrame()
			if err != nil {
				return err
			}
			frames = append(frames, frame)
			duration += frame.Duration
		}

	case ByAudioDuration:
		if audioFile == "" {
			return errors.New("cannot find audio file")
		}
		d, err := utils.GetAudioDuration(audioFile)
		if err != nil {
			return err
		}
		duration = int(d.Milliseconds())
		for dur := w.introFrame.Duration + w.outroFrame.Duration; dur < duration; {
			frame, err := getFrame()
			if err != nil {
				return err
			}
			frames = append(frames, frame)
			dur += frame.Duration
		}

	default:
		return errors.New("unknown duration type")
	}

	if w.Outro != "" {
		frames = append(frames, &w.outroFrame)
	}

	w.wg.Add(len(frames))

	framesChannel := make(chan *Frame, w.Threads)

	for i := 0; i < w.Threads; i++ {
		go w.startCreateFramesWorker(framesChannel)
	}

	w.filesWorker(frames, framesChannel)

	w.wg.Wait()

	if w.isCancelled(true) {
		return nil
	}

	logger.Notice("frames generated")

	resultfile := strconv.Itoa(i)

	if audioFile != "" {
		filename := filepath.Base(audioFile)
		name := strings.TrimSuffix(filename, path.Ext(filename))
		resultfile = name + "_" + strconv.Itoa(i)
	}

	return w.concatFrames(frames, resultfile, audioFile)
}

func (w *VideoFFmpegGeneratorWorker) filesWorker(files []*Frame, framesChannel chan *Frame) {
	for _, file := range files {
		framesChannel <- file
	}
	close(framesChannel)
}

func (w *VideoFFmpegGeneratorWorker) startCreateFramesWorker(framesChannel chan *Frame) {
	var err error
	for frame, opened := <-framesChannel; opened; frame, opened = <-framesChannel {
		if w.isCancelled(true) || frame.VideoFile != "" || err != nil {
			w.wg.Done()
			continue
		}

		filename := filepath.Base(frame.File)
		name := strings.TrimSuffix(filename, path.Ext(filename))
		durationMs := utils.RandRange(w.SlideDurationFrom, w.SlideDurationTo) * 1000
		frame.VideoFile = filepath.Join(paths.Temp, fmt.Sprintf("%s_%d%s", name, time.Now().Nanosecond(), w.Ext))
		cmd := exec.Command(config.GetFFmpegBin(), "-y",
			"-loop", "1",
			"-i", frame.File,
			"-c:v", "libx264",
			"-preset", "ultrafast",
			"-crf", "18",
			"-pix_fmt", "yuv420p",
			"-filter:v", fmt.Sprintf("fps=fps=%d,scale=1920:1080", w.FPS),
			"-t", fmt.Sprintf("%dms", durationMs),
			frame.VideoFile)
		var out []byte
		out, err = cmd.CombinedOutput()
		if err != nil {
			logger.Error(errors.New("error on generate frame, out: " + string(out)))
			w.wg.Done()
			w.handleError(err)
			continue
		}
		w.wg.Done()
	}

}

func (w *VideoFFmpegGeneratorWorker) getAudio() string {
	if w.AudioFile == "" {
		return ""
	}

	isDir, err := utils.IsDir(w.AudioFile)
	if err != nil {
		return ""
	}

	if isDir {
		files, err := utils.GetFilesFromPath(w.AudioFile, []string{"audio"})
		if err != nil {
			return ""
		}
		if len(files) == 0 {
			return ""
		}
		return files[utils.RandRange(0, len(files)-1)]
	}

	return w.AudioFile
}

func (w *VideoFFmpegGeneratorWorker) concatFrames(frames []*Frame, filename string, audioFile string) error {
	defer func() {
		for _, f := range frames {
			if f.File == "" || f.File == w.introFrame.File || f.File == w.outroFrame.File {
				continue
			}
			os.Remove(f.VideoFile)
			if w.DeleteImages && f.File != "" {
				os.Remove(f.File)
			}
		}
	}()

	files := make([]string, len(frames))
	for i, frame := range frames {
		abs, _ := filepath.Abs(frame.VideoFile)
		files[i] = fmt.Sprintf("file '%s'", abs)
	}

	filesPath := filepath.Join(w.ResultPath, "files.txt")

	if err := ioutil.WriteFile(filesPath, []byte(strings.Join(files, "\r\n")), 0666); err != nil {
		return err
	}
	defer os.Remove(filesPath)

	iVideoArg := filesPath

	duration := getFramesDuration(frames)

	audioArg := []string{}

	if audioFile != "" {
		dur, err := utils.GetAudioDuration(audioFile)
		if err != nil {
			return err
		}
		audioDuration := int(dur.Milliseconds())

		aFilename := filepath.Base(audioFile)

		if w.IfAudioLonger == 20 && audioDuration > duration {
			audioArg = []string{"-i", filepath.Join(paths.Temp, aFilename)}
			err = utils.CopyFile(audioFile, audioArg[1])
			if err != nil {
				return err
			}
			err = utils.AdaptAudio(audioArg[1], float64(duration)/1000, w.IfAudioLonger)
			if err != nil {
				return err
			}
		} else if w.IfVideoLonger == 20 && duration > audioDuration {
			audioArg = []string{
				"-filter_complex",
				fmt.Sprintf("amovie='%s':loop=0,asetpts=N/SR/TB", strings.ReplaceAll(strings.ReplaceAll(audioFile, "\\", "/"), ":/", "\\:/")),
			}
		} else {
			audioArg = []string{"-i", audioFile}
		}
	}

	args := append(audioArg, []string{
		"-y",
		"-f", "concat",
		"-safe", "0",
		"-i", iVideoArg,
		"-shortest",
		//"-t", fmt.Sprintf("%dms", duration),
		"-c:v", "copy",
	}...)

	resFile := filepath.Join(w.ResultPath, fmt.Sprintf("%s%s", filename, w.Ext))

	if !w.Overlay.Enabled {
		args = append(args, resFile)
		cmd := exec.Command(config.GetFFmpegBin(), args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(string(out))
			log.Println(err)
		}

		return err
	}

	tempFile := filepath.Join(w.ResultPath, fmt.Sprintf("temp_%s%s", filename, w.Ext))
	args = append(args, tempFile)
	cmd := exec.Command(config.GetFFmpegBin(), args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Println(err)
		return err
	}
	defer os.Remove(tempFile)

	return utils.SetOverlay(tempFile, resFile, w.Overlay)
}

func getFramesDuration(frames []*Frame) (res int) {
	for _, frame := range frames {
		res += frame.Duration
	}
	return
}
