package candleDao

import (
	"errors"
	"time"

	"github.com/paper-trade-chatbot/be-candle/models/dbModels"
	"github.com/paper-trade-chatbot/be-common/pagination"
	"github.com/paper-trade-chatbot/be-proto/general"

	"gorm.io/gorm"
)

const table = "candle"

type OrderColumn int

const (
	OrderColumn_None OrderColumn = iota
	OrderColumn_Start
	OrderColumn_ProductID
)

type OrderDirection int

const (
	OrderDirection_None = 0
	OrderDirection_ASC  = 1
	OrderDirection_DESC = -1
)

type Order struct {
	Column    OrderColumn
	Direction OrderDirection
}

// QueryModel set query condition, used by queryChain()
type QueryModel struct {
	ProductID    uint64
	IntervalType dbModels.IntervalType
	Start        *time.Time //只查單一時間點
	ProductIDIn  []uint64
	StartFrom    *time.Time
	StartTo      *time.Time
	OrderBy      []*Order
	Offset       int
	Limit        int
}

// New a row
func New(db *gorm.DB, model *dbModels.CandleModel) (int, error) {

	err := db.Table(table).
		Create(model).Error

	if err != nil {
		return 0, err
	}
	return 1, nil
}

// New rows
func News(db *gorm.DB, m []*dbModels.CandleModel) (int, error) {

	err := db.Transaction(func(tx *gorm.DB) error {

		err := tx.Table(table).
			CreateInBatches(m, 3000).Error

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return len(m), nil
}

// Get return a record as raw-data-form
func Get(tx *gorm.DB, query *QueryModel) (*dbModels.CandleModel, error) {

	result := &dbModels.CandleModel{}
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Gets return records as raw-data-form
func Gets(tx *gorm.DB, query *QueryModel) ([]dbModels.CandleModel, error) {
	result := make([]dbModels.CandleModel, 0)
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []dbModels.CandleModel{}, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetsWithPagination(tx *gorm.DB, query *QueryModel, paginate *general.Pagination) ([]dbModels.CandleModel, *general.PaginationInfo, error) {

	var rows []dbModels.CandleModel
	var count int64 = 0
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Count(&count).
		Scopes(paginateChain(paginate)).
		Scan(&rows).Error

	offset, _ := pagination.GetOffsetAndLimit(paginate)
	paginationInfo := pagination.SetPaginationDto(paginate.Page, paginate.PageSize, int32(count), int32(offset))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []dbModels.CandleModel{}, paginationInfo, nil
	}

	if err != nil {
		return []dbModels.CandleModel{}, nil, err
	}

	return rows, paginationInfo, nil
}

func queryChain(query *QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(productIDEqualScope(query.ProductID)).
			Scopes(intervalTypeEqualScope(query.IntervalType)).
			Scopes(startEqualScope(query.Start)).
			Scopes(productIDInScope(query.ProductIDIn)).
			Scopes(startBetweenScope(query.StartFrom, query.StartTo)).
			Scopes(orderByScope(query.OrderBy)).
			Scopes(offsetScope(query.Offset)).
			Scopes(limitScope(query.Limit))

	}
}

func paginateChain(paginate *general.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset, limit := pagination.GetOffsetAndLimit(paginate)
		return db.
			Scopes(offsetScope(offset)).
			Scopes(limitScope(limit))

	}
}

func productIDEqualScope(productID uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if productID != 0 {
			return db.Where(table+".product_id = ?", productID)
		}
		return db
	}
}

func intervalTypeEqualScope(intervalType dbModels.IntervalType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if intervalType != dbModels.IntervalType_None {
			return db.Where(table+".interval_type = ?", intervalType)
		}
		return db
	}
}

func startEqualScope(start *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if start != nil {
			return db.Where(table+".start = ?", start)
		}
		return db
	}
}

func productIDInScope(productIDIn []uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(productIDIn) > 0 {
			return db.Where(table+".product_id IN = ?", productIDIn)
		}
		return db
	}
}

func startBetweenScope(startFrom, startTo *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if startFrom != nil && startTo != nil {
			return db.Where(table+".start BEWTEEN ? AND ?", startFrom, startTo)
		}
		return db
	}
}

func orderByScope(order []*Order) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(order) > 0 {
			for _, o := range order {
				orderClause := ""
				switch o.Column {
				case OrderColumn_Start:
					orderClause += "start"
				default:
					continue
				}

				switch o.Direction {
				case OrderDirection_ASC:
					orderClause += " ASC"
				case OrderDirection_DESC:
					orderClause += " DESC"
				}

				db = db.Order(orderClause)
			}
			return db
		}
		return db
	}
}

func limitScope(limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit > 0 {
			return db.Limit(limit)
		}
		return db
	}
}

func offsetScope(offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if offset > 0 {
			return db.Limit(offset)
		}
		return db
	}
}
