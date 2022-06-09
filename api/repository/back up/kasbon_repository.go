package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// KasbonRepository database structure
type KasbonRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewKasbonRepository creates a new kasbon repository
func NewKasbonRepository(db lib.Database, logger lib.Logger) KasbonRepository {
	return KasbonRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a KasbonRepository) WithTrx(trxHandle *gorm.DB) KasbonRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a KasbonRepository) Query(param *models.KasbonQueryParam) (*models.KasbonQueryResult, error) {
	db := a.db.ORM.Model(&models.Kasbon{}).Preload("Company").Preload("Branch").Preload("Department").Preload("Employee")

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.Description; v != "" {
		v = "%" + v + "%"
		db = db.Where("description LIKE ?", v)
	}

	if v := param.Status; v != "" {
		db = db.Where("status=?", v)
	}

	if v := param.EmployeeID; v != "" {
		db = db.Where("employee_id=?", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.BranchID; v != "" {
		db = db.Where("branch_id=?", v)
	}

	if v := param.DeptID; v != "" {
		db = db.Where("department_id=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("num LIKE ? OR description LIKE ?", v, v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Kasbons, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.KasbonQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a KasbonRepository) Get(id string) (*models.Kasbon, error) {
	kasbon := new(models.Kasbon)

	if ok, err := QueryOne(a.db.ORM.Model(kasbon).Where("id=?", id), kasbon); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return kasbon, nil
}

func (a KasbonRepository) Create(kasbon *models.Kasbon) error {
	result := a.db.ORM.Model(kasbon).Create(kasbon)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a KasbonRepository) Update(id string, kasbon *models.Kasbon) error {
	result := a.db.ORM.Model(kasbon).Where("id=?", id).Updates(kasbon)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a KasbonRepository) Delete(id string) error {
	kasbon := new(models.Kasbon)

	result := a.db.ORM.Model(kasbon).Where("id=?", id).Delete(kasbon)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a KasbonRepository) UpdateStatus(id string, status int) error {
	kasbon := new(models.Kasbon)

	result := a.db.ORM.Model(kasbon).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
