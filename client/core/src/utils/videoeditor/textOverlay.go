package videoeditor

import (
	"fmt"
	"strings"
)

type TextOverlayOpts struct {
	Bold            bool      `json:"bold"`
	Color           string    `json:"color"`
	Font            string    `json:"font"`
	Italic          bool      `json:"italic"`
	Position        Point     `json:"position"`
	Preview         string    `json:"preview"`
	Size            int       `json:"size"`
	Text            string    `json:"text"`
	Time            TimeRange `json:"time"`
	Background      bool      `json:"background"`
	BackgroundColor string    `json:"backgroundColor"`
}

func getTextOverlayFilter(overlayOpts TextOverlayOpts) (string, error) {
	fontfileopt := ""
	if overlayOpts.Font != "" {
		fontfile, err := GetGoogleFontPath(overlayOpts.Font, overlayOpts.Bold, overlayOpts.Italic)
		if err != nil {
			return "", err
		}
		fontfile = strings.ReplaceAll(fontfile, "\\", "/")
		fontfile = strings.Replace(fontfile, ":/", "\\\\:/", 1)
		fontfileopt = fmt.Sprintf(`fontfile=%s:`, fontfile)
	}

	enable := getEnableParam(overlayOpts.Time)

	var boxFilter string
	if overlayOpts.Background {
		color := overlayOpts.BackgroundColor
		if len(color) == 3 {
			color = color + color
		}
		boxFilter = fmt.Sprintf(":box=1:boxcolor=0x%s@1.0", color)
	}

	color := overlayOpts.Color
	if len(color) == 3 {
		color = color + color
	}

	vfilter := fmt.Sprintf(`drawtext=%sfontsize=%d:fontcolor=0x%s:x=%d:y=%d:text='%s'%s%s`,
		fontfileopt, overlayOpts.Size, color,
		overlayOpts.Position.X, overlayOpts.Position.Y, overlayOpts.Text, boxFilter, enable,
	)

	return vfilter, nil
}
