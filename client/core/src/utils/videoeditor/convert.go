package videoeditor

import (
	"fmt"
	"os/exec"

	"github.com/spf13/viper"
)

func ConvertVideo(file string, resultFile string) error {
	opts := getDefaultFfmpegOpts()
	opts.Inputs = []string{file}

	out, err := getFfmpegCmd().Output(resultFile).Start(opts)
	if err != nil {
		err = NewGenerateError("converting", err, string(out))
		return err
	}

	return nil
}

func ScaleVideo(file string, info VideoInfo, resultFile string) error {
	inputInfo, err := GetVideoInfo(file)
	if err != nil {
		err = NewGenerateError("get video info", err, "")
		return err
	}

	opts := getDefaultFfmpegOpts()
	opts.Inputs = []string{file}

	if inputInfo.AudioCodec == "" {
		iformat := "lavfi"
		shortest := true
		opts.InputFormat = &iformat
		opts.Inputs = []string{fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=%d", info.AudioTimeBase), file}
		opts.Shortest = &shortest
	}

	if inputInfo.AudioCodec != info.AudioCodec {
		opts.AudioCodec = &info.AudioCodec
	}

	vfilter := fmt.Sprintf("fps=fps=%d,scale=%d:%d", info.FPS, info.Width, info.Height)
	opts.VideoFilter = &vfilter
	opts.OutputExtraArgs = map[string]interface{}{
		"-video_track_timescale": info.VideoTimeBase,
	}

	out, err := getFfmpegCmd().Output(resultFile).Start(opts)
	if err != nil {
		err = NewGenerateError("scaling", err, string(out))
		return err
	}

	return nil
}

func ConvertImageToVideo(file string, info VideoInfo, durationMs int, resultFile string) error {
	cmd := exec.Command(
		viper.GetString("ffmpeg_bin"), "-y",
		"-loop", "1",
		"-i", file,
		"-f", "lavfi",
		"-i", fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=%d", info.AudioTimeBase),
		"-t", fmt.Sprintf("%dms", durationMs),
		"-c:v", "libx264",
		"-t", fmt.Sprintf("%dms", durationMs),
		"-preset", "ultrafast",
		"-crf", "18",
		"-c:a", info.AudioCodec,
		"-video_track_timescale", fmt.Sprintf("%d", info.VideoTimeBase),
		"-vf", fmt.Sprintf("fps=fps=%d,scale=%d:%d", info.FPS, info.Width, info.Height),
		"-pix_fmt", "yuv420p",
		resultFile,
	)
	return cmd.Run()
}
