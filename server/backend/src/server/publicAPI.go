package server

import (
	"net/http"

	"github.com/Nicks344/moneytube/server/backend/src/modules/users"

	"github.com/gin-gonic/gin"
)

func servePublicAPI(router *gin.Engine) {
	apiGroup := router.Group("/api/public")
	{
		apiGroup.POST("/activate", func(c *gin.Context) {
			var data struct {
				Key  string
				HWID string
			}
			if err := c.ShouldBindJSON(&data); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			if key, name, err := users.Activate(data.Key, data.HWID); err != nil {
				fail(c, http.StatusInternalServerError, err)
			} else {
				success(c, map[string]string{
					"key":  key,
					"name": name,
				})
			}
		})
	}
}
