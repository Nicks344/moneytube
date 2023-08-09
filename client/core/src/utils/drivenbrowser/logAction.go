package drivenbrowser

import (
	"context"
	"fmt"

	"github.com/meandrewdev/logger"

	"github.com/chromedp/chromedp"
)

type LogAction struct {
	logMess string
	action  chromedp.Action
}

func (this *LogAction) Do(ctx context.Context) error {
	logger.Notice(fmt.Sprintf("[%s]: %s", ctx.Value("session"), this.logMess))
	return this.action.Do(ctx)
}

func NewLogAction(a chromedp.Action, m string) chromedp.Action {
	return &LogAction{
		logMess: m,
		action:  a,
	}
}
