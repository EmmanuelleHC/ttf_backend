package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// SaldoMonthRepository database structure
type SaldoMonthRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewSaldoMonthRepository creates a new saldomonth repository
func NewSaldoMonthRepository(db lib.Database, logger lib.Logger) SaldoMonthRepository {
	return SaldoMonthRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a SaldoMonthRepository) WithTrx(trxHandle *gorm.DB) SaldoMonthRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a SaldoMonthRepository) Query(param *models.SaldoMonthQueryParam) (*models.SaldoMonthQueryResult, error) {
	db := a.db.ORM.Model(&models.SaldoMonth{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.BranchID; v != "" {
		db = db.Where("branch_id=?", v)
	}

	if v := param.MonthYear; v != "" {
		db = db.Where("month_year=?", v)
	}

	if v := param.ActiveFlag; v {
		db = db.Where("active_flag=?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.SaldoMonths, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.SaldoMonthQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a SaldoMonthRepository) Get(id string) (*models.SaldoMonth, error) {
	saldomonth := new(models.SaldoMonth)

	if ok, err := QueryOne(a.db.ORM.Model(saldomonth).Where("id=?", id), saldomonth); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldomonth, nil
}

func (a SaldoMonthRepository) GetbyCompanyAndBranch(companyId string, branchId string) (*models.SaldoMonth, error) {
	saldomonth := new(models.SaldoMonth)

	if ok, err := QueryOne(a.db.ORM.Model(saldomonth).Where("company_id=? AND branch_id=?", companyId, branchId), saldomonth); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldomonth, nil
}

func (a SaldoMonthRepository) GetbyCompanyAndBranchMonth(companyId string, branchId string, monthYear string) (*models.SaldoMonth, error) {
	saldomonth := new(models.SaldoMonth)

	if ok, err := QueryOne(a.db.ORM.Model(saldomonth).Where("company_id=? AND branch_id=? AND month_year=?", companyId, branchId, monthYear), saldomonth); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldomonth, nil
}

func (a SaldoMonthRepository) Create(saldomonth *models.SaldoMonth) error {
	result := a.db.ORM.Model(saldomonth).Create(saldomonth)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoMonthRepository) Update(id string, saldomonth *models.SaldoMonth) error {
	result := a.db.ORM.Model(saldomonth).Where("id=?", id).Updates(saldomonth)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoMonthRepository) Delete(id string) error {
	saldomonth := new(models.SaldoMonth)

	result := a.db.ORM.Model(saldomonth).Where("id=?", id).Delete(saldomonth)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoMonthRepository) UpdateStatus(id string, status int) error {
	saldomonth := new(models.SaldoMonth)

	result := a.db.ORM.Model(saldomonth).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
