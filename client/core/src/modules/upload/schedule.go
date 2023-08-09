package upload

import (
	"fmt"
	"sync"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/model"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/events"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var schedulerLock sync.Mutex

func init() {
	go func() {
		for {
			time.Sleep(5 * time.Second)

			tasks, _ := model.GetUploadTasks()
			for _, task := range tasks {
				scheduleTick(task)
			}
		}
	}()
}

func scheduleTick(task moneytubemodel.UploadTask) {
	schedulerLock.Lock()
	defer schedulerLock.Unlock()

	if task.Scheduler.Enabled && task.Scheduler.SecondStartTime != "" && !task.IsRunning() {
		now := time.Now()
		t, err := time.Parse("2006-01-02 15:04-0700", task.Scheduler.SecondStartTime+now.Format("-0700"))
		if err != nil {
			return
		}

		if now.After(t) {
			fmt.Printf("%s: Start %d, scheduled at %s\r\n", time.Now().Format("2006-01-02 15:04:03"), task.ID, task.Scheduler.SecondStartTime)
			done := StartUploadTask(&task, task.Details.Scheduler.UploadCount)
			go func() {
				result := <-done

				schedulerLock.Lock()
				defer schedulerLock.Unlock()

				switch task.Details.Scheduler.StopVariant {
				case moneytubemodel.SSTV_Runs:
					task.Scheduler.Progress++

				case moneytubemodel.SSTV_Uploads:
					task.Scheduler.Progress += result.Uploads
				}

				if task.Scheduler.Progress >= task.Details.Scheduler.StopCount || result.Error != nil || task.Details.Scheduler.Interval == 0 {
					task.Scheduler.SecondStartTime = ""
				} else {
					task.Scheduler.SecondStartTime = t.Add(time.Duration(task.Details.Scheduler.Interval) * time.Hour).Format("2006-01-02 15:04")
				}

				events.OnUploadTaskUpdated(task)
				model.SaveUploadTask(&task)
				fmt.Printf("%s: Done %d, scheduled at %s\r\n", time.Now().Format("2006-01-02 15:04:03"), task.ID, task.Scheduler.SecondStartTime)
			}()
		}
	}
}
