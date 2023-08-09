package uibridge

type RenderRPC struct {
	Render func(config VideoRenderConfig) (RenderOutput, error)
}

type RenderOutput struct {
	Uid      string
	Output   string
	Workpath string
}

type VideoRenderConfig struct {
	Workpath        string
	AepFile         string
	CompositionName string
	FromRenderQueue bool
	Resolution      string
	ScriptPath      string
	ResultPath      string
	OutputExt       string
	AerenderPath    string
	Assets          []VideoRenderAsset
	Memory          int
}

type VideoRenderAsset struct {
	Type      string
	LayerName string
	Src       string
	Value     interface{}
	Property  string
}

func OnRenderProgress(action func(progress int)) (unsubscribe func()) {
	unsubscribe = func() {
		rpcEvents.remove("onRenderProgress")
	}

	rpcEvents.add("onRenderProgress", func(data interface{}) error {
		if p, ok := data.(int); ok {
			action(p)
		}
		return nil
	})

	return
}

func handleRenderProgress(progress int) {
	if event, ok := rpcEvents.get("onRenderProgress"); ok {
		event(progress)
	}
}
