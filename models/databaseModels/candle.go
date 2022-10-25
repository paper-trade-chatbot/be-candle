package databaseModels

import (
	"time"

	"github.com/shopspring/decimal"
)

type IntervalType int32

const (
	IntervalType_None IntervalType = 0
	IntervalType_1MI  IntervalType = 21
	IntervalType_2MI  IntervalType = 22
	IntervalType_5MI  IntervalType = 25
	IntervalType_10MI IntervalType = 210
	IntervalType_15MI IntervalType = 215
	IntervalType_30MI IntervalType = 230
	IntervalType_1HR  IntervalType = 31
	IntervalType_1DY  IntervalType = 41
	IntervalType_5DY  IntervalType = 45
	IntervalType_1WK  IntervalType = 51
	IntervalType_1MO  IntervalType = 61
	IntervalType_1YR  IntervalType = 71
)

type CandleModel struct {
	ProductID    uint64          `gorm:"column:product_id"`
	IntervalType IntervalType    `gorm:"column:interval_type"`
	Start        time.Time       `gorm:"column:start"`
	Open         decimal.Decimal `gorm:"column:open"`
	Close        decimal.Decimal `gorm:"column:close"`
	High         decimal.Decimal `gorm:"column:high"`
	Low          decimal.Decimal `gorm:"column:low"`
	Volume       decimal.Decimal `gorm:"column:volume"`
}
