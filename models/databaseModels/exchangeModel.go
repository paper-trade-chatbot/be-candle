package models

import (
	"database/sql"
	"time"
)

type ExchangeModel struct {
	Code                string        `gorm:"column:code; primary_key"`
	ProductType         ProductType   `gorm:"product_type"`
	Name                string        `gorm:"column:name"`
	Status              int           `gorm:"column:status"`          // 1:enabled , 2:disabled
	Display             int           `gorm:"column:display"`         // 1:enabled , 2:disabled
	CountryCode         string        `gorm:"column:country_code"`    //
	TimezoneOffset      float32       `gorm:"column:timezone_offset"` //
	OpenTime            sql.NullTime  `gorm:"column:open_time"`       //
	CloseTime           sql.NullTime  `gorm:"column:close_time"`      //
	ExchangeDay         string        `gorm:"column:exchange_day"`    // 星期幾
	ExceptionTime       string        `gorm:"column:exception_time"`
	DaylightSaving      bool          `gorm:"column:daylight_saving"`
	Location            string        `gorm:"column:location"`
	CreatedAt           time.Time     `gorm:"column:created_at"`
	UpdatedAt           time.Time     `gorm:"column:updated_at"`
	ExchangeDayParsed   ExchangeDay   `gorm:"-"`
	ExceptionTimeParsed ExceptionTime `gorm:"-"`
}

type ExchangeDay struct {
	StartDay int `json:"startDay"` // 星期幾開始
	EndDay   int `json:"endDay"`   // 星期幾結束（包含當日）
}

type ExceptionTime struct {
	Trade     []ExceptionTimeFormat `json:"trade"`
	StopTrade []ExceptionTimeFormat `json:"stopTrade"`
}

type ExceptionTimeFormat struct {
	Start time.Time `json:"start"` // 僅使用年以下之資料 ex:幾月幾日幾點
	End   time.Time `json:"end"`   // 僅使用年以下之資料 ex:幾月幾日幾點
}
