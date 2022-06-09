package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// AccountRepository database structure
type AccountRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db lib.Database, logger lib.Logger) AccountRepository {
	return AccountRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a AccountRepository) WithTrx(trxHandle *gorm.DB) AccountRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a AccountRepository) Query(param *models.AccountQueryParam) (*models.AccountQueryResult, error) {
	db := a.db.ORM.Model(&models.Account{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ?", v)
	}

	if v := param.Num; v != "" {
		db = db.Where("num LIKE ?", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.ActiveFlag; v {
		db = db.Where("active_flag=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR num LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Accounts, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.AccountQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a AccountRepository) Get(id string) (*models.Account, error) {
	account := new(models.Account)

	if ok, err := QueryOne(a.db.ORM.Model(account).Where("id=?", id), account); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return account, nil
}

func (a AccountRepository) Create(account *models.Account) error {
	result := a.db.ORM.Model(account).Create(account)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a AccountRepository) Update(id string, account *models.Account) error {
	result := a.db.ORM.Model(account).Where("id=?", id).Updates(account)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a AccountRepository) Delete(id string) error {
	account := new(models.Account)

	result := a.db.ORM.Model(account).Where("id=?", id).Delete(account)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a AccountRepository) UpdateStatus(id string, status int) error {
	account := new(models.Account)

	result := a.db.ORM.Model(account).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
