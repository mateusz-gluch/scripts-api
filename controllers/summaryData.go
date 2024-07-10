package controllers

import (
	"fmt"
	repos "scripts-api/repositories"
	"time"

	"github.com/elmodis/go-libs/api"
	"github.com/elmodis/go-libs/formatters"
	"github.com/elmodis/go-libs/models"
	"github.com/elmodis/go-libs/models/properties"
	"github.com/elmodis/go-libs/models/specs"
	"github.com/elmodis/go-libs/parsers"
	"github.com/elmodis/go-libs/repositories"
	"github.com/elmodis/go-libs/validators"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func NewSummaryDataController[T any](
	data *repos.SummaryDataRepository[T],
	table string,
	assets repositories.ImmutableRepository[properties.Asset],
	filter map[string]parsers.Parser[[]string],
	log *zerolog.Logger,
) *SummaryDataController[T] {
	ret := new(SummaryDataController[T])
	ret.Label = "summary-data"
	ret.Log = log

	ret.Summary = data
	ret.AssetRepo = assets
	ret.Filter = filter
	ret.Table = table

	ret.Span = parsers.NewDurationParser("span", log)
	ret.Timestamp = parsers.NewUnboundTimestampParser("ts", log)
	ret.AssetParser = parsers.NewNumSequenceParser("assets", log)
	ret.OrganizationValid = validators.NewIdValidator(log)

	ret.json = formatters.NewJSONFormatter(*new(T), log)
	ret.csv = formatters.NewCSVFormatter(*new(T), log)

	return ret
}

type SummaryDataController[T any] struct {
	api.ControllerTemplate
	Summary           repositories.ImmutableSpecRepository[[]T, specs.ScriptSpec]
	AssetRepo         repositories.ImmutableRepository[properties.Asset]
	Filter            map[string]parsers.Parser[[]string]
	Timestamp         parsers.Parser[models.TimeRange]
	Span              parsers.Parser[time.Duration]
	AssetParser       parsers.Parser[[]int]
	OrganizationValid validators.Validator
	Table             string

	json formatters.Formatter[T]
	csv  formatters.Formatter[T]
}

func (c *SummaryDataController[T]) GetData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		op := "getdata"
		c.HandleDebug(op, ctx)
		var fmter formatters.Formatter[T]

		content := ctx.GetHeader("Accept")
		if content == "text/csv" {
			fmter = c.csv
		} else {
			fmter = c.json
		}

		spec, err := c.parseSpec(ctx)
		if err != nil {
			c.HandleBadRequest(op, ctx, "spec: %s", err.Error())
			return
		}

		ret, err := c.Summary.SelectSpec(c.Table, spec)
		if err != nil {
			// c.HandleServerError(op, ctx, "data: %s", err.Error())
			fmter.Format([]T{}, ctx)
			return
		}

		_, err = fmter.Format(*ret, ctx)
		if err != nil {
			c.HandleServerError(op, ctx, "format: %s", err.Error())
			return
		}
	}

}

func (c *SummaryDataController[T]) parseSpec(ctx *gin.Context) (*specs.ScriptSpec, error) {
	spec := new(specs.ScriptSpec)

	spanStr := ctx.Query("span")
	span, err := c.Span.Parse(spanStr)
	if err != nil {
		return nil, fmt.Errorf("span: %s", err.Error())
	}
	if span != nil {
		spec.StartTs = time.Now().Add(-*span)
		spec.EndTs = time.Now()
	}

	tsStr := ctx.Query("ts")
	ts, err := c.Timestamp.Parse(tsStr)
	if err != nil {
		return nil, fmt.Errorf("timestamp: %s", err.Error())
	}
	if ts != nil {
		spec.StartTs = ts.Start
		spec.EndTs = ts.End
	}

	if ts == nil && span == nil {
		return nil, fmt.Errorf("span and ts is nil")
	}

	assetStr := ctx.Query("assets")
	assets, err := c.AssetParser.Parse(assetStr)
	if err != nil {
		return nil, fmt.Errorf("asset: %s", err.Error())
	}
	spec.Assets = *assets

	return spec, nil
}
