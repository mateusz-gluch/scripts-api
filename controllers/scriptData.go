package controllers

import (
	"errors"
	"fmt"
	"scripts-api/formatters"
	"scripts-api/model"
	"strconv"

	"github.com/elmodis/go-libs/models"
	"github.com/elmodis/go-libs/models/properties"
	"github.com/elmodis/go-libs/parsers"
	"github.com/elmodis/go-libs/repositories"
	"github.com/elmodis/go-libs/validators"
	"github.com/gin-gonic/gin"
)

type ScriptDataController struct {
	ScriptRepo        repositories.ImmutableSpecRepository[[]map[string]any, model.ScriptSpec]
	AssetRepo         repositories.ImmutableRepository[properties.Asset]
	Filter            map[string]parsers.Parser[[]string]
	Timestamp         parsers.Parser[models.TimeRange]
	AssetParser       parsers.Parser[[]string]
	OrganizationValid validators.Validator
}

//	@BasePath	/api/v1
//
// Root godoc
//
//	@Summary	Retrieves script data for asset context
//	@Description
//	@Description	Currently available scripts
//	@Description	- events-summary - Summarizes daily state event occurences
//	@Description	- online-summary - Summarizes daily datasource activity
//
//	@Param			scriptName		path	string	true	"Unique Script Identifier"														enum(events-summary, online-summary)
//	@Param			ts				query	string	true	"Timestamp range in format {startTs}:{endTs}. Timestamps in seconds, tz=UTC"	default(1708300800:1708387200)
//	@Param			organization	query	int		false	"Organization Integer Identifier"												default(elmodis)
//	@Param			assetId			query	string	false	"Comma separated list of Asset IDs"												default(367)
//	@Param			category		query	string	false	"(available in events-summary) Comma separated list of events"					example(motor,motor.Psum)
//
//	@Tags			Assets
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	requests.EventsResponse
//	@Router			/{organization}/assets/{assetId} [get]
func (c *ScriptDataController) GetData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var fmter formatters.Formatter[map[string]any]

		scriptName := ctx.Param("scriptName")

		content := ctx.GetHeader("Accept")

		if content == "text/csv" {
			fmter = new(formatters.CSVFormatter[map[string]any])
		} else {
			fmter = new(formatters.JSONFormatter[map[string]any])
		}

		spec, err := c.parseSpec(ctx)
		if err != nil {
			ctx.String(400, err.Error())
			return
		}

		ret, err := c.ScriptRepo.SelectSpec(scriptName, spec)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		retStr, err := fmter.Format(*ret)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}
		ctx.String(200, retStr)
	}

}

func (c *ScriptDataController) parseSpec(ctx *gin.Context) (*model.ScriptSpec, error) {
	spec := new(model.ScriptSpec)

	ts, err := c.Timestamp.Parse(ctx.Query("ts"))
	if err != nil {
		return nil, fmt.Errorf("timestamp: %s", err.Error())
	}
	if ts == nil {
		return nil, fmt.Errorf("timestamp is mandatory")
	}
	spec.StartTs = ts.Start
	spec.EndTs = ts.End

	assetStr := ctx.Query("assets")
	assets, err := c.AssetParser.Parse(assetStr)
	if err != nil {
		return nil, fmt.Errorf("asset: %s", err.Error())
	}

	organization := ctx.Query("organization")
	err = c.OrganizationValid.Validate(organization)
	if err != nil && organization != "" {
		return nil, fmt.Errorf("organization: %s", err.Error())
	}

	if organization == "" {
		if len(*assets) == 0 {
			return nil, errors.New("insufficient data")
		}

		props, err := c.AssetRepo.Select((*assets)[0])
		if err != nil {
			return nil, fmt.Errorf("assets api: %s", err.Error())
		}
		spec.Organization = int(props.OrganizationId)

	} else {
		oid, _ := strconv.Atoi(organization)
		spec.Organization = oid
	}

	for _, asset := range *assets {
		aid, _ := strconv.Atoi(asset)
		spec.Assets = append(spec.Assets, aid)
	}

	queryParams := make(map[string]string)
	ctx.BindQuery(&queryParams)
	delete(queryParams, "ts")
	delete(queryParams, "assets")
	delete(queryParams, "organization")

	spec.Options = make(map[string][]string)

	for k, v := range queryParams {
		vals, err := c.Filter[k].Parse(v)
		if err != nil {
			return nil, fmt.Errorf("filter parser: %s", err.Error())
		}

		spec.Options[k] = *vals
	}

	return spec, nil
}
