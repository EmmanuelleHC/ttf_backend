package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// BranchRepository database structure
type BranchRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewBranchRepository creates a new branch repository
func NewBranchRepository(db lib.Database, logger lib.Logger) BranchRepository {
	return BranchRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a BranchRepository) WithTrx(trxHandle *gorm.DB) BranchRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a BranchRepository) Query(param *models.BranchQueryParam) (*models.BranchQueryResult, error) {
	db := a.db.ORM.Model(&models.Branch{})

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

	list := make(models.Branchs, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.BranchQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a BranchRepository) Get(id string) (*models.Branch, error) {
	branch := new(models.Branch)

	if ok, err := QueryOne(a.db.ORM.Model(branch).Where("id=?", id), branch); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return branch, nil
}

func (a BranchRepository) GetFirst() (*models.Branch, error) {
	branch := new(models.Branch)

	if ok, err := QueryOne(a.db.ORM.Model(branch), branch); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return branch, nil
}

func (a BranchRepository) Create(branch *models.Branch) error {
	result := a.db.ORM.Model(branch).Create(branch)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BranchRepository) Update(id string, branch *models.Branch) error {
	result := a.db.ORM.Model(branch).Where("id=?", id).Updates(branch)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BranchRepository) Delete(id string) error {
	branch := new(models.Branch)

	result := a.db.ORM.Model(branch).Where("id=?", id).Delete(branch)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BranchRepository) UpdateStatus(id string, status int) error {
	branch := new(models.Branch)

	result := a.db.ORM.Model(branch).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
