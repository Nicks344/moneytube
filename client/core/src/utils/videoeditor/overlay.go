package videoeditor

import (
	"fmt"
	"strings"
)

type Overlay interface {
	GetType() string
}

type OverlayConfig struct {
	Name      string        `json:"name"`
	VideoFile string        `json:"video_file"`
	Overlays  []interface{} `json:"overlays"`
}

func DoOverlay(conf OverlayConfig, resultFile string) error {
	videoFile := conf.VideoFile

	opts := getDefaultFfmpegOpts()
	opts.Inputs = []string{videoFile}

	cmd := getFfmpegCmd()

	textFilters := []string{}
	mediaData := []MediaOverlayData{}

	for _, overlay := range conf.Overlays {
		switch overlay.(type) {
		case TextOverlayOpts:
			filter, err := getTextOverlayFilter(overlay.(TextOverlayOpts))
			if err != nil {
				err = NewGenerateError("get text overlay filter", err, fmt.Sprintf("Overlay data: %s", overlay))
				return err
			}
			textFilters = append(textFilters, filter)

		case MediaOverlayOpts:
			overlayData, err := getMediaOverlayData(overlay.(MediaOverlayOpts), conf.VideoFile)
			if err != nil {
				err = NewGenerateError("get media overlay filter", err, fmt.Sprintf("Overlay data: %s", overlay))
				return err
			}
			mediaData = append(mediaData, overlayData)
		}
	}

	lastVideoStream := "[0:v]"
	filtersComplex := []string{}

	for i, filterData := range mediaData {
		stream := fmt.Sprintf("[i_%d]", i)

		opts.Inputs = append(opts.Inputs, filterData.Input)
		istream := fmt.Sprintf("[%d:v]", len(opts.Inputs)-1)

		filterData.Filter = strings.ReplaceAll(filterData.Filter, "[input_stream]", istream)
		filterData.Filter = strings.ReplaceAll(filterData.Filter, "[last_stream]", lastVideoStream)
		filtersComplex = append(filtersComplex, fmt.Sprintf("%s%s", filterData.Filter, stream))
		lastVideoStream = stream
	}

	for i, filter := range textFilters {
		stream := fmt.Sprintf("[t_%d]", i)
		filtersComplex = append(filtersComplex, lastVideoStream+filter+stream)
		lastVideoStream = stream
	}

	if lastVideoStream != "[0:v]" {
		opts.Map = append(opts.Map, lastVideoStream)
	}

	opts.Map = append(opts.Map, "0:a?")

	filterComplex := strings.Join(filtersComplex, ";")
	opts.FilterComplex = &filterComplex

	out, err := cmd.Output(resultFile).Start(opts)
	if err != nil {
		err = NewGenerateError("generating overlay", err, string(out))
		return err
	}

	return nil
}
