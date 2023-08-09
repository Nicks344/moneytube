package videoeditor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nicks344/moneytube/client/core/src/paths"
	"github.com/meandrewdev/transcoder/ffmpeg"
)

func ConcatVideoWithReEncode(files []string, resfile string) error {
	var filter string

	for i := range files {
		filter += fmt.Sprintf("[%d]", i)
	}

	filter += fmt.Sprintf("concat=n=%d:v=1:a=1", len(files))

	opts := ffmpeg.Options{
		FilterComplex: &filter,
		Overwrite:     &overwrite,
		Inputs:        files,
	}

	_, err := getFfmpegCmd().Output(resfile).Start(opts)
	return err
}

func ConcatVideo(files []string, resfile string) error {
	for i := range files {
		files[i] = fmt.Sprintf("file '%s'", files[i])
	}

	tempFile := filepath.Join(paths.Temp, "files.txt")
	defer os.Remove(tempFile)

	err := ioutil.WriteFile(tempFile, []byte(strings.Join(files, "\r\n")), 0666)
	if err != nil {
		return err
	}

	codec := "copy"
	iformat := "concat"
	isafe := "0"
	opts := ffmpeg.Options{
		Overwrite:   &overwrite,
		VideoCodec:  &codec,
		AudioCodec:  &codec,
		InputFormat: &iformat,
		InputSafe:   &isafe,
		Inputs:      []string{tempFile},
	}

	_, err = getFfmpegCmd().Output(resfile).Start(opts)
	return err
}
