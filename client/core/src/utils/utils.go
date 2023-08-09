package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/viper"
	"github.com/tcolgate/mp3"
)

func RandRange(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func IsContextCancelled(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func GetWorkpath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func RemoveStr(list []string, i int) []string {
	return append(list[:i], list[i+1:]...)
}

func ClearSliceFromEmpty(list []string) []string {
	res := []string{}
	for _, el := range list {
		el = strings.Trim(el, "\r\n ")
		if el != "" {
			res = append(res, el)
		}
	}

	return res
}

func CopyFile(sourcePath, destPath string) error {
	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

func GetFilesFromPath(path string, assepts []string) (files []string, err error) {
	err = filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if p == path {
				return nil
			}
			return filepath.SkipDir
		}

		mime, err := mimetype.DetectFile(p)
		if err != nil {
			return err
		}

		for _, assept := range assepts {
			if strings.Contains(mime.String(), assept) {
				files = append(files, p)
				break
			}
		}

		return nil
	})
	return
}

func GetAudioDuration(filename string) (time.Duration, error) {
	result := time.Duration(0)
	skipped := 0
	file, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer file.Close()

	decoder := mp3.NewDecoder(file)
	var frame mp3.Frame
	for {
		if err := decoder.Decode(&frame, &skipped); err != nil {
			return result, nil
		}
		result += frame.Duration()
	}
}

func GetVideoDuration(filename string) (time.Duration, error) {
	cmd := exec.Command(viper.GetString("ffprobe_bin"),
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filename)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}
	t := strings.Trim(string(out), " \n\r")
	return time.ParseDuration(t + "s")
}

const (
	CropAudio       = 10
	AccelerateAudio = 20
)

func AdaptAudio(filename string, toSec float64, variant int) error {
	dur, err := GetAudioDuration(filename)
	if err != nil {
		return err
	}
	if dur.Seconds() <= toSec {
		return nil
	}
	switch variant {
	case CropAudio:
		cmd := exec.Command(viper.GetString("ffmpeg_bin"), "-y", "-i", filename, "-ss", "00", "-to", fmt.Sprintf("%.2f", toSec), "-c", "copy", filename+".cropped.mp3")
		err = cmd.Run()

	case AccelerateAudio:
		tempo := dur.Seconds() / toSec
		cmd := exec.Command(viper.GetString("ffmpeg_bin"), "-y", "-i", filename, "-filter:a", fmt.Sprintf("atempo=%.2f", tempo), "-c:a", "libmp3lame", "-q:a", "0", filename+".cropped.mp3")
		err = cmd.Run()

	default:
		err = errors.New("unknown crop variant")
	}

	if err != nil {
		return err
	}
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	err = os.Rename(filename+".cropped.mp3", filename)
	if err != nil {
		return err
	}
	return nil
}

func IsDir(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

type ImageOverlayData struct {
	OverlaySrc string
	Enabled    bool
	X          int
	Y          int
	From       int
	To         int
}

type TextOverlayData struct {
	Text            string
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

func SetOverlay(srcfile string, dstfile string, overlayData ImageOverlayData) error {
	enable := ""
	if overlayData.From > 0 || overlayData.To > 0 {
		enable = fmt.Sprintf(":enable='between(t,%d,%d)'", overlayData.From, overlayData.To)
	}

	cmd := exec.Command(
		viper.GetString("ffmpeg_bin"), "-y",
		"-i", srcfile,
		"-i", overlayData.OverlaySrc,
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-crf", "18",
		"-filter_complex", fmt.Sprintf("[0:v][1:v] overlay=%d:%d%s", overlayData.X, overlayData.Y, enable),
		"-pix_fmt", "yuv420p",
		dstfile,
	)
	return cmd.Run()
}

func IsFileExists(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
