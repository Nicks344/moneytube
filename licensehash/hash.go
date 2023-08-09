package licensehash

import (
	"crypto/sha1"
	"encoding/hex"
)

func GetInfoHash(key, hwid string) string {
	return hash(key + hwid + "sdfsdfsdfvcxhjyu34423")
}

func hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
