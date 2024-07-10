package repos

import (
	"errors"
	model "scripts-api/models"

	"github.com/elmodis/go-libs/models/specs"
	"github.com/elmodis/go-libs/repositories"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const dateFmt = "2006-01-02"

func NewSummaryDataRepository[T any](t T, fmt model.ModelFormatter[T], db *gorm.DB, log *zerolog.Logger, cat bool) *SummaryDataRepository[T] {
	ret := new(SummaryDataRepository[T])
	ret.Label = "summary_data"
	ret.Log = log
	ret.Engine = db
	ret.Formatter = fmt
	ret.categoryFlag = cat
	return ret
}

type SummaryDataRepository[T any] struct {
	repositories.SpecRepositoryTemplate[[]T, specs.ScriptSpec]
	Engine       *gorm.DB
	Formatter    model.ModelFormatter[T]
	categoryFlag bool
}

func (repo *SummaryDataRepository[T]) SelectSpec(table string, spec *specs.ScriptSpec) (*[]T, error) {
	op := "select"
	repo.HandleDebug(table, spec, op)

	ret := new([]T)

	q := repo.Engine.Table(table)

	startDate := spec.StartTs.Format(dateFmt)
	endDate := spec.EndTs.Format(dateFmt)

	q = q.Where("asset_id", spec.Assets)
	q = q.Where("date BETWEEN ? AND ?", startDate, endDate)

	if repo.categoryFlag {
		q = q.Order("date DESC, category")
	} else {
		q = q.Order("date DESC")
	}

	err := q.Find(ret).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return ret, nil
	}

	for idx, e := range *ret {
		(*ret)[idx] = repo.Formatter.Format(e)
	}

	if err != nil {
		return nil, repo.HandleError(table, spec, op, "db: %s", err.Error())
	}

	return ret, nil
}
