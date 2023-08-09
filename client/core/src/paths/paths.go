package paths

import (
	"os"
	"path/filepath"
)

const Data = "data"
const Bin = "bin"

var Sessions = filepath.Join(Data, "sessions")
var Templates = filepath.Join(Data, "templates")
var AETemplates = filepath.Join(Templates, "ae")
var UniqueTemplates = filepath.Join(Templates, "unique")
var Fonts = filepath.Join(Data, "fonts")
var Temp = filepath.Join(os.TempDir(), "moneytube")
var ChromeExe = filepath.Join(Bin, "chrome-win", "chrome.exe")
