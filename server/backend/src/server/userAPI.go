package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func serveUserAPI(router *gin.Engine) {
	apiGroup := router.Group("/api/user")
	{
		apiGroup.Use(userAuth())

		v1Group := apiGroup.Group("/v1")
		{
			v1Group.GET("/check", func(c *gin.Context) { c.Status(http.StatusOK) })

			serveMacrosesAPI(v1Group)

			serveYoutubeAPI(v1Group)

			serveAudioAPI(v1Group)

			serveAccountsAPI(v1Group)

			serveUploadDatasAPI(v1Group)

			serveUploadTasksAPI(v1Group)

			serveBugReportAPI(v1Group)

			serveUploadDataTemplatesAPI(v1Group)
		}
	}
}
