package moneytubemodel

const (
	UTSStopped   = 10 //Остановлен
	UTSStopping  = 15 //В процессе остановки
	UTSInProcess = 20 //В процессе
	UTSWaiting   = 25 //Ожидает паузу
	UTSReady     = 30 //Выполнен
	UTSError     = 40 //Ошибка
)

type UploadTask struct {
	ID           int `bson:"_id"`
	AccountID    int
	Account      Account     `bson:"-" json:"-" flag:""`
	DetailsID    int         `flag:""`
	Details      *UploadData `bson:"-" json:"-"`
	Count        int
	Progress     int    `flag:""`
	Status       int    `flag:""`
	ErrorMessage string `flag:""`
	IsScheduled  bool   `flag:""`
	IsFromAPI    bool   `flag:""`
	APIID        int    `flag:""`
	ScheduleTime string `flag:""`

	Scheduler UploadTaskScheduler `flag:""`
}

func (this *UploadTask) IsRunning() bool {
	return this.Status == UTSInProcess || this.Status == UTSWaiting
}

type UploadTaskScheduler struct {
	Enabled         bool
	SecondStartTime string
	Progress        int
}
