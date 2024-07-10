package model

type EventSummary struct {
	ID       int    `gorm:"id" json:"-" mapstructure:"-"`
	Date     string `gorm:"date" json:"workingDay" mapstructure:"workingDay"`
	Tz       string `gorm:"tz" json:"timezone" mapstructure:"timezone"`
	AssetID  uint64 `gorm:"asset_id" json:"assetId" mapstructure:"assetId"`
	Category string `gorm:"category" json:"category" mapstructure:"category"`
	Warnings uint64 `gorm:"warnings" json:"countWarnings" mapstructure:"countWarnings"`
	Alarms   uint64 `gorm:"alarms" json:"countAlarms" mapstructure:"countAlarms"`
	Online   string `gorm:"online" json:"online" mapstructure:"online"`
	Status   string `gorm:"status" json:"status" mapstructure:"status"`
}
