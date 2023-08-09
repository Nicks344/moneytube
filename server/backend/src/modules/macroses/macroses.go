package macroses

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Nicks344/moneytube/client/core/src/utils"
	"github.com/Nicks344/moneytube/server/backend/src/model"
)

func tagCharLength(t string) int {
	length := utf8.RuneCountInString(t)
	if strings.Contains(t, " ") {
		length += 2
	}
	return length
}

func charLength(t string) int {
	return utf8.RuneCountInString(t)
}

func ExecuteUserMacroses(key string, text string) string {
	macrosesInText := []string{}
	var re = regexp.MustCompile(`(?m)\[(\d+-\d+|t?less-\d+)\:.*?\]`)

	for _, match := range re.FindAllString(text, -1) {
		macrosesInText = append(macrosesInText, match)
	}

	for _, macros := range macrosesInText {
		values := []string{}

		mTrim := strings.Trim(macros, "[]")
		mSplit := strings.Split(mTrim, ":")
		mCount := mSplit[0]
		mName := mSplit[1]

		macrosData, err := model.GetMacros(key, mName)
		if err != nil {
			continue
		}

		if strings.Contains(mCount, "less") {
			mCountSplit := strings.Split(mCount, "-")
			isTags := mCountSplit[0] == "tless"
			charCount, err := strconv.Atoi(mCountSplit[1])
			if err != nil {
				continue
			}

			var length int

			for {
				val := macrosData.Data[utils.RandRange(0, len(macrosData.Data)-1)]

				if isTags {
					length += tagCharLength(val) + 1
				} else {
					length += charLength(val) + 2
				}

				if length >= charCount {
					break
				}

				values = append(values, val)
			}

		} else {
			mRangeStr := strings.Split(mCount, "-")

			mValues := macrosData.Data

			mRange := make([]int64, 2, 2)
			mRange[0], err = strconv.ParseInt(mRangeStr[0], 10, 32)
			if err != nil {
				continue
			}
			mRange[1], err = strconv.ParseInt(mRangeStr[1], 10, 32)
			if err != nil {
				continue
			}

			count := utils.RandRange(int(mRange[0]), int(mRange[1]))
			if count > len(mValues) {
				count = len(mValues)
				values = mValues
			} else {
				for i := 0; i < count; i++ {
					m := utils.RandRange(0, len(mValues)-1)
					values = append(values, mValues[m])
					mValues = utils.RemoveStr(mValues, m)
				}
			}
		}

		value := strings.Join(values, ", ")
		text = strings.Replace(text, macros, value, 1)
	}

	return text
}
