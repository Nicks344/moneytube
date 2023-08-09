package tools

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meandrewdev/logger"
	"github.com/Nicks344/moneytube/client/core/src/config"
	"github.com/Nicks344/moneytube/client/core/src/paths"
	"github.com/Nicks344/moneytube/client/core/src/uibridge"
	"github.com/Nicks344/moneytube/client/core/src/utils"

	"gopkg.in/go-playground/colors.v1"
)

type VideoGeneratorData struct {
	AepFile         string
	CompositionName string
	FromRenderQueue bool
	Resolution      string
	ResultPath      string
	Ext             string
	OutputExt       string
	Quality         int
	IsCompress      bool
	Data            []map[string]LayerData `type:"video-generator-data-json"`
}

type VideoGeneratorWorker struct {
	*ToolService
	VideoGeneratorData `mapstructure:",squash" flag:"!noprefix"`
}

type LayerData struct {
	Data string
	Type string
}

func (w *VideoGeneratorWorker) Start() {
	defer w.cancel()

	go func() {
		<-w.ctx.Done()
		stopAEProcess()
	}()

	w.maxProgress = len(w.Data) * 100
	w.handleProgressChanged("")

	outputScriptPath, err := filepath.Abs(filepath.Join(paths.Temp, "outputmodule.jsx"))
	if err != nil {
		w.handleError(err)
		return
	}

	err = ioutil.WriteFile(outputScriptPath, []byte(aeScript), 0666)
	if err != nil {
		w.handleError(err)
		return
	}
	defer os.Remove(outputScriptPath)

	for i, data := range w.Data {
		if w.isCancelled(true) {
			return
		}

		w.renderVideo(i, data, outputScriptPath)

		if w.isCancelled(true) {
			return
		}
	}
}

