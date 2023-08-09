package videoeditor

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/paths"
	"github.com/Nicks344/moneytube/client/core/src/utils"

	"github.com/meandrewdev/transcoder/ffmpeg"
)

type Cutmethod int

const (
	CM_DontCutSeconds    Cutmethod = 0
	CM_CutStart          Cutmethod = 10
	CM_CutEnd            Cutmethod = 20
	CM_CutStartAndEnd    Cutmethod = 30
	CM_CutRandStartOrEnd Cutmethod = 40
	CM_CutRandFromMiddle Cutmethod = 50
)

func CutFragments(file string, cutMethod Cutmethod, seconds float64, fragmentCount int, result string) error {
	info, err := GetVideoInfo(file)
	if err != nil {
		return err
	}

	if seconds*float64(fragmentCount) >= info.Duration.Seconds() {
		return errors.New("количество секунд для вырезания больше или равно длительности видео")
	}

	copyTs := true
	opts := getDefaultFfmpegOpts()
	opts.Inputs = []string{file}
	opts.CopyTs = &copyTs
	opts.OutputExtraArgs = map[string]interface{}{
		"-avoid_negative_ts": "1",
	}

	switch cutMethod {
	case CM_DontCutSeconds:
		return errors.New("no cutting")

	case CM_CutStart:
		return cutStart(&opts, info, seconds, result)

	case CM_CutEnd:
		return cutEnd(&opts, info, seconds, result)

	case CM_CutStartAndEnd:
		return cutStartAndEnd(&opts, info, seconds, result)

	case CM_CutRandStartOrEnd:
		mode := utils.RandRange(0, 2)
		switch mode {
		case 0:
			return cutStart(&opts, info, seconds, result)

		case 1:
			return cutEnd(&opts, info, seconds, result)

		case 2:
			return cutStartAndEnd(&opts, info, seconds, result)

		default:
			return fmt.Errorf("unknown cut method: %d, %d", cutMethod, mode)
		}

	case CM_CutRandFromMiddle:
		partDuration := info.Duration.Seconds() / float64(fragmentCount)
		partTimes := make([]float64, fragmentCount)
		for i := range partTimes {
			from := int((partDuration * float64(100*i)) + seconds*100)
			to := int(partDuration * float64(100*(i+1)))
			rnd := utils.RandRange(from, to)
			partTimes[i] = float64(rnd) / 100
		}

		files := make([]string, fragmentCount+1)

		for i := 0; i < len(partTimes)+1; i++ {
			tmp := filepath.Join(paths.Temp, fmt.Sprintf("%d_%s", time.Now().Nanosecond(), path.Ext(result)))

			var from string
			if i == 0 {
				from = "0"
			} else {
				from = fmt.Sprintf("%f", partTimes[i-1])
			}

			var to string
			if i == len(partTimes) {
				to = fmt.Sprintf("%f", info.Duration.Seconds())
			} else {
				to = fmt.Sprintf("%f", partTimes[i]-seconds)
			}

			if err := cut(&opts, from, to, tmp); err != nil {
				return err
			}

			files[i] = tmp
			defer os.Remove(tmp)
		}

		return ConcatVideo(files, result)

	default:
		return fmt.Errorf("unknown cut method: %d", cutMethod)
	}

}

func cutEnd(opts *ffmpeg.Options, info VideoInfo, seconds float64, result string) error {
	return cut(opts, "0", fmt.Sprintf("%f", info.Duration.Seconds()-seconds), result)
}

func cutStart(opts *ffmpeg.Options, info VideoInfo, seconds float64, result string) error {
	return cut(opts, fmt.Sprintf("%f", seconds), fmt.Sprintf("%f", info.Duration.Seconds()), result)
}

func cutStartAndEnd(opts *ffmpeg.Options, info VideoInfo, seconds float64, result string) error {
	return cut(opts, fmt.Sprintf("%f", seconds), fmt.Sprintf("%f", info.Duration.Seconds()-seconds), result)
}

func cut(opts *ffmpeg.Options, seekFrom, seekTo, result string) error {
	opts.SeekTime = &seekFrom
	opts.SeekTimeTo = &seekTo

	_, err := getFfmpegCmd().Output(result).Start(opts)
	return err
}
