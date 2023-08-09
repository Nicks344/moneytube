package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/server/serverutils"
)

func Start(port int) {
	mainServer := http.NewServeMux()
	mainServer.Handle("/ws/gql", serverutils.GetGQLWsHandler(config.GetApiKey()))
	mainServer.Handle("/gql", serverutils.GetGQLHTTPHandler(config.GetApiKey()))
	mainServer.HandleFunc("/file", file)
	mainServer.HandleFunc("/network-errors", netErrors)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), mainServer)
}

func netErrors(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error(err)
	}
	logger.Notice(string(body))
}

func file(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	isFrame := r.URL.Query().Get("frame") == "true"

	if isFrame {
		cmd := exec.Command(config.GetFFmpegBin(), "-i", path, "-vframes", "1", "-f", "singlejpeg", "-")
		var buffer bytes.Buffer
		cmd.Stdout = &buffer
		if cmd.Run() != nil {
			panic("could not generate frame")
		}
		w.Write(buffer.Bytes())
		return
	}
	http.ServeFile(w, r, path)
}
