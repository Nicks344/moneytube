package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"os"
	"strings"

	"github.com/hprose/hprose-golang/rpc"
	_ "github.com/hprose/hprose-golang/rpc/websocket"
	"github.com/meandrewdev/go-flagsfiller"
	"github.com/Nicks344/moneytube/client/core/src/modules/tools"
	"github.com/Nicks344/moneytube/client/core/src/server/gqlserver/events"
	"github.com/Nicks344/moneytube/moneytubemodel"
)

var commands = map[string]interface{}{
	tools.GenerateVideo:       tools.VideoGeneratorWorker{},
	tools.GenerateImages:      tools.ImagesGeneratorWorker{},
	tools.GenerateAudio:       tools.AudioGeneratorWorker{},
	tools.GenerateCopies:      tools.CopiesGeneratorWorker{},
	tools.CommentAndLike:      tools.CommentAndLikeWorker{},
	tools.ChangeDescription:   tools.ChangeDescriptionWorker{},
	tools.CreatePlaylist:      tools.CreatePlaylistWorker{},
	tools.DeleteVideo:         tools.DeleteVideoWorker{},
	tools.GetLinks:            tools.GetLinksWorker{},
	tools.GenerateVideoFFmpeg: tools.VideoFFmpegGeneratorWorker{},
	"AddUploadTask":           moneytubemodel.UploadData{},
}

func init() {
	flagsfiller.AddConverter("video-generator-data-json", func(s string) (interface{}, error) {
		var data []map[string]tools.LayerData
		if err := json.Unmarshal([]byte(s), &data); err != nil {
			return nil, err
		}
		return data, nil
	})

	flagsfiller.AddConverter("string-slice-urlencoded", func(s string) (interface{}, error) {
		split := strings.Split(s, ",")
		for i := range split {
			var err error
			split[i], err = url.QueryUnescape(split[i])
			if err != nil {
				return nil, err
			}
		}

		return split, nil
	})
}

func main() {
	var temp string
	var help bool
	var printJSON bool

	flag.StringVar(&temp, "cmd", "", "Command [!required]")
	flag.BoolVar(&help, "help", false, "Print help")
	flag.BoolVar(&printJSON, "print-tmpl-json", false, "Print template json for command")

	if len(os.Args) < 2 || !strings.Contains(os.Args[1], "-cmd=") {
		finish("no command specified", "")
	}

	cmd := strings.ReplaceAll(os.Args[1], "-cmd=", "")

	if cmd == "" {
		finish("no command specified", "")
	}

	var data interface{}
	var err error

	switch cmd {
	case tools.GenerateVideo:
		data, err = parse(tools.VideoGeneratorWorker{})

	case tools.GenerateImages:
		data, err = parse(tools.ImagesGeneratorWorker{})

	case tools.GenerateAudio:
		data, err = parse(tools.AudioGeneratorWorker{})

	case tools.GenerateCopies:
		data, err = parse(tools.CopiesGeneratorWorker{})

	case tools.CommentAndLike:
		data, err = parse(tools.CommentAndLikeWorker{})

	case tools.ChangeDescription:
		data, err = parse(tools.ChangeDescriptionWorker{})

	case tools.CreatePlaylist:
		data, err = parse(tools.CreatePlaylistWorker{})

	case tools.DeleteVideo:
		data, err = parse(tools.DeleteVideoWorker{})

	case tools.GetLinks:
		data, err = parse(tools.GetLinksWorker{})

	case tools.GenerateVideoFFmpeg:
		data, err = parse(tools.VideoFFmpegGeneratorWorker{})

	case "AddUploadTask":
		data, err = parseUploadData(help || printJSON)

	case "StartUploadTask":
		var id int
		flag.IntVar(&id, "id", 0, "Task ID [!required]")
		flag.Parse()

		if id == 0 && !(help || printJSON) {
			finish("template id not specified", "")
		}

		data = id

	default:
		finish(`command "`+cmd+`" not found`, "")
	}

	if printJSON {
		data, ok := commands[cmd]
		if !ok {
			finish("no json struct for this command", "")
		}

		str, err := json.Marshal(data)
		if err != nil {
			finish(err.Error(), "")
		}

		fmt.Println(string(str))
		return
	}

	if help {
		flag.PrintDefaults()
		return
	}

	if err != nil {
		finish(err.Error(), "")
	}

	connect()

	ans, err := stub.RunTool(cmd, data)
	if err != nil {
		finish(err.Error(), "")
	}

	time.Sleep(2 * time.Second)

	finish("", ans)
}

type Stub struct {
	RunTool           func(string, interface{}) (string, error)
	GetUploadTemplate func(name string) (tmpl moneytubemodel.UploadDataTemplate, err error)
}

var stub Stub
var connected bool

func connect() {
	if connected {
		return
	}

	client := rpc.NewClient("ws://127.0.0.1:10030/")
	client.SetTimeout(time.Hour * 10)
	client.UseService(&stub)

	client.Subscribe("ToolResult", "", nil, func(result events.ToolsResultInput) {
		strBytes, _ := json.Marshal(result)
		fmt.Println(string(strBytes))
	})

	connected = true
}

func parse[T any](tmp T) (T, error) {
	var tmplJSON, tmplFile string
	flag.StringVar(&tmplJSON, "tmpl-json", "", "JSON data [!required if tmpl-file is empty]")
	flag.StringVar(&tmplFile, "tmpl-file", "", "Template file [!required if tmpl-json is empty]")

	filler := flagsfiller.New()
	filler.Fill(flag.CommandLine, &tmp)
	flag.Parse()

	if err := fillByJSON(tmplJSON, tmplFile, &tmp); err != nil {
		return tmp, err
	}

	flag.Parse()
	return tmp, nil
}

func parseUploadData(help bool) (interface{}, error) {
	var tmplName string
	flag.StringVar(&tmplName, "tmpl-name", "", "Template name [!required]")

	var uplData moneytubemodel.UploadData

	filler := flagsfiller.New()
	filler.Fill(flag.CommandLine, &uplData)
	flag.Parse()

	if tmplName == "" && !help {
		return nil, errors.New("template name not specified")
	}

	if help {
		return nil, nil
	}

	connect()

	tmpl, err := stub.GetUploadTemplate(tmplName)
	if err != nil {
		finish(err.Error(), "")
	}

	uplData.AccountIDs = []int{}
	uplData.UploadDataFields = tmpl.UploadDataFields

	flag.Parse()

	return uplData, nil
}

func fillByJSON(tmplJSON, tmplFile string, result interface{}) error {
	if tmplJSON != "" {
		if err := json.Unmarshal([]byte(tmplJSON), &result); err != nil {
			return errors.New("parsing template json failed: " + err.Error())
		}
	} else if tmplFile != "" {
		jsonBytes, err := ioutil.ReadFile(tmplFile)
		if err != nil {
			return errors.New("reading template file failed: " + err.Error())
		}

		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			return errors.New("parsing template json failed: " + err.Error())
		}
	} else {
		return errors.New("no data specified")
	}

	return nil
}

type output struct {
	Error   string `json:",omitempty"`
	Result  string `json:",omitempty"`
	Success bool
}

func finish(err string, result string) {
	str, _ := json.Marshal(output{
		Error:   err,
		Result:  result,
		Success: err == "",
	})
	fmt.Println(string(str))
	code := 0
	if err != "" {
		code = 1
	}
	os.Exit(code)
}
