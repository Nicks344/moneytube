package moneytubemodel

import (
	"errors"
	"strings"
	"sync"

	"github.com/Nicks344/moneytube/client/core/src/utils"
)

const (
	SSTV_Runs    = 10
	SSTV_Uploads = 20
)

const (
	CM_AllowAll   = 0
	CM_Automated  = 10
	CM_Approved   = 20
	CM_DisableAll = 30
)

type UploadData struct {
	ID         int   `bson:"_id"`
	AccountIDs []int `flag:"account-ids"`

	UploadDataFields `bson:",inline" mapstructure:",squash" flag:"!noprefix"`
}

type UploadDataFields struct {
	WithProcessing               bool
	DisableComments              bool
	AgeRestrictions              bool
	OrderComments                bool
	ShowRating                   bool
	NotifySubscribers            bool
	ShowStatistic                bool
	Category                     string
	Language                     string
	PauseFrom                    int
	PauseTo                      int
	UploadCountFrom              int
	UploadCountTo                int
	IsScheduled                  bool
	ScheduleTime                 string
	ScheduleStep                 int
	SkipErrors                   bool
	WaitVideoInFolder            bool
	ClearFilesAfterSuccessUpload bool
	IsDeferred                   bool
	DeferTime                    string
	DeferStep                    int
	CommentMode                  int
	Videos                       UploadVideoOptions
	Envelopes                    UploadEnvelopesOptions
	Titles                       UploadTitlesOptions
	Tags                         UploadOptions
	Descriptions                 UploadOptions
	Comments                     UploadCommentsOptions
	Hints                        HintsOptions
	Scheduler                    UploadDataScheduler
}

type UploadDataScheduler struct {
	Enabled     bool
	StartTime   string
	Interval    int
	UploadCount int
	StopVariant int
	StopCount   int
}

type UploadVideoOptions struct {
	UploadOptions `mapstructure:",squash" flag:"!noprefix"`
	RenameToTitle bool
}

type UploadTitlesOptions struct {
	UploadOptions `mapstructure:",squash" flag:"!noprefix"`
	IsGetFilename bool
}

type UploadEnvelopesOptions struct {
	UploadOptions       `mapstructure:",squash" flag:"!noprefix"`
	IsRandomFromPropose bool
}

type UploadCommentsOptions struct {
	UploadOptions `mapstructure:",squash" flag:"!noprefix"`
	AddComment    bool
	FixComment    bool
}

type HintsOptions struct {
	AddHints  bool
	HintsList []HintOptions
}

type HintOptions struct {
	Type    int
	Time    string
	Data    UploadOptions
	Message UploadOptions
	Teaser  UploadOptions
}

type UploadOptions struct {
	sync.Mutex `flag:""`

	List     []string `type:"string-slice-urlencoded"`
	Cycle    bool
	IsRandom bool
	Index    int `flag:""`
}

func (uo *UploadOptions) ClearFromEmpty() {
	uo.Lock()
	defer uo.Unlock()

	uo.List = utils.ClearSliceFromEmpty(uo.List)
}

func (uo *UploadOptions) GetOne() (string, error) {
	uo.Lock()
	defer uo.Unlock()

	if len(uo.List) == 0 {
		return "", errors.New("list is empty")
	}

	if uo.IsRandom {
		uo.Index = utils.RandRange(0, len(uo.List)-1)
	}

	defer func() {
		if !uo.Cycle {
			uo.List = utils.RemoveStr(uo.List, uo.Index)
		} else if !uo.IsRandom {
			uo.Index++
			if uo.Index >= len(uo.List) {
				uo.Index = 0
			}
		}
	}()

	return strings.Trim(uo.List[uo.Index], "\r\n "), nil
}

func (uo *UploadOptions) Return(entry string) {
	uo.Lock()
	defer uo.Unlock()

	if uo.Cycle {
		return
	}

	uo.List = append(uo.List, entry)
}
