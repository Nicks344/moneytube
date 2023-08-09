package enigmakeygen

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GenerateKey(projectFileName, name, hwid string, days int) (string, error) {
	cmd := exec.Command("./data/keygen/keygen")
	env := fmt.Sprintf(`QUERY_STRING=Action=GenerateKeyFromProject&FileName=%s&RegName=%s&Hardware=%s&Days=%d`, projectFileName, name, hwid, days)
	cmd.Env = append(os.Environ(), env)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	outStr := strings.ReplaceAll(string(out), "\n", "")
	outStr = strings.ReplaceAll(outStr, "\r", "")
	outStr = strings.ReplaceAll(outStr, "Content-Type: text/html", "")
	return outStr, nil
}
