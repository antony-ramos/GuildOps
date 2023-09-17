package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//	@BasePath	/api/v1

// GetPlayer godoc
//
//	@Summary		get player from database
//	@Schemes		http
//	@Description	Get player by name from database
//	@Tags			Player
//	@Accept			json
//	@Param			name	query	string	true	"player name"
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Router			/player/ [get]
func (api *API) GetPlayer(g *gin.Context) {
	playerName := g.Request.URL.Query().Get("name")
	if playerName == "" {
		g.JSON(http.StatusBadRequest, "missing name query param")
		return
	}
	player, err := api.ReadPlayer(g, playerName)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	g.JSON(http.StatusOK, player)
}

// DelPlayer godoc
//
//	@Summary		delete player from database
//	@Schemes		http
//	@Description	Delete player by name from database
//	@Tags			Player
//	@Accept			json
//	@Param			name	query	string	true	"player name"
//	@Produce		json
//	@Success		200	{string}	string	"OK"
//	@Router			/player/ [delete]
func (api *API) RemovePlayer(g *gin.Context) {
	playerName := g.Request.URL.Query().Get("name")
	if playerName == "" {
		g.JSON(http.StatusBadRequest, "missing name query param")
		return
	}
	err := api.DeletePlayer(g, playerName)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	g.JSON(http.StatusOK, "deleted")
}
