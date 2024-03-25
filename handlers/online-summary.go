package handlers

import (
	"scripts-api/controllers"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1
//
// Root godoc
//
//	@Summary	Retrieves online-summary script data for asset context
//	@Description
//
//	@Param			ts				query	string	true	"Timestamp range in format {startTs}:{endTs}. Timestamps in seconds, tz=UTC"	default(1708300800:1708387200)
//	@Param			organization	query	int		false	"Organization Integer Identifier"												default(213)
//	@Param			assetId			query	string	false	"Comma separated list of Asset IDs"												default(367,333)
//
//	@Tags			Data
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]map[string]any
//	@Router			/online-summary/data [get]
func OnlineSummary(c *controllers.ScriptDataController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := c.GetData()
		h(ctx)
	}
}
