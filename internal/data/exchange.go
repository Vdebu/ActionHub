package data

import (
	"gorm.io/gorm"
	"time"
)

// 存储汇率信息
type ExchangeRate struct {
	ID           uint      `json:"_id" gorm:"primarykey"`
	FromCurrency string    `json:"fromCurrency" binding:"required"`
	ToCurrency   string    `json:"toCurrency" binding:"required"`
	Rate         float64   `json:"rate" binding:"required"`
	Date         time.Time `json:"date"`
}

type ExchangeRateModel struct {
	db *gorm.DB
}

func (e ExchangeRate) TableName() string {
	return "exchange_rates"
}

// 自动迁移
func (m ExchangeRateModel) Migrate(exchangeRate *ExchangeRate) error {
	return m.db.AutoMigrate(exchangeRate)
}

// 创建记录
func (m ExchangeRateModel) Create(exchangeRate *ExchangeRate) error {
	return m.db.Create(exchangeRate).Error
}

// 返回最新的汇率信息
func (m ExchangeRateModel) GetLatest(exchangeRates *[]*ExchangeRate) error {
	return m.db.Find(exchangeRates).Error
}
