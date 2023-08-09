package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Nicks344/moneytube/server/backend/src/model"
	"github.com/Nicks344/moneytube/server/backend/src/modules/users"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func serveAdminAPI(router *gin.Engine, debug bool) {
	var authorized *gin.RouterGroup

	if !debug {
		authHandler := gin.BasicAuth(gin.Accounts{
			"moneytube-admin": "S1xcOgCnAbio89L3blbj",
		})
		authorized = router.Group("/", authHandler)
	} else {
		authorized = router.Group("/")
	}

	router.Use(static.Serve("/", static.LocalFile("data/ui", true)))

	adminApiGroup := authorized.Group("/api/admin")
	{
		adminApiGroup.GET("/users", func(c *gin.Context) {
			u, err := model.GetUsers()
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
			} else {
				success(c, u)
			}
		})

		adminApiGroup.POST("/user", func(c *gin.Context) {
			var data struct {
				Name           string
				Days           int
				DaysReactivate int
				Version        string
				Count          int
			}
			if err := c.ShouldBindJSON(&data); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			result := []model.User{}
			for i := 0; i < data.Count; i++ {
				name := data.Name
				if data.Count > 1 {
					name = fmt.Sprintf("%s-%d", data.Name, i)
				}
				u, err := users.CreateUser(name, data.Days, data.DaysReactivate, data.Version)
				if err != nil {
					fail(c, http.StatusInternalServerError, err)
					return
				} else {
					result = append(result, u)
				}
			}

			success(c, result)
		})

		adminApiGroup.DELETE("/user/:id", func(c *gin.Context) {
			key := c.Param("id")
			if err := model.DeleteUser(key); err != nil {
				fail(c, http.StatusInternalServerError, err)
			} else {
				success(c, nil)
			}
		})

		adminApiGroup.GET("/bugreports", func(c *gin.Context) {
			reports, err := model.GetReports()
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
			} else {
				success(c, reports)
			}
		})

		adminApiGroup.GET("/bugreports/:id", func(c *gin.Context) {
			c.File(fmt.Sprintf("data/reports/%s.zip", c.Param("id")))
		})

		adminApiGroup.DELETE("/bugreports/:id", func(c *gin.Context) {
			id := c.Param("id")
			if err := model.DeleteReport(id); err != nil {
				fail(c, http.StatusInternalServerError, err)
			} else {
				os.Remove(fmt.Sprintf("data/reports/%s.zip", id))
				success(c, nil)
			}
		})

		adminApiGroup.POST("/bugreports/:id/resolve", func(c *gin.Context) {
			id := c.Param("id")
			if err := model.ResolveReport(id); err != nil {
				fail(c, http.StatusInternalServerError, err)
			} else {
				success(c, nil)
			}
		})
	}
}
