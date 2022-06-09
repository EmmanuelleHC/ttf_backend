package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// InvoiceDetailRepository database structure
type InvoiceDetailRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewInvoiceDetailRepository creates a new invociedetail repository
func NewInvoiceDetailRepository(db lib.Database, logger lib.Logger) InvoiceDetailRepository {
	return InvoiceDetailRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a InvoiceDetailRepository) WithTrx(trxHandle *gorm.DB) InvoiceDetailRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a InvoiceDetailRepository) Query(param *models.InvoiceDetailQueryParam) (*models.InvoiceDetailQueryResult, error) {
	db := a.db.ORM.Model(&models.InvoiceDetail{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Status; v != "" {
		db = db.Where("status=?", v)
	}

	if v := param.InvoiceHeaderID; v != "" {
		db = db.Where("invoice_header_id=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("lines_desc LIKE ?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.InvoiceDetails, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.InvoiceDetailQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a InvoiceDetailRepository) Get(id string) (*models.InvoiceDetail, error) {
	invociedetail := new(models.InvoiceDetail)

	if ok, err := QueryOne(a.db.ORM.Model(invociedetail).Where("id=?", id), invociedetail); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return invociedetail, nil
}

func (a InvoiceDetailRepository) Create(invociedetail *models.InvoiceDetail) error {
	result := a.db.ORM.Model(invociedetail).Create(invociedetail)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceDetailRepository) CreateBatch(invociedetails models.InvoiceDetails) error {
	result := a.db.ORM.Model(invociedetails).Create(invociedetails)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceDetailRepository) Update(id string, invociedetail *models.InvoiceDetail) error {
	result := a.db.ORM.Model(invociedetail).Where("id=?", id).Updates(invociedetail)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceDetailRepository) Delete(id string) error {
	invociedetail := new(models.InvoiceDetail)

	result := a.db.ORM.Model(invociedetail).Where("id=?", id).Delete(invociedetail)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a InvoiceDetailRepository) UpdateStatus(id string, status int) error {
	invociedetail := new(models.InvoiceDetail)

	result := a.db.ORM.Model(invociedetail).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
