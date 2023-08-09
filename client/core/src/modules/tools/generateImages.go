package tools

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Nicks344/moneytube/client/core/src/paths"
)

type ImagesGeneratorWorker struct {
	*ToolService

	PsdFile    string
	ResultPath string
	Data       []map[string]LayerData
}

func (w *ImagesGeneratorWorker) Start() {
	w.maxProgress = 1
	w.handleProgressChanged("")

	d, _ := json.Marshal(w.Data)
	data := string(d)

	script := strings.ReplaceAll(psJSX, "{{PsdFile}}", strings.ReplaceAll(w.PsdFile, "\\", "/"))
	script = strings.ReplaceAll(script, "{{ResultPath}}", strings.ReplaceAll(w.ResultPath, "\\", "/"))
	script = strings.ReplaceAll(script, "{{Data}}", data)

	scriptFile := filepath.Join(paths.Temp, "script.jsx")
	os.Remove(scriptFile)
	ioutil.WriteFile(scriptFile, append([]byte{0xEF, 0xBB, 0xBF}, []byte(script)...), 0777)

	exec.Command("C:\\Windows\\SysWOW64\\wscript.exe", filepath.Join(paths.Bin, "ps.vbs"), scriptFile).Run()
	os.Remove(scriptFile)

	w.progress = 1
	w.handleProgressChanged("")
}

const psJSX = `
var doc = app.open(new File("{{PsdFile}}"));

var allLayers = collectAllLayers(doc, []);

function collectAllLayers(doc, allLayers) {
    for (var m = 0; m < doc.layers.length; m++) {
        var theLayer = doc.layers[m];
        if (theLayer.typename === "ArtLayer") {
            allLayers.push(theLayer);
        } else {
            collectAllLayers(theLayer, allLayers);
        }
    }
    return allLayers;
}

function findLayer(name) {
    for (var i = 0; i < allLayers.length; i++) {
        if (allLayers[i].name === name) {
            return allLayers[i]
        }
    }
}

function replaceImage(layer, newFile) {
    function changeSizeTo(layer, bounds) {
        var sourceWidth = bounds[2].value - bounds[0].value;
        var sourceHeight = bounds[3].value - bounds[1].value;
        var changeWidth = layer.bounds[2].value - layer.bounds[0].value;
        var changeHeight = layer.bounds[3].value - layer.bounds[1].value;
        var newWidth = sourceWidth / changeWidth * 100;
        var newHeight = sourceHeight / changeHeight * 100;
        layer.resize(newWidth, newHeight, AnchorPosition.MIDDLECENTER);
    }

    function moveLayerTo(layer, bounds) {
        var currBounds = layer.bounds;
        var x = bounds[0].value - currBounds[0].value;
        var y = bounds[1].value - currBounds[1].value;
      
        layer.translate(x, y);
    }

    doc.activeLayer = layer
    var bounds = doc.activeLayer.bounds
    var idplacedLayerReplaceContents = stringIDToTypeID("placedLayerReplaceContents");
    var desc3 = new ActionDescriptor();
    var idnull = charIDToTypeID("null");
    desc3.putPath(idnull, new File(newFile));
    var idPgNm = charIDToTypeID("PgNm");
    desc3.putInteger(idPgNm, 1);
    executeAction(idplacedLayerReplaceContents, desc3, DialogModes.NO);
    changeSizeTo(doc.activeLayer, bounds)
    moveLayerTo(doc.activeLayer, bounds)
};

function replaceText(layer, text) {
    layer.textItem.contents = text
}

var saveFolder = new Folder('{{ResultPath}}');
var fileName = doc.name.split('.')[0];

var jpgOptions = new JPEGSaveOptions();
jpgOptions.quality = 12;
jpgOptions.embedColorProfile = true;
jpgOptions.formatOptions = FormatOptions.PROGRESSIVE;
jpgOptions.scans = 5;
jpgOptions.matte = MatteType.NONE;

var data = {{Data}}
var layers = {}

for (var i = 0; i < data.length; i++) {
    var line = data[i]
    for (var layerName in line) {
        if (!layers[layerName]) {
            layers[layerName] = findLayer(layerName)
        }
        var layer = layers[layerName]
        var layerData = line[layerName]
        if (layerData.Type === 'data') {
            replaceText(layer, layerData.Data)
        } else if (layerData.Type === 'image') {
            replaceImage(layer, layerData.Data)
        }
    }
    doc.saveAs(new File(saveFolder + '/' + fileName + '_' + i + '.jpg'), jpgOptions, true, Extension.LOWERCASE);
}

doc.close(SaveOptions.DONOTSAVECHANGES)
`
