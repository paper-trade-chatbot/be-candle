package models

import (
	"database/sql"
	"time"
)

type ExchangeHolidayType int

const (
	ExchangeHolidayType_None ExchangeHolidayType = iota
	ExchangeHolidayType_FullDay
	ExchangeHolidayType_HalfDay
)

type ExchangeHolidayModel struct {
	ID               uint64              `gorm:"column:id; primary_key"`
	Name             string              `gorm:"column:name"`
	Date             time.Time           `gorm:"column:date"`
	EndDate          sql.NullTime        `gorm:"column:end_date"`
	Type             ExchangeHolidayType `gorm:"type:enum('ExchangeHolidayType_None', 'ExchangeHolidayType_FullDay', 'ExchangeHolidayType_HalfDay')";"column:type"`
	ExchangeCode     string              `gorm:"column:exchange_code"`
	UpdatedAt        time.Time           `gorm:"column:updated_at"`
	HalfDayCloseTime sql.NullTime        `gorm:"column:half_day_close_time"`
	Memo             string              `gorm:"column:memo"`
}
