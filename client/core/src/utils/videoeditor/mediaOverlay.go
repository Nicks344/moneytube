package videoeditor

import (
	"fmt"
	"strings"
)

const (
	OP_Chormakey   = "chromakey"
	OP_Perspective = "perspective"
	OP_Coordinate  = "coordinate"
)

type MediaOverlayOpts struct {
	Blur     BlurOpts     `json:"blur"`
	Color    ColorOpts    `json:"color"`
	File     string       `json:"file"`
	Position PositionOpts `json:"position"`

	Preview string    `json:"preview"`
	Time    TimeRange `json:"time"`
}

func (opts *MediaOverlayOpts) IsBrightnessChanged() bool {
	return opts.Color.Brightness != 100
}

type ColorOpts struct {
	Brightness int    `json:"brightness"`
	Curves     string `json:"curves"`
}

type BlurOpts struct {
	Degree      int  `json:"degree"`
	Transparent int  `json:"transparent"`
	Enabled     bool `json:"enabled"`
	Padding     int  `json:"padding"`
}

type PositionOpts struct {
	Type           string          `json:"type"`
	ChromakeyColor string          `json:"chromakey"`
	Perspective    PerspectiveOpts `json:"perspective"`
	Coordinate     Point           `json:"coordinate"`
}

type PerspectiveOpts struct {
	X0 int `json:"x0"`
	X1 int `json:"x1"`
	X2 int `json:"x2"`
	X3 int `json:"x3"`
	Y0 int `json:"y0"`
	Y1 int `json:"y1"`
	Y2 int `json:"y2"`
	Y3 int `json:"y3"`
}

type MediaOverlayData struct {
	Input  string
	Filter string
}

func getMediaOverlayData(overlayOpts MediaOverlayOpts, mainInput string) (res MediaOverlayData, err error) {
	res.Input = overlayOpts.File

	info, err := GetVideoInfo(mainInput)
	if err != nil {
		return
	}

	enable := getEnableParam(overlayOpts.Time)

	prefixOpts := []string{}

	if overlayOpts.Color.Curves != "" {
		acvFile := overlayOpts.Color.Curves

		acvFile = strings.ReplaceAll(acvFile, "\\", "/")
		acvFile = strings.Replace(acvFile, ":/", "\\\\:/", 1)

		prefixOpts = append(prefixOpts, fmt.Sprintf("curves=psfile=%s", acvFile))
	}

	if overlayOpts.IsBrightnessChanged() {
		prefixOpts = append(prefixOpts, fmt.Sprintf("eq=brightness=%f", (float64(overlayOpts.Color.Brightness)/100)-1))
	}

	if overlayOpts.Position.Type != OP_Coordinate {
		prefixOpts = append(prefixOpts, fmt.Sprintf("scale=%d:%d", info.Width, info.Height))
	}

	if overlayOpts.Blur.Enabled {
		wPad := fmt.Sprintf("(in_w/100)*%d", overlayOpts.Blur.Padding)
		hPad := fmt.Sprintf("(in_h/100)*%d", overlayOpts.Blur.Padding)
		blur := "split=2[eq1][eq2];[eq1]crop=in_w-({wp}*2):in_h-({hp}*2):{wp}:{hp}[crop];[eq2]format=argb,colorchannelmixer=aa=%f,gblur=%d[blur];[blur][crop]overlay=(W/100)*%d:(H/100)*%d"
		blur = strings.ReplaceAll(blur, "{wp}", wPad)
		blur = strings.ReplaceAll(blur, "{hp}", hPad)
		prefixOpts = append(prefixOpts, fmt.Sprintf(blur, float64(overlayOpts.Blur.Transparent-100)/-100, overlayOpts.Blur.Degree, overlayOpts.Blur.Padding, overlayOpts.Blur.Padding))
	}

	stream := "[input_stream]"
	var prefix string
	if len(prefixOpts) > 0 {
		stream = "[prefix]"
		prefix = "[input_stream]" + strings.Join(prefixOpts, ",") + stream + ";"
	}

	switch overlayOpts.Position.Type {
	case "chromakey":
		res.Filter = fmt.Sprintf("%s[last_stream]colorkey=0x%s:0.3:0.1[m];%s[m]overlay%s",
			prefix, strings.TrimPrefix(overlayOpts.Position.ChromakeyColor, "#"), stream, enable,
		)
		return

	case "perspective":
		res.Filter = fmt.Sprintf("%s%spad=iw+4:ih+4:2:2:black@0,perspective=x0=%d:y0=%d:x1=%d:y1=%d:x2=%d:y2=%d:x3=%d:y3=%d:sense=1:eval=0[m];[last_stream][m]overlay=0:0%s",
			prefix, stream,
			overlayOpts.Position.Perspective.X0, overlayOpts.Position.Perspective.Y0, overlayOpts.Position.Perspective.X1, overlayOpts.Position.Perspective.Y1,
			overlayOpts.Position.Perspective.X2, overlayOpts.Position.Perspective.Y2, overlayOpts.Position.Perspective.X3, overlayOpts.Position.Perspective.Y3,
			enable,
		)

	case "coordinate":
		res.Filter = fmt.Sprintf("%s[last_stream]%soverlay=%d:%d%s", prefix, stream, overlayOpts.Position.Coordinate.X, overlayOpts.Position.Coordinate.Y, enable)
	}

	return
}
