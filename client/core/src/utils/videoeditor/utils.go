package videoeditor

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/meandrewdev/transcoder/ffmpeg"
)

type TimeRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func getEnableParam(t TimeRange) string {
	if t.From > 0 && t.To < 0 {
		return fmt.Sprintf(":enable='gte(t,%d)'", t.From)
	}
	if t.From >= 0 && t.To > 0 {
		return fmt.Sprintf(":enable='between(t,%d,%d)'", t.From, t.To)
	}

	return ""
}

func randName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return prefix + hex.EncodeToString(randBytes) + suffix
}

func appendVideoFilter(opts ffmpeg.Options, filter string) ffmpeg.Options {
	if opts.VideoFilter == nil {
		opts.VideoFilter = &filter
	} else {
		vf := *opts.VideoFilter
		vf += "," + filter
		opts.VideoFilter = &vf
	}

	return opts
}
