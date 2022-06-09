package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// CostCentreRepository database structure
type CostCentreRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewCostCentreRepository creates a new costcentre repository
func NewCostCentreRepository(db lib.Database, logger lib.Logger) CostCentreRepository {
	return CostCentreRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a CostCentreRepository) WithTrx(trxHandle *gorm.DB) CostCentreRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a CostCentreRepository) Query(param *models.CostCentreQueryParam) (*models.CostCentreQueryResult, error) {
	db := a.db.ORM.Model(&models.CostCentre{}).Preload("Company")

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ?", v)
	}

	if v := param.Code; v != "" {
		db = db.Where("code LIKE ?", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.ActiveFlag; v {
		db = db.Where("active_flag=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("name LIKE ? OR code LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.CostCentres, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.CostCentreQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a CostCentreRepository) Get(id string) (*models.CostCentre, error) {
	costcentre := new(models.CostCentre)

	if ok, err := QueryOne(a.db.ORM.Model(costcentre).Where("id=?", id), costcentre); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return costcentre, nil
}

func (a CostCentreRepository) Create(costcentre *models.CostCentre) error {
	result := a.db.ORM.Model(costcentre).Create(costcentre)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CostCentreRepository) Update(id string, costcentre *models.CostCentre) error {
	result := a.db.ORM.Model(costcentre).Where("id=?", id).Updates(costcentre)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CostCentreRepository) Delete(id string) error {
	costcentre := new(models.CostCentre)

	result := a.db.ORM.Model(costcentre).Where("id=?", id).Delete(costcentre)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CostCentreRepository) UpdateStatus(id string, status int) error {
	costcentre := new(models.CostCentre)

	result := a.db.ORM.Model(costcentre).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
