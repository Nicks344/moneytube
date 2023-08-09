package paths

import (
	"os"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/utils"
)

func init() {
	initPath(Data)
	initPath(Sessions)
	initPath(Templates)
	initPath(AETemplates)
	initPath(UniqueTemplates)
	initPath(Fonts)
	initPath(Temp)
}

func initPath(path string) {
	if ok, err := utils.IsDir(path); !ok || err != nil {
		err := os.Mkdir(path, 0777)
		if err != nil {
			logger.Error(err)
			panic(err)
		}
	}
}
