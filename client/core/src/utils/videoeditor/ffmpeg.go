package videoeditor

import (
	"github.com/meandrewdev/transcoder"
	"github.com/meandrewdev/transcoder/ffmpeg"
)

var overwrite = true
var vcodec = "libx264"
var acodec = "libmp3lame"
var pixlFormat = "yuv420p"
var crf uint32 = 18
var abitrate = "320k"
var tb = "60"

var ffmpegBinPath, ffprobeBinPath string

func SetPaths(ffmpegPath, ffprobePath string) {
	ffmpegBinPath = ffmpegPath
	ffprobeBinPath = ffprobePath
}

func getDefaultFfmpegOpts() ffmpeg.Options {
	return ffmpeg.Options{
		Overwrite:    &overwrite,
		VideoCodec:   &vcodec,
		AudioCodec:   &acodec,
		PixFmt:       &pixlFormat,
		Crf:          &crf,
		AudioBitrate: &abitrate,
	}
}

func getFfmpegCmd() transcoder.Transcoder {
	ffmpegConf := &ffmpeg.Config{
		FfmpegBinPath:  ffmpegBinPath,
		FfprobeBinPath: ffprobeBinPath,
		//Verbose:        true,
		ProgressEnabled: true,
	}

	return ffmpeg.New(ffmpegConf)
}
