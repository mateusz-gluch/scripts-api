package controllers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/elmodis/go-libs/api/logging"
	"github.com/elmodis/go-libs/formatters"
	"github.com/elmodis/go-libs/models"
	"github.com/elmodis/go-libs/models/properties"
	"github.com/elmodis/go-libs/models/specs"
	"github.com/elmodis/go-libs/parsers"
	"github.com/elmodis/go-libs/repositories"
	"github.com/elmodis/go-libs/validators"
	"github.com/gin-gonic/gin"
)

type ScriptDataController struct {
	ScriptRepo        repositories.ImmutableSpecRepository[[]map[string]any, specs.ScriptSpec]
	AssetRepo         repositories.ImmutableRepository[properties.Asset]
	Filter            map[string]parsers.Parser[[]string]
	Timestamp         parsers.Parser[models.TimeRange]
	AssetParser       parsers.Parser[[]string]
	OrganizationValid validators.Validator
}

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
			logging.ControllerError("scriptData", "getData", ctx.Request.Method, ctx.Request.RequestURI, nil, err)
			ctx.String(400, err.Error())
			return
		}

		ret, err := c.ScriptRepo.SelectSpec(scriptName, spec)
		if err != nil {
			logging.ControllerError("scriptData", "getData", ctx.Request.Method, ctx.Request.RequestURI, nil, err)
			ctx.String(500, err.Error())
			return
		}

		_, err = fmter.Format(*ret, ctx)
		if err != nil {
			logging.ControllerError("scriptData", "getData", ctx.Request.Method, ctx.Request.RequestURI, nil, err)
			ctx.String(500, err.Error())
			return
		}
	}

}

func (c *ScriptDataController) parseSpec(ctx *gin.Context) (*specs.ScriptSpec, error) {
	spec := new(specs.ScriptSpec)

	tsStr := ctx.Query("ts")
	if tsStr == "" {
		tsStr = ctx.Param("ts")
	}

	ts, err := c.Timestamp.Parse(tsStr)
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
