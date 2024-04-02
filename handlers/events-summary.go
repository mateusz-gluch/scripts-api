package handlers

import (
	"fmt"
	"scripts-api/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

//	@BasePath	/api/v1
//
// Root godoc
//
//	@Summary	Retrieves script data for asset context
//
//	@Param			ts				query	string	false	"Timestamp range in format {startTs}:{endTs}. Timestamps in seconds, tz=UTC"	default(1708300800:1708387200)
//	@Param			organization	query	int		false	"Organization Integer Identifier"												default(213)
//	@Param			assets			query	string	false	"Comma separated list of Asset IDs"												default(367,333)
//	@Param			category		query	string	false	"(available in events-summary) Comma separated list of event categories"		example(MACHINE,DATA)
//
//	@Tags			Data
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]map[string]any
//	@Router			/events-summary/data [get]
func EventsSummary(c *controllers.ScriptDataController) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if ctx.Query("ts") == "" {
			end := time.Now()
			delay := time.Duration(7 * 24 * time.Hour)
			start := end.Add(-delay)
			ctx.AddParam("ts", fmt.Sprintf("%d:%d", start.Unix(), end.Unix()))
		}

		ctx.AddParam("scriptName", "events-summary")
		h := c.GetData()
		h(ctx)
	}
}
