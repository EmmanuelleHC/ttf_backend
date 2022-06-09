package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// EmployeeRepository database structure
type EmployeeRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewEmployeeRepository creates a new employee repository
func NewEmployeeRepository(db lib.Database, logger lib.Logger) EmployeeRepository {
	return EmployeeRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a EmployeeRepository) WithTrx(trxHandle *gorm.DB) EmployeeRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a EmployeeRepository) Query(param *models.EmployeeQueryParam) (*models.EmployeeQueryResult, error) {
	db := a.db.ORM.Model(&models.Employee{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Name; v != "" {
		q := "%" + v + "%"
		db = db.Where("name LIKE ?", q)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.BranchID; v != "" {
		db = db.Where("branch_id=?", v)
	}

	if v := param.ActiveFlag; v {
		db = db.Where("active_flag=?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Employees, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.EmployeeQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a EmployeeRepository) Get(id string) (*models.Employee, error) {
	employee := new(models.Employee)

	if ok, err := QueryOne(a.db.ORM.Model(employee).Where("id=?", id), employee); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return employee, nil
}

func (a EmployeeRepository) Create(employee *models.Employee) error {
	result := a.db.ORM.Model(employee).Create(employee)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a EmployeeRepository) Update(id string, employee *models.Employee) error {
	result := a.db.ORM.Model(employee).Where("id=?", id).Updates(employee)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a EmployeeRepository) Delete(id string) error {
	employee := new(models.Employee)

	result := a.db.ORM.Model(employee).Where("id=?", id).Delete(employee)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a EmployeeRepository) UpdateStatus(id string, status int) error {
	employee := new(models.Employee)

	result := a.db.ORM.Model(employee).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
