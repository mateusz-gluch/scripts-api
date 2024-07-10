package model

type OnlineSummary struct {
	ID           int    `gorm:"id" json:"-" mapstructure:"-"`
	Date         string `gorm:"date" json:"workingDay" mapstructure:"workingDay"`
	Tz           string `gorm:"tz" json:"timezone" mapstructure:"timezone"`
	AssetID      uint64 `gorm:"asset_id" json:"assetId" mapstructure:"assetId"`
	DatasourceID uint64 `gorm:"datasource_id" json:"dsId" mapstructure:"dsId"`
	Online       bool `gorm:"online" json:"online" mapstructure:"online"`
	Timestamp    string `gorm:"timestamp" json:"timestamp" mapstructure:"timestamp"`
}
