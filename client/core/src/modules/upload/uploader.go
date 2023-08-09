package upload

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/modules/macros"
	"github.com/Nicks344/moneytube/client/core/src/modules/ybrowser"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/events"
	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/moneytubemodel"

	"github.com/chromedp/cdproto/network"
)

const (
	maxTitleSymbols       = 100
	maxDescriptionSymbols = 5000
	maxTagsSymbols        = 500
)

func NewUploader(task *moneytubemodel.UploadTask) VideoUploader {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	return VideoUploader{
		task:   task,
		ctx:    ctx,
		cancel: cancel,
	}
}

type VideoUploader struct {
	task      *moneytubemodel.UploadTask
	ctx       context.Context
	cancel    context.CancelFunc
	ytBrowser *ybrowser.YoutubeBrowser
	cookies   []*network.Cookie
}

func (vu *VideoUploader) Run(count int) (uploads int, err error) {
	vu.task.Status = moneytubemodel.UTSInProcess
	vu.task.ErrorMessage = ""
	if !vu.task.IsFromAPI {
		vu.saveTask()
	}

	vu.task.Details.Videos.ClearFromEmpty()
	vu.task.Details.Titles.ClearFromEmpty()
	vu.task.Details.Descriptions.ClearFromEmpty()
	vu.task.Details.Tags.ClearFromEmpty()
	vu.task.Details.Envelopes.ClearFromEmpty()
	vu.task.Details.Comments.ClearFromEmpty()
	for i := range vu.task.Details.Hints.HintsList {
		vu.task.Details.Hints.HintsList[i].Data.ClearFromEmpty()
		vu.task.Details.Hints.HintsList[i].Message.ClearFromEmpty()
		vu.task.Details.Hints.HintsList[i].Teaser.ClearFromEmpty()
	}

	for (count == 0 || uploads < count) && vu.task.Progress < vu.task.Count {

		if vu.task.Details.WaitVideoInFolder {
			vu.task.Details.Videos.IsRandom = true
			vu.task.Details.Videos.Cycle = true
		}

		var video string
		video, err = vu.task.Details.Videos.GetOne()
		if err != nil {
			vu.task.ErrorMessage = "Закончились видео"
			err = errors.New(vu.task.ErrorMessage)
			return
		}

		var videosFolder string = video

		for {

			if utils.IsContextCancelled(vu.ctx) {
				vu.task.Status = moneytubemodel.UTSStopped
				if !vu.task.IsFromAPI {
					vu.saveTask()
				}
				return
			}
			if vu.task.Details.WaitVideoInFolder {
				for {
					files, err1 := ioutil.ReadDir(videosFolder)
					if err1 != nil {
						vu.task.ErrorMessage = "Не удалось прочитать файлы в папке"
						err = fmt.Errorf("Не удалось прочитать файлы в папке %s", videosFolder)
						return
					}
					var founded bool = false
					for _, file := range files {
						if filepath.Ext(file.Name()) == ".mp4" {
							video = filepath.Join(videosFolder, file.Name())
							founded = true
							break
						}
					}
					if founded {
						break
					}
					time.Sleep(3 * time.Second)
					if utils.IsContextCancelled(vu.ctx) {
						vu.task.Status = moneytubemodel.UTSStopped
						if !vu.task.IsFromAPI {
							vu.saveTask()
						}
						return
					}
				}
			}

			if _, err = os.Stat(video); os.IsNotExist(err) {
				if vu.task.Details.WaitVideoInFolder {
					time.Sleep(3 * time.Second)
					continue
				} else {
					vu.task.ErrorMessage = "Файл не найден"
					err = fmt.Errorf("Файл %s не найден", video)
					return
				}
			}

			var title string
			if vu.task.Details.Titles.IsGetFilename {
				name := filepath.Base(video)
				title = strings.ReplaceAll(name, filepath.Ext(name), "")
			} else {
				title, err = vu.task.Details.Titles.GetOne()
				if err != nil {
					vu.task.ErrorMessage = "Закончились заголовки"
					err = errors.New(vu.task.ErrorMessage)
					return
				}
			}

			var description string
			description, err = vu.task.Details.Descriptions.GetOne()
			if err != nil {
				vu.task.ErrorMessage = "Закончились описания"
				err = errors.New(vu.task.ErrorMessage)
				return
			}

			var tags string
			tags, err = vu.task.Details.Tags.GetOne()
			if err != nil {
				vu.task.ErrorMessage = "Закончились теги"
				err = errors.New(vu.task.ErrorMessage)
				return
			}

			var envelope string
			if !vu.task.Details.Envelopes.IsRandomFromPropose {
				envelope, err = vu.task.Details.Envelopes.GetOne()
				if err != nil {
					vu.task.ErrorMessage = "Закончились обложки"
					err = errors.New(vu.task.ErrorMessage)
					return
				}
			}

			macrosData := macros.StaticMacroses{
				VideoTitle:   title,
				VideoTags:    tags,
				ChannelTitle: vu.task.Account.ChannelName,
				ChannelLink:  vu.task.Account.GetChannelLink(),
			}

			title = macros.Execute(title, macrosData)
			split := strings.Split(title, " ")
			split[0] = strings.Title(split[0])
			title = strings.Join(split, " ")
			title = cutStringByWords(title, maxTitleSymbols)

			err = vu.uploadVideo(video, title, description, tags, envelope, macrosData)
			if !vu.task.IsFromAPI {
				_, e := model.SaveUploadData(vu.task.Details)
				if e != nil {
					err = e
					return
				}
			}

			if vu.task.Details.ClearFilesAfterSuccessUpload {
				os.Remove(video)
				os.Remove(envelope)
			} else {
				if vu.task.Details.WaitVideoInFolder {
					os.Remove(video)
					os.Remove(envelope)
				}
			}

			if err != nil {
				if vu.task.Details.SkipErrors && err != limitError {
					logger.Error(err)
					logger.Notice("skip error and continue")
				} else {
					vu.task.ErrorMessage = err.Error()
					err = fmt.Errorf("Видео: %s\n\nТекст ошибки: %s", video, err.Error())
					return
				}
			} else {
				uploads++
				vu.task.Progress++
				logger.Notice("video uploaded")
			}

			if utils.IsContextCancelled(vu.ctx) {
				vu.task.Status = moneytubemodel.UTSStopped
				if !vu.task.IsFromAPI {
					vu.saveTask()
				}
				return
			}

			vu.task.Status = moneytubemodel.UTSWaiting
			if !vu.task.IsFromAPI {
				vu.saveTask()
			}

			pause := time.Duration(utils.RandRange(vu.task.Details.PauseFrom, vu.task.Details.PauseTo))
			select {
			case <-vu.ctx.Done():
				vu.task.Status = moneytubemodel.UTSStopped

			case <-time.After(pause * time.Second):
				vu.task.Status = moneytubemodel.UTSInProcess
			}

			if !vu.task.IsFromAPI {
				vu.saveTask()
			}

			if !vu.task.Details.WaitVideoInFolder {
				break
			}
		}
	}

	vu.task.Status = moneytubemodel.UTSStopped
	if !vu.task.IsFromAPI {
		vu.saveTask()
	}

	return
}

func (vu *VideoUploader) Stop() {
	vu.cancel()
	vu.task.Status = moneytubemodel.UTSStopped
	vu.saveTask()
	if vu.ytBrowser != nil {
		vu.ytBrowser.Browser.Stop()
	}
}

func (vu *VideoUploader) saveTask() {
	model.SaveUploadTask(vu.task)
	events.OnUploadTaskUpdated(*vu.task)
}