func (w *VideoGeneratorWorker) renderVideo(multiplier int, data map[string]LayerData, outputScriptPath string) {
	logger.Notice("start render video " + strconv.Itoa(multiplier))
	unsubscribe := uibridge.OnRenderProgress(func(progress int) {
		if progress < 2 {
			return
		}
		w.progress = (progress - 1) + (100 * multiplier)
		w.handleProgressChanged("")
	})
	defer unsubscribe()

	script := `
			var selectLayer = function(layer) {
				if(layer.source && layer.source instanceof CompItem) {
					if(layer.source.layers.length > 1) {
						throw new Error('Error: nexrender: Layer "' + layer.name + '" are composition and has more than one child')
					} else if(layer.source.layers.length > 0) {
                        layer = layer.source.layers[1]
                    }
				}
			
				return layer
			}
		`
	assets := []uibridge.VideoRenderAsset{}
	for layerName, layerData := range data {
		asset := uibridge.VideoRenderAsset{
			Type:      layerData.Type,
			LayerName: layerName,
		}

		switch layerData.Type {
		case "image":
			if _, err := os.Stat(layerData.Data); os.IsNotExist(err) {
				w.handleError(errors.New("Не найден файл " + layerData.Data))
				return
			}
			script += fmt.Sprintf(`
					nexrender.selectLayersByName(null, '%s', function(layer) {
						var orig_layer = selectLayer(layer)
						var orig_w = orig_layer.width
						var orig_h = orig_layer.height

						nexrender.replaceFootage(orig_layer, '%s')

						var new_layer = selectLayer(layer)
						var new_w = new_layer.width
						var new_h = new_layer.height

						var w_scale = (orig_w / new_w) * 100
						var h_scale = (orig_h / new_h) * 100

						new_layer.transform.scale.setValue([w_scale, h_scale])
					});
				`, escapeInput(layerName), escapeInput(layerData.Data))
			continue

		case "data":
			script += fmt.Sprintf(`
					nexrender.selectLayersByName(null, '%s', function(layer) {
						selectLayer(layer).property("Source Text").setValue('%s')
					});
				`, escapeInput(layerName), escapeInput(layerData.Data))
			continue

		case "background-color":
			asset.Type = "data"
			hex, err := colors.ParseHEX(layerData.Data)
			if err != nil {
				w.handleError(err)
				return
			}

			rgb := hex.ToRGB()

			value := strings.Join([]string{
				strconv.FormatFloat(float64(rgb.R)/255, 'f', 2, 32),
				strconv.FormatFloat(float64(rgb.G)/255, 'f', 2, 32),
				strconv.FormatFloat(float64(rgb.B)/255, 'f', 2, 32),
				"1",
			}, ", ")

			property := ""
			if config.GetAELang() == "RU" {
				property = "Effects.Заливка.Цвет"
			} else {
				property = "Effects.Fill.Color"
			}

			script += fmt.Sprintf(`
					nexrender.selectLayersByName(null, '%s', function(layer) {
						var keys = "%s".split('.')
						var property = selectLayer(layer)
						for(var i = 0; i < keys.length; i++) {
							property = property.property(keys[i])
						}
						property.setValue([%s])
					});
				`, escapeInput(layerName), property, value)
			continue

		default:
			if _, err := os.Stat(layerData.Data); os.IsNotExist(err) {
				w.handleError(errors.New("Не найден файл " + layerData.Data))
				return
			}
			asset.Src = "file:///" + layerData.Data
			break
		}

		assets = append(assets, asset)
	}

	scriptPath, err := filepath.Abs(filepath.Join(paths.Temp, "images.jsx"))
	if err != nil {
		w.handleError(err)
		return
	}

	err = ioutil.WriteFile(scriptPath, []byte(script), 0666)
	if err != nil {
		w.handleError(err)
		return
	}
	defer os.Remove(scriptPath)

	assets = append(assets, uibridge.VideoRenderAsset{
		Type: "script",
		Src:  "file:///" + scriptPath,
	})

	renderInput := uibridge.VideoRenderConfig{
		Workpath:        paths.Temp,
		AepFile:         w.AepFile,
		CompositionName: w.CompositionName,
		Resolution:      w.Resolution,
		FromRenderQueue: w.FromRenderQueue,
		ScriptPath:      outputScriptPath,
		ResultPath:      w.ResultPath,
		OutputExt:       w.OutputExt,
		Assets:          assets,
		Memory:          config.GetAEMaxMemoryPercent(),
		AerenderPath:    config.GetAerenderPath(),
	}

	if w.isCancelled(true) {
		return
	}

	res, err := uibridge.Endpoint.Render(renderInput)

	if w.isCancelled(true) {
		return
	}

	if err != nil {
		w.handleError(err)
		return
	}

	logger.Notice("render complited, compress")

	projectName := strings.ReplaceAll(filepath.Base(w.AepFile), ".aep", "")
	resultName := filepath.Join(w.ResultPath, fmt.Sprintf("%s_%d", projectName, multiplier))

	if w.IsCompress {
		err = compressVideo(res.Output, resultName+w.Ext, w.Quality)
		if err != nil {
			w.handleError(err)
			return
		}
	} else {
		err = utils.MoveFile(res.Output, resultName+"."+w.OutputExt)
		if err != nil {
			w.handleError(err)
			return
		}
	}

	logFile := filepath.Join(paths.Temp, fmt.Sprintf("aerender-%s.log", res.Uid))

	err = utils.MoveFile(logFile, resultName+".log")
	if err != nil {
		w.handleError(err)
		return
	}

	logger.Notice("compress complited")

	err = os.RemoveAll(res.Workpath)
	if err != nil {
		w.handleError(err)
		return
	}
	w.progress = 100 + (100 * multiplier)
	w.handleProgressChanged("")
}

func escapeInput(str string) string {
	str = strings.ReplaceAll(str, "\\", "\\\\")
	str = strings.ReplaceAll(str, "'", "\\'")
	str = strings.ReplaceAll(str, "\r", "\\r")
	return strings.ReplaceAll(str, "\n", "\\n")
}

func stopAEProcess() {
	exec.Command("taskkill", "/F", "/IM", "AfterFX.com").Start()
}

func compressVideo(inputFile string, outputFile string, quality int) error {
	crf := 51 - int(0.51*float64(quality))
	cmd := exec.Command(config.GetFFmpegBin(), "-y", "-i", inputFile, "-c:v", "libx264", "-preset", "ultrafast", "-crf", strconv.Itoa(crf), outputFile)
	return cmd.Run()
}

const aeScript = `
var selectCompositionByName = function (name) {
    var len = app.project.items.length
    for (var i = 1; i <= len; i++) {
        var item = app.project.items[i]
        if (!(item instanceof CompItem)) continue;

        if (item.name === name) return item
    }

    throw new Error("nexrender: Can not find composition with name (" + name + ")")
}

var clearRenderQueue = function () {
    while (app.project.renderQueue.numItems > 0) {
        app.project.renderQueue.item(app.project.renderQueue.numItems).remove();
    }
}

var setResizeToOutputModule = function (module, x, y) {
    module.setSetting('Resize', 'true')
    module.setSetting('Resize to', '{x: ' + x + ', y: ' + y + '}')
}

clearRenderQueue()
app.project.renderQueue.items.add(selectCompositionByName(nexrender.renderCompositionName))
setResizeToOutputModule(app.project.renderQueue.item(1).outputModule(1), NX.get('x'), NX.get('y'))
`
