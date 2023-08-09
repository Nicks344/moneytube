package uibridge

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Nicks344/moneytube/client/core/src/utils/update"

	rpc "github.com/hprose/hprose-golang/rpc/websocket"
)

type RPCService struct {
	Launch func(string, string)
	Login  func()

	UpdateRPC
	RenderRPC
}

var RPCClient *rpc.WebSocketClient
var Endpoint *RPCService

func Connect(port int) {
	RPCClient = rpc.NewWebSocketClient(fmt.Sprintf("ws://127.0.0.1:%d/", port))
	RPCClient.SetMaxConcurrentRequests(10)
	RPCClient.SetTimeout(time.Hour * 10)
	RPCClient.UseService(&Endpoint)
}

func Serve(port int) {
	service := rpc.NewWebSocketService()

	service.AddFunction("OnRenderProgress", handleRenderProgress)
	service.AddFunction("OnUpdateCancelled", handleUpdateCancelled)
	service.AddFunction("Activate", activate)
	service.AddFunction("GetUpdateProgress", update.GetProgress)

	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), service)
}

func OnActivate(action func(key string) error) {
	rpcEvents.add("onActivate", func(data interface{}) error {
		if k, ok := data.(string); ok {
			action(k)
			return nil
		}
		return errors.New("invalid key")
	})
}

func activate(key string) error {
	if f, ok := rpcEvents.get("onActivate"); ok {
		return f(key)
	}
	return errors.New("unknown error")
}
