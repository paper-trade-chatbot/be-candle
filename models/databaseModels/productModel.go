package models

import (
	"database/sql"
	"time"
)

type ProductType int

const (
	ProductType_None ProductType = iota
	ProductType_Stock
	ProductType_Crypto
	ProductType_Forex
	ProductType_Futures
)

type ProductModel struct {
	ID           uint64          `gorm:"column:id; primary_key"`
	Type         ProductType     `gorm:"column:type"`
	ExchangeCode string          `gorm:"column:exchange_code"`
	Code         string          `gorm:"column:code"`
	Name         string          `gorm:"column:name"`
	Status       int             `gorm:"column:status"`  // 1:enabled , 2:disabled
	Display      int             `gorm:"column:display"` // 1:enabled , 2:disabled
	CurrencyCode string          `gorm:"column:currency_code"`
	TickUnit     float64         `gorm:"column:tick_unit"`
	MinimumOrder sql.NullFloat64 `gorm:"column:minimum_order"`
	IconID       sql.NullString  `gorm:"column:icon_id"`
	CreatedAt    time.Time       `gorm:"column:created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at"`
}
