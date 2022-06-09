package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// InvoiceHeaderRepository database structure
type InvoiceHeaderRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewInvoiceHeaderRepository creates a new invoiceheader repository
func NewInvoiceHeaderRepository(db lib.Database, logger lib.Logger) InvoiceHeaderRepository {
	return InvoiceHeaderRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a InvoiceHeaderRepository) WithTrx(trxHandle *gorm.DB) InvoiceHeaderRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a InvoiceHeaderRepository) Query(param *models.InvoiceHeaderQueryParam) (*models.InvoiceHeaderQueryResult, error) {
	db := a.db.ORM.Model(&models.InvoiceHeader{}).Preload("Company").Preload("Branch")

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Num; v != "" {
		v = "%" + v + "%"
		db = db.Where("num LIKE ?", v)
	}

	if v := param.Status; v != "" {
		db = db.Where("status=?", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.BranchID; v != "" {
		db = db.Where("branch_id=?", v)
	}

	if v := param.StatusApprove; v != "" {
		db = db.Where("status_approve=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("num LIKE ?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.InvoiceHeaders, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.InvoiceHeaderQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a InvoiceHeaderRepository) Get(id string) (*models.InvoiceHeader, error) {
	invoiceheader := new(models.InvoiceHeader)

	if ok, err := QueryOne(a.db.ORM.Model(invoiceheader).Where("id=?", id), invoiceheader); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return invoiceheader, nil
}

func (a InvoiceHeaderRepository) Create(invoiceheader *models.InvoiceHeader) error {
	result := a.db.ORM.Model(invoiceheader).Omit("Branch").Omit("Company").Create(invoiceheader)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceHeaderRepository) Update(id string, invoiceheader *models.InvoiceHeader) error {
	result := a.db.ORM.Model(invoiceheader).Where("id=?", id).Updates(invoiceheader)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceHeaderRepository) Delete(id string) error {
	invoiceheader := new(models.InvoiceHeader)

	result := a.db.ORM.Model(invoiceheader).Where("id=?", id).Delete(invoiceheader)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceHeaderRepository) UpdateStatus(id string, status int) error {
	invoiceheader := new(models.InvoiceHeader)

	result := a.db.ORM.Model(invoiceheader).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceHeaderRepository) UpdateApprove(id string, status int) error {
	invoiceheader := new(models.InvoiceHeader)

	result := a.db.ORM.Model(invoiceheader).Where("id=?", id).Update("status_approve", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceHeaderRepository) UpdateApproves(ids []string, status int) error {
	invoiceheader := new(models.InvoiceHeader)

	result := a.db.ORM.Model(invoiceheader).Where("id IN (?)", ids).Update("status_approve", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
