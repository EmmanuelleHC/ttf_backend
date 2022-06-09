package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// CompanyRepository database structure
type CompanyRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewCompanyRepository creates a new company repository
func NewCompanyRepository(db lib.Database, logger lib.Logger) CompanyRepository {
	return CompanyRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a CompanyRepository) WithTrx(trxHandle *gorm.DB) CompanyRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a CompanyRepository) Query(param *models.CompanyQueryParam) (*models.CompanyQueryResult, error) {
	db := a.db.ORM.Model(&models.Company{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		q := "%" + v + "%"
		db = db.Where("name LIKE ?", q)
	}

	if v := param.Address; v != "" {
		db = db.Where("address LIKE ?", v)
	}

	if v := param.Num; v != "" {
		db = db.Where("num LIKE ?", v)
	}

	if v := param.PaymentFlag; v {
		db = db.Where("payment_flag=?", v)
	}

	if v := param.BDCFlag; v {
		db = db.Where("bdc_flag=?", v)
	}

	if v := param.ApprovalFlag; v {
		db = db.Where("approval_flag=?", v)
	}

	if v := param.ActiveFlag; v {
		db = db.Where("active_flag=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR address LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Companies, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.CompanyQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a CompanyRepository) Get(id string) (*models.Company, error) {
	company := new(models.Company)

	if ok, err := QueryOne(a.db.ORM.Model(company).Where("id=?", id), company); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return company, nil
}

func (a CompanyRepository) Create(company *models.Company) error {
	result := a.db.ORM.Model(company).Create(company)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CompanyRepository) Update(id string, company *models.Company) error {
	result := a.db.ORM.Model(company).Where("id=?", id).Updates(company)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CompanyRepository) Delete(id string) error {
	company := new(models.Company)

	result := a.db.ORM.Model(company).Where("id=?", id).Delete(company)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CompanyRepository) UpdateStatus(id string, status int) error {
	company := new(models.Company)

	result := a.db.ORM.Model(company).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
