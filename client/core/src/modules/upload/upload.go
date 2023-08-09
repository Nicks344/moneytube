package upload

import (
	"context"
	"errors"
	"fmt"

	"github.com/meandrewdev/logger"

	"strings"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/modules/macros"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils/parseLinks"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/chromedp/chromedp"
)

func (vu *VideoUploader) uploadVideo(video, title, description, tags, envelope string, macrosData macros.StaticMacroses) (err error) {
	title = macros.Execute(title, macrosData)
	title = removeInvalidChars(title)
	title = cutStringByWords(title, maxTitleSymbols)

	video = strings.ReplaceAll(video, "\\", "/")
	if vu.task.Details.Videos.RenameToTitle {
		video, err = renameFile(video, title)
		if err != nil {
			return
		}
	}

	vu.ytBrowser, err = ybrowser.Start(vu.ctx, &vu.task.Account)
	if err != nil {
		return
	}
	defer vu.ytBrowser.Stop()

	vu.ytBrowser.Browser.BlockRequest(vu.ctx, "metadata_update")

	err = vu.ytBrowser.GoTo(fmt.Sprintf("https://studio.youtube.com/channel/%s/videos/upload?d=ud", vu.ytBrowser.Account.ChannelID))
	if err != nil {
		return
	}

	uploadCtx, cancelListenUpload := context.WithCancel(vu.ctx)
	defer cancelListenUpload()
	videoUploaded := vu.watchForUpload(uploadCtx)

	processingCtx, cancelListenProcessing := context.WithCancel(vu.ctx)
	defer cancelListenProcessing()
	videoProcessed := vu.watchForProcessed(processingCtx)

	err = vu.ytBrowser.Browser.Run(
		// Upload video file and waiting until the link will be generated
		chromedp.WaitVisible("select-files-button", chromedp.ByID),
		drivenbrowser.NewLogAction(chromedp.SetUploadFiles("input[type=file]", []string{video}, chromedp.ByQuery), "set video file"))
	if err != nil {
		return
	}

	logger.Notice("wait for video ready")
	metadataParams, err := vu.waitForVideoReady()
	if err != nil {
		return err
	}

	macrosData.VideoLink = "https://youtu.be/" + metadataParams.Metadata.EncryptedVideoID

	description = macros.Execute(description, macrosData)
	description = removeInvalidChars(description)
	description = cutStringByWords(description, maxDescriptionSymbols)

	tags = macros.Execute(tags, macrosData)
	tags = removeInvalidChars(tags)
	tags = strings.ReplaceAll(tags, ", ", ",")
	tags = strings.ReplaceAll(tags, ",", ", ")
	tags = cutTags(tags, maxTagsSymbols)

	metadataParams.Metadata.Title = &Title{
		NewTitle:      title,
		ShouldSegment: true,
	}

	metadataParams.Metadata.Description = &Description{
		NewDescription: description,
		ShouldSegment:  true,
	}

	metadataParams.Metadata.TargetedAudience = &TargetedAudience{
		Operation:           "MDE_TARGETED_AUDIENCE_UPDATE_OPERATION_SET",
		NewTargetedAudience: "MDE_TARGETED_AUDIENCE_TYPE_ALL",
	}

	metadataParams.Metadata.Tags = &Tags{
		NewTags: strings.Split(tags, ", "),
	}

	langCode, ok := LangCodes[vu.task.Details.Language]
	if !ok {
		return errors.New("can't find code for language " + vu.task.Details.Language)
	}
	metadataParams.Metadata.AudioLanguage = &AudioLanguage{
		NewAudioLanguage: langCode,
	}

	categoryCode, ok := CategoriesCodes[vu.task.Details.Category]
	if !ok {
		return errors.New("can't find code for category " + vu.task.Details.Category)
	}
	metadataParams.Metadata.Category = &Category{
		NewCategoryID: categoryCode,
	}
	//logger.Notice("4")

	if !vu.task.Details.Envelopes.IsRandomFromPropose {
		var uploadEnvelopeDisabled bool
		var temp string
		err = vu.ytBrowser.Browser.Run(chromedp.AttributeValue("ytcp-thumbnails-compact-editor-uploader-old", "feature-disabled", &temp, &uploadEnvelopeDisabled, chromedp.ByQuery))
		if err != nil {
			return
		}
		if uploadEnvelopeDisabled {
			err = errors.New("для загрузки обложки вам нужно подтвердить аккаунт")
			return
		}

		var dataURI string
		dataURI, err = getFileDataURI(envelope)
		if err != nil {
			return
		}

		metadataParams.Metadata.VideoStill = &VideoStill{
			Operation: "UPLOAD_CUSTOM_THUMBNAIL",
			Image: &Image{
				DataURI: dataURI,
			},
		}
	} else if vu.task.Details.WithProcessing {
		/*
			metadataParams.Metadata.VideoStill = &VideoStill{
				Operation:  "SET_AUTOGEN_STILL",
				NewStillID: utils.RandRange(1, 3),
			}
		*/
	}

	if vu.task.Details.AgeRestrictions {
		metadataParams.Metadata.TargetedAudience = &TargetedAudience{
			NewTargetedAudience: "MDE_TARGETED_AUDIENCE_TYPE_AGE_RESTRICTED",
			Operation:           "MDE_TARGETED_AUDIENCE_UPDATE_OPERATION_SET",
		}
	}

	if vu.task.Details.Comments.AddComment && vu.task.Details.CommentMode == moneytubemodel.CM_DisableAll {
		metadataParams.Metadata.CommentOptions = &CommentOptions{
			NewAllowComments:     true,
			NewAllowCommentsMode: "ALL_COMMENTS",
			NewCanViewRatings:    vu.task.Details.ShowStatistic,
			NewDefaultSortOrder:  "MDE_COMMENT_SORT_ORDER_TOP",
		}
	} else {
		metadataParams.Metadata.CommentOptions = vu.getCommentOptions()
	}

	metadataParams.Metadata.PublishingOptions = &PublishingOptions{
		NewPostToFeed: vu.task.Details.NotifySubscribers,
	}

	metadataParams.Metadata.DraftState = &DraftState{
		Operation: "MDE_DRAFT_STATE_UPDATE_OPERATION_REMOVE_DRAFT_STATE",
	}
	//logger.Notice("6")

	if vu.task.Details.IsDeferred {
		var t time.Time
		t, err = time.Parse("2006-01-02 15:04-0700", vu.task.Details.DeferTime+time.Now().Format("-0700"))
		if err != nil {
			return
		}
		t = t.Add(time.Duration(vu.task.Details.DeferStep*vu.task.Progress) * time.Minute)
		metadataParams.Metadata.PrivacyState = &PrivacyState{
			NewPrivacy: "PRIVATE",
		}
		metadataParams.Metadata.ScheduledPublishing = &ScheduledPublishing{
			Set: Set{
				Privacy: "PUBLIC",
				TimeSec: fmt.Sprintf("%d", t.Unix()),
			},
		}
	} else {
		metadataParams.Metadata.PrivacyState = &PrivacyState{
			NewPrivacy: "PUBLIC",
		}
	}

	logger.Notice("wait for loading")
	//time.Sleep(65 * time.Second)
	if err := waitForTrue(uploadCtx, videoUploaded); err != nil {
		return errors.New("upload timeout")
	}
	//logger.Notice("7")

	vu.cookies, err = vu.ytBrowser.Browser.GetCookies()
	if err != nil {
		return err
	}

	metadataParams.Headers["Cookie"] = getActualCookiesStr(vu.cookies, nil)

	logger.Notice("send metadata")
	if err := vu.submit(&metadataParams); err != nil {
		return err
	}

	if vu.task.Details.WithProcessing {
		logger.Notice("wait for processing")

		if err := waitForTrue(vu.ctx, videoProcessed); err != nil {
			return errors.New("processing timeout")
		}

		// Add hints
		if vu.task.Details.Hints.AddHints {
			hinstData := make([]HintInfo, len(vu.task.Details.Hints.HintsList))

			for i, hInfo := range vu.task.Details.Hints.HintsList {
				tStr := strings.Replace(hInfo.Time, ":", "h", 1)
				tStr = strings.Replace(tStr, ":", "m", 1)
				tStr += "s"
				t, err := time.ParseDuration(tStr)
				if err != nil {
					return err
				}

				data, err := hInfo.Data.GetOne()
				if err != nil {
					return errors.New("закончились ссылки для подсказок")
				}

				hType := HintType(hInfo.Type)

				switch hType {
				case HT_Video:
					data = parseLinks.GetVideoIDsByLinks([]string{data})[0]
					if err != nil {
						return err
					}

				case HT_Playlist:
					data, err = parseLinks.GetPlaylistIDByLink(data)
					if err != nil {
						return err
					}

				case HT_Channel:
					data, err = parseLinks.GetChannelIDByLink(data)
					if err != nil {
						return err
					}

				default:
					return fmt.Errorf("неизвестный тип подсказки: %d", hInfo.Type)
				}

				message, err := hInfo.Message.GetOne()
				if hType == HT_Channel && err != nil {
					return errors.New("закончились сообщения для подсказок")
				}

				teaser, err := hInfo.Teaser.GetOne()
				if hType == HT_Channel && err != nil {
					return errors.New("закончились тизеры для подсказок")
				}

				hinstData[i] = HintInfo{
					Type:    hType,
					Time:    t,
					Data:    data,
					Message: message,
					Teaser:  teaser,
				}
			}

			hintsParams := VideoHintsParams{
				Context: metadataParams.Metadata.Context,
				Query:   metadataParams.Query,
				Headers: metadataParams.Headers,
				VideoID: metadataParams.Metadata.EncryptedVideoID,
				Hints:   hinstData,
			}

			if resp, err := editVideoHints(hintsParams, vu.ytBrowser.Account.Proxy); err != nil {
				return err
			} else {
				setCookies := resp.Response().Cookies()
				metadataParams.Headers["Cookie"] = getActualCookiesStr(vu.cookies, setCookies)
			}
		}

		// Add comment
		if vu.task.Details.Comments.AddComment {
			err = vu.comment(macrosData)
			if err != nil {
				return
			}

			if vu.task.Details.CommentMode == moneytubemodel.CM_DisableAll {
				err = vu.setCommentsOptions(&metadataParams)
				if err != nil {
					return
				}
			}
		}
	}

	return
}

func (vu *VideoUploader) submit(metadataParams *VideoMetadataParams) error {
	resp, err := vu.updateMetadata(*metadataParams)
	if err != nil {
		return err
	}

	setCookies := resp.Response().Cookies()
	metadataParams.Headers["Cookie"] = getActualCookiesStr(vu.cookies, setCookies)

	return nil
}
