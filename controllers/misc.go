package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MiscController struct {
	RootMessage string
	PingValue   int
	Version     string
}

//	@BasePath		/api/v1
//
// Root godoc
//
//	@Summary		Retrieves basic API information
//	@Description	Retrieves and prints API information. This information contains basic API description.
//	@Tags			Misc
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Router			/ [get]
func (c MiscController) Root() func(*gin.Context) {
	return func(ctx *gin.Context) {
		ret := make(map[string]any, 1)
		ret["message"] = c.RootMessage
		ctx.JSON(http.StatusOK, ret)
	}
}

//	@BasePath		/api/v1
//
// Ping godoc
//
//	@Summary		Pings API
//	@Description	Retrieves and prints ping message.
//	@Tags			Misc
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Router			/ping [get]
func (c MiscController) Ping() func(*gin.Context) {
	return func(ctx *gin.Context) {
		ret := make(map[string]any, 1)
		ret["pong"] = c.PingValue
		ctx.JSON(http.StatusOK, ret)
	}
}

//	@BasePath		/api/v1
//
// Version godoc
//
//	@Summary		Retrieves API Version
//	@Description	Retrieves and prints API version. This information contains API app version information
//	@Tags			Misc
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	string
//	@Router			/version [get]
func (c MiscController) Ver() func(*gin.Context) {
	return func(ctx *gin.Context) {
		ret := make(map[string]any, 1)
		ret["version"] = c.Version
		ctx.JSON(http.StatusOK, ret)
	}
}
