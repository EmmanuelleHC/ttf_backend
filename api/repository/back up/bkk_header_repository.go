package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// BKKHeaderRepository database structure
type BKKHeaderRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewBKKHeaderRepository creates a new bkkheader repository
func NewBKKHeaderRepository(db lib.Database, logger lib.Logger) BKKHeaderRepository {
	return BKKHeaderRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a BKKHeaderRepository) WithTrx(trxHandle *gorm.DB) BKKHeaderRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a BKKHeaderRepository) Query(param *models.BKKHeaderQueryParam) (*models.BKKHeaderQueryResult, error) {
	db := a.db.ORM.Model(&models.BKKHeader{}).Preload("Company").Preload("Branch").Preload("BKKDetails")

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

	if v := param.StatusApprove; v != "" {
		db = db.Where("status_approve=?", v)
	}

	if v := param.CompanyID; v != "" {
		db = db.Where("company_id=?", v)
	}

	if v := param.BranchID; v != "" {
		db = db.Where("branch_id=?", v)
	}

	if v := param.DateQuery; len(v) != 0 {
		db = db.Where("created_at BETWEEN ? AND ?", param.DateQuery[0], param.DateQuery[1])
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("num LIKE ?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.BKKHeaders, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.BKKHeaderQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a BKKHeaderRepository) Get(id string) (*models.BKKHeader, error) {
	bkkheader := new(models.BKKHeader)

	if ok, err := QueryOne(a.db.ORM.Model(bkkheader).Where("id=?", id), bkkheader); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return bkkheader, nil
}

func (a BKKHeaderRepository) Create(bkkheader *models.BKKHeader) error {
	result := a.db.ORM.Model(bkkheader).
		Select("ID", "Num", "CompanyID", "BranchID", "ReleaseDate",
			"PaidDate", "TotalAmount", "KasbonID", "InvoiceID", "Status", "BKKDetails",
			"CreatedAt", "CreatedBy", "UpdateBy").Create(bkkheader)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKHeaderRepository) Update(id string, bkkheader *models.BKKHeader) error {
	result := a.db.ORM.Model(bkkheader).Where("id=?", id).Select("ID", "Num", "CompanyID", "BranchID", "ReleaseDate",
		"PaidDate", "TotalAmount", "KasbonID", "InvoiceID", "Status", "BKKDetails",
		"CreatedAt", "CreatedBy", "UpdatedAt", "UpdateBy").Updates(bkkheader)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKHeaderRepository) Delete(id string) error {
	bkkheader := new(models.BKKHeader)

	result := a.db.ORM.Model(bkkheader).Where("id=?", id).Delete(bkkheader)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKHeaderRepository) UpdateStatus(id string, status string) error {
	bkkheader := new(models.BKKHeader)

	result := a.db.ORM.Model(bkkheader).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKHeaderRepository) UpdateApprove(ids []string, status int) error {
	bkkheader := new(models.BKKHeader)

	result := a.db.ORM.Model(bkkheader).Where("id IN (?)", ids).Update("status_approve", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
