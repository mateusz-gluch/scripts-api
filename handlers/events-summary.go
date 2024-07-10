package handlers

import (
	"scripts-api/controllers"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1
//
// Root godoc
//
//	@Summary	Retrieves script data for asset context
//
//	@Param			assets			query	string	true	"Comma separated list of Asset IDs"												default(367,333)
//	@Param			ts				query	string	false	"Timestamp range in format {startTs}:{endTs}. Timestamps in seconds, tz=UTC"	default(1708300800:1708387200)
//	@Param			span				query	string	false	"Span description string"	default(144h)
//	@Param			category		query	string	false	"Comma separated list of event categories"		example(MACHINE,DATA)
//
//	@Tags			Data
//	@Accept			json
//	@Produce		json
//	@Accept			csv
//	@Produce		csv
//	@Success		200	{object}	[]map[string]any
//	@Router			/events-summary/data [get]
func EventsSummary(c *controllers.ScriptDataController) gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
