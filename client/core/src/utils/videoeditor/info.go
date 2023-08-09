package videoeditor

import (
	"math"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/meandrewdev/transcoder"

	"github.com/spf13/viper"
)

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

type VideoInfo struct {
	Width         int
	Height        int
	FPS           int
	Duration      time.Duration
	VideoTimeBase int
	VideoCodec    string
	AudioTimeBase int
	AudioCodec    string
}

func GetVideoInfo(filename string) (info VideoInfo, err error) {
	meta, err := getFfmpegCmd().Probe(filename)
	if err != nil {
		return
	}

	vstream := getFirstVideoStream(meta)
	if vstream != nil {
		info.VideoCodec = vstream.CodecName
		info.Width = vstream.Width
		info.Height = vstream.Height

		fpsSplit := strings.Split(vstream.RFrameRrate, "/")
		var fpsAvg, fpsDelimiter float64
		fpsAvg, err = strconv.ParseFloat(fpsSplit[0], 64)
		if err != nil {
			return
		}
		fpsDelimiter, err = strconv.ParseFloat(fpsSplit[1], 64)
		if err != nil {
			return
		}
		info.FPS = int(math.Round(fpsAvg / fpsDelimiter))

		tbSplit := strings.Split(vstream.TimeBase, "/")
		var tbAvg, tbDelimiter float64
		tbAvg, err = strconv.ParseFloat(tbSplit[1], 64)
		if err != nil {
			return
		}
		tbDelimiter, err = strconv.ParseFloat(tbSplit[0], 64)
		if err != nil {
			return
		}
		info.VideoTimeBase = int(math.Round(tbAvg / tbDelimiter))

		info.Duration, err = time.ParseDuration(vstream.Duration + "s")
		if err != nil {
			return
		}
	}

	astream := getFirstAudioStream(meta)
	if astream != nil {
		info.AudioCodec = astream.CodecName
		tbSplit := strings.Split(astream.TimeBase, "/")
		var tbAvg, tbDelimiter float64
		tbAvg, err = strconv.ParseFloat(tbSplit[1], 64)
		if err != nil {
			return
		}
		tbDelimiter, err = strconv.ParseFloat(tbSplit[0], 64)
		if err != nil {
			return
		}
		info.AudioTimeBase = int(math.Round(tbAvg / tbDelimiter))
	}

	return
}

func getFirstVideoStream(meta transcoder.Metadata) *transcoder.Streams {
	for _, stream := range meta.Streams {
		if stream.CodecType == "video" {
			return &stream
		}
	}

	return nil
}

func getFirstAudioStream(meta transcoder.Metadata) *transcoder.Streams {
	for _, stream := range meta.Streams {
		if stream.CodecType == "audio" {
			return &stream
		}
	}

	return nil
}
