package upload

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Nicks344/moneytube/client/core/src/utils"

	"github.com/chromedp/cdproto/network"
	"github.com/gabriel-vasile/mimetype"
)

func getActualCookiesStr(cookies []*network.Cookie, setCookies []*http.Cookie) string {
	var cookiesStr string

	setCookiesMap := map[string]string{}
	if setCookies != nil {
		for _, cookie := range setCookies {
			setCookiesMap[cookie.Name] = cookie.Value
		}
	}

	for _, cookie := range cookies {
		if strings.Contains(cookie.Domain, "youtube") {
			value := cookie.Value
			if v, ok := setCookiesMap[cookie.Name]; ok {
				value = v
			}
			cookiesStr += cookie.Name + "=" + value + "; "
		}
	}

	return cookiesStr
}

func renameFile(path string, name string) (newPath string, err error) {
	var re = regexp.MustCompile(`(?m)[\+\=\[\]\:\;«\,\.\/\?\/\:\*\<\>\|]`)
	name = strings.Trim(re.ReplaceAllString(name, ""), "\r\n ")
	filename := filepath.Base(path)
	ext := filepath.Ext(filename)
	newPath = strings.ReplaceAll(path, filename, name+ext)
	err = os.Rename(path, newPath)
	return
}

func cutTags(tags string, maxSymbols int) string {
	count := func(t string) int {
		length := utf8.RuneCountInString(t)
		if strings.Contains(t, " ") {
			length += 2
		}
		return length
	}

	tagsSplit := strings.Split(tags, ", ")

	var length int
	for _, t := range tagsSplit {
		length += count(t)
	}
	length += len(tagsSplit) - 1

	for length > maxSymbols {
		i := len(tagsSplit) - 1
		lastTag := tagsSplit[i]
		l := count(lastTag)
		tagsSplit = utils.RemoveStr(tagsSplit, i)

		length -= l + 1
	}

	return strings.Join(tagsSplit, ", ")
}

func cutStringByWords(text string, maxSymbols int) string {
	var length int
	lines := strings.Split(text, "\n")

	linesWords := make([][]string, len(lines))

	for i, line := range lines {
		line = strings.Trim(line, "\r")

		length += utf8.RuneCountInString(line)

		linesWords[i] = strings.Split(line, " ")
	}

	if length <= maxSymbols {
		return text
	}

	for length > maxSymbols {
		i := len(linesWords) - 1
		wi := len(linesWords[i]) - 1
		word := linesWords[i][wi]
		wLength := utf8.RuneCountInString(word)

		if length-wLength < 0 {
			linesWords[i][wi] = string([]rune(word)[:maxSymbols])
			break
		}

		linesWords[i] = utils.RemoveStr(linesWords[i], wi)
		if len(linesWords[i]) == 0 {
			linesWords = append(linesWords[:i], linesWords[i+1:]...)
		}

		length -= wLength + 1
	}

	newLines := make([]string, len(linesWords))

	for i := range linesWords {
		newLines[i] = strings.Join(linesWords[i], " ")
	}

	newLines = utils.ClearSliceFromEmpty(newLines)

	return strings.Join(newLines, "\r\n")
}

func removeInvalidChars(str string) string {
	var re = regexp.MustCompile(`(?m)[«\<\>]`)
	return strings.Trim(re.ReplaceAllString(str, ""), "\r\n ")
}

func waitForTrue(ctx context.Context, variable *bool) error {
	for !*variable {
		if utils.IsContextCancelled(ctx) {
			return errors.New("timeout")
		}

		time.Sleep(time.Millisecond * 100)
	}

	return nil
}

func getFileDataURI(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	m, err := mimetype.DetectFile(filename)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:%s;base64,%s", m.String(), base64.StdEncoding.EncodeToString(data)), nil
}
