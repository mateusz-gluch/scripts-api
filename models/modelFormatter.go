package model

import (
	"time"
)

type ModelFormatter[T any] interface {
	Format(e T) T
}

type EventFormatter struct{}

func (EventFormatter) Format(e EventSummary) EventSummary {
	dt, _ := time.Parse(time.RFC3339, e.Date)

	return EventSummary{
		Date:     dt.Format("2006-01-02"),
		Tz:       e.Tz,
		AssetID:  e.AssetID,
		Category: e.Category,
		Warnings: e.Warnings,
		Alarms:   e.Alarms,
		Online:   e.Online,
		Status:   e.Status,
	}
}

type OnlineFormatter struct{}

func (OnlineFormatter) Format(e OnlineSummary) OnlineSummary {
	dt, _ := time.Parse(time.RFC3339, e.Date)

	return OnlineSummary{
		Date:         dt.Format("2006-01-02"),
		Tz:           e.Tz,
		AssetID:      e.AssetID,
		DatasourceID: e.DatasourceID,
		Online:       e.Online,
		Timestamp:    e.Timestamp,
	}
}
