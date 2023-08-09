package server

import (
	"net/http"

	"github.com/Nicks344/moneytube/server/backend/src/modules/youtubeAPI"

	"github.com/gin-gonic/gin"
)

func serveYoutubeAPI(group *gin.RouterGroup) {
	youtubeGroup := group.Group("/youtube")
	{
		youtubeGroup.GET("/videoIDs/by/channelID", func(c *gin.Context) {
			var params struct {
				Apikey    string `form:"apikey"`
				ChannelID string `form:"channelID"`
			}

			if err := c.BindQuery(&params); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			result, err := youtubeAPI.GetVideoIDsByChannelID(params.Apikey, params.ChannelID, "")
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, result)
		})

		youtubeGroup.GET("/videoIDs/by/playlistID", func(c *gin.Context) {
			var params struct {
				Apikey     string `form:"apikey"`
				PlaylistID string `form:"playlistID"`
			}

			if err := c.BindQuery(&params); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			result, err := youtubeAPI.GetVideoIDsByPlaylistID(params.Apikey, params.PlaylistID, "")
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, result)
		})

		youtubeGroup.GET("/channelID/by/username", func(c *gin.Context) {
			var params struct {
				Apikey   string `form:"apikey"`
				Username string `form:"username"`
			}

			if err := c.BindQuery(&params); err != nil {
				fail(c, http.StatusBadRequest, err)
				return
			}

			result, err := youtubeAPI.GetChannelIDByUsername(params.Apikey, params.Username)
			if err != nil {
				fail(c, http.StatusInternalServerError, err)
				return
			}

			success(c, result)
		})
	}
}
