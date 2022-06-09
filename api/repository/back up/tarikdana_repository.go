package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// TarikDanaRepository database structure
type TarikDanaRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewTarikDanaRepository creates a new tarikdana repository
func NewTarikDanaRepository(db lib.Database, logger lib.Logger) TarikDanaRepository {
	return TarikDanaRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a TarikDanaRepository) WithTrx(trxHandle *gorm.DB) TarikDanaRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a TarikDanaRepository) Query(param *models.TarikDanaQueryParam) (*models.TarikDanaQueryResult, error) {
	db := a.db.ORM.Model(&models.TarikDana{}).Preload("Company").Preload("Branch")
	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Description; v != "" {
		v = "%" + v + "%"
		db = db.Where("description LIKE ?", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.BranchID; v != "" {
		db = db.Where("branch_id=?", v)
	}

	if v := param.Date; v != "" {
		db = db.Where("date=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("num LIKE ? OR description LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.TarikDanas, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.TarikDanaQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a TarikDanaRepository) Get(id string) (*models.TarikDana, error) {
	tarikdana := new(models.TarikDana)

	if ok, err := QueryOne(a.db.ORM.Model(tarikdana).Where("id=?", id), tarikdana); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return tarikdana, nil
}

func (a TarikDanaRepository) Create(tarikdana *models.TarikDana) error {
	result := a.db.ORM.Model(tarikdana).Create(tarikdana)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a TarikDanaRepository) Update(id string, tarikdana *models.TarikDana) error {
	result := a.db.ORM.Model(tarikdana).Where("id=?", id).Updates(tarikdana)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a TarikDanaRepository) Delete(id string) error {
	tarikdana := new(models.TarikDana)

	result := a.db.ORM.Model(tarikdana).Where("id=?", id).Delete(tarikdana)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a TarikDanaRepository) UpdateStatus(id string, status int) error {
	tarikdana := new(models.TarikDana)

	result := a.db.ORM.Model(tarikdana).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
