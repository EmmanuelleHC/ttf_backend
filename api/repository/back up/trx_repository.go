package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// TrxRepository database structure
type TrxRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewTrxRepository creates a new trx repository
func NewTrxRepository(db lib.Database, logger lib.Logger) TrxRepository {
	return TrxRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a TrxRepository) WithTrx(trxHandle *gorm.DB) TrxRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a TrxRepository) Query(param *models.TrxQueryParam) (*models.TrxQueryResult, error) {
	db := a.db.ORM.Model(&models.Trx{}).Preload("Company").Preload("Branch").Preload("CostCentre").Preload("Account").Preload("Department")

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ?", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.BranchID; v != "" {
		db = db.Where("branch_id=?", v)
	}

	if v := param.CCID; v != "" {
		db = db.Where("cc_id=?", v)
	}

	if v := param.AccID; v != "" {
		db = db.Where("account_id=?", v)
	}

	if v := param.DeptID; v != "" {
		db = db.Where("department_id=?", v)
	}

	if v := param.ActiveFlag; v {
		db = db.Where("active_flag=?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Trxs, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.TrxQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a TrxRepository) Get(id string) (*models.Trx, error) {
	trx := new(models.Trx)

	if ok, err := QueryOne(a.db.ORM.Model(trx).Where("id=?", id), trx); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return trx, nil
}

func (a TrxRepository) Create(trx *models.Trx) error {
	result := a.db.ORM.Model(trx).Create(trx)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a TrxRepository) Update(id string, trx *models.Trx) error {
	result := a.db.ORM.Model(trx).Where("id=?", id).Updates(trx)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a TrxRepository) Delete(id string) error {
	trx := new(models.Trx)

	result := a.db.ORM.Model(trx).Where("id=?", id).Delete(trx)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a TrxRepository) UpdateStatus(id string, status int) error {
	trx := new(models.Trx)

	result := a.db.ORM.Model(trx).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
