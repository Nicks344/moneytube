package upload

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/tidwall/gjson"
)

func (vu *VideoUploader) watchForUpload(ctx context.Context) *bool {
	var uploaded bool

	go func() {
		reqBodyCh := vu.ytBrowser.ListenForRequestBody(ctx, "upload.youtube.com/?authuser")

		for {
			body := <-reqBodyCh
			if body != "" {
				uploaded = gjson.Get(body, "status").String() == "STATUS_SUCCESS"
				if uploaded {
					return
				}
			}
		}
	}()

	return &uploaded
}

func (vu *VideoUploader) watchForProcessed(ctx context.Context) *bool {
	var processed bool

	processingBody := vu.ytBrowser.ListenForRequestBody(ctx, "get_creator_videos")
	go func() {
		for !processed {
			body := <-processingBody
			if body == "" {
				continue
			}

			status := gjson.Get(body, "videos.0.status").String()

			if status == "VIDEO_STATUS_PROCESSED" {
				processed = true
				break
			}
		}
	}()

	return &processed
}

func (vu *VideoUploader) waitForVideoReady() (params VideoMetadataParams, err error) {
	ctx, cancel := context.WithTimeout(vu.ctx, 5*time.Minute)
	defer cancel()

	bodyCh := vu.ytBrowser.WaitForRequestBody(ctx, "upload/createvideo")
	reqCh := vu.ytBrowser.ListenForRequests(ctx, "upload/createvideo")

	for {
		select {
		case <-ctx.Done():
			err = errors.New("wait for new session timeout")
			return

		case requestInfo, _ := <-reqCh:
			var query string
			query, err = url.QueryUnescape(strings.Split(requestInfo.Request.URL, "?")[1])
			if err != nil {
				return
			}

			params.Query = query

			err = mapstructure.Decode(requestInfo.Request.Headers, &params.Headers)
			if err != nil {
				return
			}

			err = json.Unmarshal([]byte(requestInfo.Request.PostData), &params.Metadata)
			if err != nil {
				return
			}
			params.Metadata.VideoReadMask = GetVideoReadMask()
			reqCh = nil

		case body, _ := <-bodyCh:
			jbody := gjson.Parse(body)
			status := jbody.Get("contents.uploadFeedbackItemRenderer.contents.0.uploadStatus.uploadStatus").String()
			if status == "REJECTED" {
				reason := jbody.Get("contents.uploadFeedbackItemRenderer.contents.0.uploadStatus.uploadStatusReason").String()
				if reason == "UPLOAD_STATUS_REASON_RATE_LIMIT_EXCEEDED" {
					err = limitError
					return
				}
				err = errors.New(jbody.Get("contents.uploadFeedbackItemRenderer.contents.0.uploadStatus.detailedMessage.simpleText").String())
				return
			}
			params.Metadata.EncryptedVideoID = gjson.Get(body, "videoId").String()
			bodyCh = nil
		}

		if reqCh == nil && bodyCh == nil {
			return
		}
	}
}
