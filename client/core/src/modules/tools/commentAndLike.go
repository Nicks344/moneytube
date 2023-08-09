package tools

import (
	"encoding/json"
	"fmt"

	"github.com/meandrewdev/logger"

	"strings"
	"sync"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/macros"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/chromedp/cdproto/cdp"

	"github.com/chromedp/chromedp"
)

type CommentAndLikeWorker struct {
	*ToolService

	AccountIDs  []int
	MaxComments int
	VideoLinks  []string
	Comment     moneytubemodel.UploadOptions
	Like        bool

	sync.Mutex
	accs []moneytubemodel.Account
	wg   sync.WaitGroup
}

func (w *CommentAndLikeWorker) Start() {
	w.wg = sync.WaitGroup{}
	w.maxProgress = len(w.VideoLinks)

	w.handleProgressChanged("")

	w.Comment.ClearFromEmpty()

	var err error
	w.accs, err = model.GetAccountsByIDs(w.AccountIDs)
	if err != nil {
		w.handleError(err)
		return
	}

	w.doWork()

	if w.isCancelled(true) {
		return
	}

	if len(w.VideoLinks) > 0 {
		w.handleError(fmt.Errorf("необработанных ссылок: %d, они оставлены в списке", len(w.VideoLinks)))
		return
	}
}

func (w *CommentAndLikeWorker) onLinkProcessed(link string) {
	w.progress++
	vLinksStr, err := json.Marshal(w.VideoLinks)
	if err != nil {
		w.handleError(err)
		return
	}

	w.handleProgressChanged(string(vLinksStr))
}

func (w *CommentAndLikeWorker) doWork() {
	w.wg.Add(len(w.accs))
	for _, acc := range w.accs {
		go w.startAccount(acc)
	}
	w.wg.Wait()
}

func (w *CommentAndLikeWorker) getLink() string {
	w.Lock()
	defer w.Unlock()

	if w.isCancelled(false) {
		return ""
	}

	if len(w.VideoLinks) == 0 {
		return ""
	}

	link := w.VideoLinks[0]
	w.VideoLinks = utils.RemoveStr(w.VideoLinks, 0)
	return link
}

func (w *CommentAndLikeWorker) returnLink(link string) {
	w.Lock()
	defer w.Unlock()

	w.VideoLinks = append(w.VideoLinks, link)
}

func (w *CommentAndLikeWorker) startAccount(acc moneytubemodel.Account) {
	ybrowser, err := ybrowser.Start(w.ctx, &acc)
	if err != nil {
		return
	}
	defer ybrowser.Stop()
	defer w.wg.Done()

	for i := 0; i < w.MaxComments; i++ {
		macrosData := macros.StaticMacroses{
			ChannelLink:  acc.GetChannelLink(),
			ChannelTitle: acc.ChannelName,
		}
		comment, err := w.Comment.GetOne()
		if err != nil {
			return
		}

		link := w.getLink()
		if link == "" {
			return
		}

		err = ybrowser.GoTo(link)
		if err != nil {
			return
		}

		if len(w.Comment.List) != 0 {
			err = ybrowser.Comment(link, macros.Execute(comment, macrosData))
			if err != nil {
				return
			}
		}

		if w.Like {
			logger.Notice("click to like")
			var likeNodes []*cdp.Node
			err := ybrowser.Browser.Run(chromedp.Nodes("#top-level-buttons ytd-toggle-button-renderer", &likeNodes, chromedp.ByQuery))
			if err != nil {
				return
			}
			if len(likeNodes) == 0 {
				return
			}
			likeNode := likeNodes[0]
			if !strings.Contains(likeNode.AttributeValue("class"), "style-default-active") {
				err = ybrowser.Browser.RunWithTimeout(30*time.Second,
					chromedp.Click("#top-level-buttons ytd-toggle-button-renderer", chromedp.ByQuery),
					chromedp.WaitVisible("paper-toast", chromedp.ByQuery))
				if err != nil {
					return
				}
			}
		}

		w.onLinkProcessed(link)
	}
}
