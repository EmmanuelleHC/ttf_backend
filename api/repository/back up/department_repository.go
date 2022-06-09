package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// DepartmentRepository database structure
type DepartmentRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewDepartmentRepository creates a new departmentrepository
func NewDepartmentRepository(db lib.Database, logger lib.Logger) DepartmentRepository {
	return DepartmentRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a DepartmentRepository) WithTrx(trxHandle *gorm.DB) DepartmentRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a DepartmentRepository) Query(param *models.DepartmentQueryParam) (*models.DepartmentQueryResult, error) {
	db := a.db.ORM.Model(&models.Department{})

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

	if v := param.Desc; v != "" {
		db = db.Where("desc LIKE ?", v)
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

	list := make(models.Departments, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.DepartmentQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a DepartmentRepository) Get(id string) (*models.Department, error) {
	department := new(models.Department)

	if ok, err := QueryOne(a.db.ORM.Model(department).Where("id=?", id), department); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return department, nil
}

func (a DepartmentRepository) Create(department *models.Department) error {
	result := a.db.ORM.Model(department).Create(department)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a DepartmentRepository) Update(id string, department *models.Department) error {
	result := a.db.ORM.Model(department).Where("id=?", id).Updates(department)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a DepartmentRepository) Delete(id string) error {
	department := new(models.Department)

	result := a.db.ORM.Model(department).Where("id=?", id).Delete(department)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a DepartmentRepository) UpdateStatus(id string, status int) error {
	department := new(models.Department)

	result := a.db.ORM.Model(department).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
