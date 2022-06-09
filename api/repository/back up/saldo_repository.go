package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// SaldoRepository database structure
type SaldoRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewSaldoRepository creates a new saldo repository
func NewSaldoRepository(db lib.Database, logger lib.Logger) SaldoRepository {
	return SaldoRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a SaldoRepository) WithTrx(trxHandle *gorm.DB) SaldoRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a SaldoRepository) Query(param *models.SaldoQueryParam) (*models.SaldoQueryResult, error) {
	db := a.db.ORM.Model(&models.Saldo{})

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

	list := make(models.Saldos, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.SaldoQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a SaldoRepository) Get(id string) (*models.Saldo, error) {
	saldo := new(models.Saldo)

	if ok, err := QueryOne(a.db.ORM.Model(saldo).Where("id=?", id), saldo); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldo, nil
}

func (a SaldoRepository) GetRoot() (*models.Saldo, error) {
	saldo := new(models.Saldo)

	if ok, err := QueryOne(a.db.ORM.Model(saldo), saldo); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldo, nil
}

func (a SaldoRepository) GetbyCompanyAndBranch(companyId string, branchId string) (*models.Saldo, error) {
	saldo := new(models.Saldo)

	if ok, err := QueryOne(a.db.ORM.Model(saldo).Where("company_id=? AND branch_id=?", companyId, branchId), saldo); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldo, nil
}

func (a SaldoRepository) GetbyCompanyAndBranchMonth(companyId string, branchId string, monthYear string) (*models.Saldo, error) {
	saldo := new(models.Saldo)

	if ok, err := QueryOne(a.db.ORM.Model(saldo).Where("company_id=? AND branch_id=? AND month_year=?", companyId, branchId, monthYear), saldo); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldo, nil
}

func (a SaldoRepository) Create(saldo *models.Saldo) error {
	result := a.db.ORM.Model(saldo).Create(saldo)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoRepository) Update(id string, saldo *models.Saldo) error {
	result := a.db.ORM.Model(saldo).Where("id=?", id).Updates(saldo)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoRepository) Delete(id string) error {
	saldo := new(models.Saldo)

	result := a.db.ORM.Model(saldo).Where("id=?", id).Delete(saldo)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoRepository) UpdateStatus(id string, status int) error {
	saldo := new(models.Saldo)

	result := a.db.ORM.Model(saldo).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
