package uibridge

type UpdateRPC struct {
	OnUpdating    func(string)
	OnUpdated     func()
	OnUpdateError func(string)
}

func OnceUpdateCancelled(action func()) {
	rpcEvents.add("onUpdateCancelled", func(interface{}) error {
		action()
		return nil
	})
}

func handleUpdateCancelled() {
	if event, ok := rpcEvents.get("onUpdateCancelled"); ok {
		event(nil)
		rpcEvents.remove("onUpdateCancelled")
	}
}
