package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func serveUpdateAPI(router *gin.Engine) {
	apiGroup := router.Group("/api/update")
	{
		apiGroup.Use(userAuth())

		apiGroup.GET("/version", func(c *gin.Context) {
			version := c.GetString("version")
			newVersion, err := getLastVersion()
			major := getMajorVersion(version)
			newMajor := getMajorVersion(newVersion)

			if err != nil || major != newMajor {
				success(c, "0.0")
			} else {
				success(c, newVersion)
			}
		})

		apiGroup.GET("/download", func(c *gin.Context) {
			version := c.GetString("version")
			newVersion, err := getLastVersion()
			major := getMajorVersion(version)
			newMajor := getMajorVersion(newVersion)

			if err == nil && newVersion != "0.0" && major == newMajor {
				c.File("data/updates/" + newVersion)
			} else {
				fail(c, http.StatusNoContent, errors.New("no new version"))
			}
			/*
					ver, err := getLastVersion()
				if err != nil {
					fail(c, http.StatusNoContent, errors.New("no new version"))
					return
				}
				if ver != "0.0" {
					file, err := os.Open("data/bin/client/" + ver)
					if err != nil {
						fail(c, http.StatusInternalServerError, err)
						return
					}
					defer file.Close()
					fi, err := file.Stat()
					if err != nil {
						fail(c, http.StatusInternalServerError, err)
						return
					}
					w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))

				}
			*/
		})
	}
}

func getLastVersion() (string, error) {
	var lastTime time.Time
	lastVersion := "0.0"
	files, err := ioutil.ReadDir("data/updates")
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if f.ModTime().UnixNano() > lastTime.UnixNano() {
			lastTime = f.ModTime()
			lastVersion = f.Name()
		}
	}
	return lastVersion, nil
}

func getMajorVersion(version string) string {
	if len(version) == 0 {
		return ""
	}
	return version[:1]
}
