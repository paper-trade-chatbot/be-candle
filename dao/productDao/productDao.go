package productDao

import (
	"errors"

	models "github.com/paper-trade-chatbot/be-candle/models/databaseModels"
	"github.com/paper-trade-chatbot/be-common/pagination"
	"github.com/paper-trade-chatbot/be-proto/general"

	"gorm.io/gorm"
)

const table = "product"

// QueryModel set query condition, used by queryChain()
type QueryModel struct {
	ID            uint64
	ExchangeCode  string
	Code          string
	ProductType   []models.ProductType
	ExchangeCodes []string
	Status        int
	Display       int
	Offset        int
	Limit         int
}

// New a row
func New(tx *gorm.DB, model *models.ProductModel) (uint64, error) {

	err := tx.Table(table).
		Create(model).Error

	if err != nil {
		return 0, err
	}
	return model.ID, nil
}

// Get return a record as raw-data-form
func Get(tx *gorm.DB, query *QueryModel) (*models.ProductModel, error) {

	result := &models.ProductModel{}
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
func Gets(tx *gorm.DB, query *QueryModel) ([]models.ProductModel, error) {
	result := make([]models.ProductModel, 0)
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Scan(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.ProductModel{}, nil
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetsWithPagination(tx *gorm.DB, query *QueryModel, paginate *general.Pagination) ([]models.ProductModel, *general.PaginationInfo, error) {

	var rows []models.ProductModel
	var count int64 = 0
	err := tx.Table(table).
		Scopes(queryChain(query)).
		Count(&count).
		Scopes(paginateChain(paginate)).
		Scan(&rows).Error

	offset, _ := pagination.GetOffsetAndLimit(paginate)
	paginationInfo := pagination.SetPaginationDto(paginate.Page, paginate.PageSize, int32(count), int32(offset))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return []models.ProductModel{}, paginationInfo, nil
	}

	if err != nil {
		return []models.ProductModel{}, nil, err
	}

	return rows, paginationInfo, nil
}

func queryChain(query *QueryModel) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Scopes(idEqualScope(query.ID)).
			Scopes(codeEqualScope(query.Code)).
			Scopes(exchangeCodeEqualScope(query.ExchangeCode)).
			Scopes(productTypeInScope(query.ProductType)).
			Scopes(exchangeCodesInScope(query.ExchangeCodes)).
			Scopes(statusEqualScope(query.Status)).
			Scopes(displayEqualScope(query.Display)).
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

func idEqualScope(id uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if id != 0 {
			return db.Where(table+".id = ?", id)
		}
		return db
	}
}

func codeEqualScope(code string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if code != "" {
			return db.Where(table+".code = ?", code)
		}
		return db
	}
}

func exchangeCodeEqualScope(exchangeCode string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if exchangeCode != "" {
			return db.Where(table+".exchange_code = ?", exchangeCode)
		}
		return db
	}
}

func productTypeInScope(productType []models.ProductType) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(productType) > 0 {
			return db.Where(table+".type = ?", productType)
		}
		return db
	}
}

func exchangeCodesInScope(exchanges []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(exchanges) > 0 {
			return db.Where(table+".exchange_code = ?", exchanges)
		}
		return db
	}
}

func statusEqualScope(status int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if status != 0 {
			return db.Where(table+".status = ?", status)
		}
		return db
	}
}

func displayEqualScope(display int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if display != 0 {
			return db.Where(table+".display = ?", display)
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
