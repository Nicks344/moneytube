package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/Nicks344/moneytube/server/backend/src/modules/audio"

	"github.com/gin-gonic/gin"
)

func serveAudioAPI(group *gin.RouterGroup) {
	audioGroup := group.Group("/audio")
	{
		audioGroup.POST("/generate/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			var data moneytubemodel.AudioGenerateInput

			if err := c.BindJSON(&data); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			resultFile := fmt.Sprintf("./data/temp/%s-audio-%d.mp3", key, time.Now().UnixNano())

			if err := audio.Generate(data, resultFile); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			c.File(resultFile)
			os.Remove(resultFile)
		})
	}
}

func serveAudioAPI2(group *gin.RouterGroup) {
	audioGroup := group.Group("/audio")
	{
		audioGroup.POST("/generate/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			var data moneytubemodel.AudioGenerateInput

			if err := c.BindJSON(&data); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			resultFile := fmt.Sprintf("./data/temp/%s-audio-%d.mp3", key, time.Now().UnixNano())

			if err := audio.Generate(data, resultFile); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			c.File(resultFile)
			os.Remove(resultFile)
		})
	}
}
