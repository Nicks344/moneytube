package upload

import (
	"errors"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/macros"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils/drivenbrowser"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/chromedp/chromedp"
)

func (vu *VideoUploader) comment(macrosData macros.StaticMacroses) error {
	var err error
	vu.ytBrowser.Stop()
	vu.ytBrowser, err = ybrowser.Start(vu.ctx, &vu.task.Account)
	if err != nil {
		return err
	}
	defer vu.ytBrowser.Stop()

	comment, err := vu.task.Details.Comments.GetOne()
	if err != nil {
		return errors.New("Закончились комментарии")
	}
	_, err = model.SaveUploadData(vu.task.Details)
	if err != nil {
		return err
	}
	err = vu.ytBrowser.Comment(macrosData.VideoLink, macros.Execute(comment, macrosData))
	if err != nil {
		return err
	}
	if vu.task.Details.Comments.FixComment {
		maxTries := 2
		for i := 0; i < maxTries; i++ {
			err = vu.ytBrowser.Browser.RunWithTimeout(time.Second*5,
				drivenbrowser.NewLogAction(chromedp.Click("#action-menu button", chromedp.ByQuery), "click to actions"),
				drivenbrowser.NewLogAction(chromedp.Click("#items ytd-menu-navigation-item-renderer a", chromedp.ByQuery), "click to fix comment"),
				drivenbrowser.NewLogAction(chromedp.WaitVisible(".yt-confirm-dialog-renderer #confirm-button a", chromedp.ByQuery), "wait for confirm"),
				drivenbrowser.NewLogAction(chromedp.Click(".yt-confirm-dialog-renderer #confirm-button a", chromedp.ByQuery), "click to confirm"),
				drivenbrowser.NewLogAction(chromedp.WaitVisible("ytd-pinned-comment-badge-renderer", chromedp.ByQuery), "wait for comment fixed"),
			)
			if err != nil {
				vu.ytBrowser.Browser.Run(chromedp.EvaluateAsDevTools("location.reload()", &[]byte{}))
				continue
			} else {
				break
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (vu *VideoUploader) getCommentOptions() *CommentOptions {
	allowComments := vu.task.Details.CommentMode != moneytubemodel.CM_DisableAll
	commentsMode := ""
	switch vu.task.Details.CommentMode {
	case moneytubemodel.CM_AllowAll:
		commentsMode = "ALL_COMMENTS"

	case moneytubemodel.CM_Automated:
		commentsMode = "AUTOMATED_COMMENTS"

	case moneytubemodel.CM_Approved:
		commentsMode = "APPROVED_COMMENTS"

	case moneytubemodel.CM_DisableAll:
		commentsMode = "UNKNOWN_COMMENT_ALLOWED_MODE"
	}

	return &CommentOptions{
		NewAllowComments:     allowComments,
		NewAllowCommentsMode: commentsMode,
		NewCanViewRatings:    vu.task.Details.ShowStatistic,
		NewDefaultSortOrder:  "MDE_COMMENT_SORT_ORDER_TOP",
	}
}

func (vu *VideoUploader) setCommentsOptions(metadataParams *VideoMetadataParams) error {
	metadataParams.Metadata = VideoMetadata{
		CommentOptions:    vu.getCommentOptions(),
		Context:           metadataParams.Metadata.Context,
		DelegationContext: metadataParams.Metadata.DelegationContext,
		VideoReadMask:     metadataParams.Metadata.VideoReadMask,
		EncryptedVideoID:  metadataParams.Metadata.EncryptedVideoID,
	}

	resp, err := vu.updateMetadata(*metadataParams)
	if err != nil {
		return err
	}

	setCookies := resp.Response().Cookies()
	metadataParams.Headers["Cookie"] = getActualCookiesStr(vu.cookies, setCookies)

	return nil
}
