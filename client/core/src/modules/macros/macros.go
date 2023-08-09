package macros

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/serverAPI"

	"github.com/m1/gospin"
)

type StaticMacroses struct {
	VideoTitle   string
	VideoLink    string
	ChannelTitle string
	ChannelLink  string
	VideoTags    string
}

func Execute(text string, macroses StaticMacroses) string {
	text, _ = gospin.New(nil).Spin(text)
	text = executeStaticMacroses(text, macroses)
	text = executeUserMacroses(text)

	return text
}

func executeStaticMacroses(text string, macroses StaticMacroses) string {
	val := reflect.ValueOf(macroses)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		macros := fmt.Sprintf("[%s]", typeField.Name)
		text = strings.ReplaceAll(text, macros, valueField.Interface().(string))
	}
	return text
}

func executeUserMacroses(text string) string {
	result, err := serverAPI.ExecuteUserMacroses(text)
	if err != nil {
		logger.Error(err)
	}
	return result

}
