package server

import (
	"net/http"

	"github.com/Nicks344/moneytube/moneytubemodel"
	"github.com/Nicks344/moneytube/server/backend/src/model"
	"github.com/Nicks344/moneytube/server/backend/src/modules/macroses"

	"github.com/gin-gonic/gin"
)

func serveMacrosesAPI(group *gin.RouterGroup) {
	macrosGroup := group.Group("/macroses")
	{
		macrosGroup.POST("/execute", func(c *gin.Context) {
			key := c.GetHeader("Key")

			data, err := c.GetRawData()
			if err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			text := macroses.ExecuteUserMacroses(key, string(data))

			success(c, text)
		})

		macrosGroup.GET("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			macroses, err := model.GetMacroses(key)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, macroses)
		})

		macrosGroup.POST("/", func(c *gin.Context) {
			key := c.GetHeader("Key")

			var macros moneytubemodel.Macros

			if err := c.BindJSON(&macros); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			if err := model.SaveMacros(key, macros); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, nil)
		})

		macrosGroup.GET("/:name", func(c *gin.Context) {
			key := c.GetHeader("Key")
			name := c.Param("name")

			macros, err := model.GetMacros(key, name)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, macros)
		})

		macrosGroup.DELETE("/:name", func(c *gin.Context) {
			key := c.GetHeader("Key")
			name := c.Param("name")

			if err := model.DeleteMacros(key, name); err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, nil)
		})
	}
}
