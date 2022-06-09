package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// SaldoHistoryRepository database structure
type SaldoHistoryRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewSaldoHistoryRepository creates a new saldohistory repository
func NewSaldoHistoryRepository(db lib.Database, logger lib.Logger) SaldoHistoryRepository {
	return SaldoHistoryRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a SaldoHistoryRepository) WithTrx(trxHandle *gorm.DB) SaldoHistoryRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a SaldoHistoryRepository) Query(param *models.SaldoHistoryQueryParam) (*models.SaldoHistoryQueryResult, error) {
	db := a.db.ORM.Model(&models.SaldoHistory{})

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
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

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.SaldoHistories, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.SaldoHistoryQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a SaldoHistoryRepository) Get(id string) (*models.SaldoHistory, error) {
	saldohistory := new(models.SaldoHistory)

	if ok, err := QueryOne(a.db.ORM.Model(saldohistory).Where("id=?", id), saldohistory); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldohistory, nil
}

func (a SaldoHistoryRepository) GetbyCompanyAndBranch(companyId string, branchId string) (*models.SaldoHistory, error) {
	saldohistory := new(models.SaldoHistory)

	if ok, err := QueryOne(a.db.ORM.Model(saldohistory).Order("record_id desc").Where("company_id=? AND branch_id=?", companyId, branchId), saldohistory); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return saldohistory, nil
}

func (a SaldoHistoryRepository) GetSaldoForAwal(companyId string, branchId string, dateFrom string, dateTo string) (int64, int64, error) {
	saldohistory := new(models.SaldoHistory)
	type resultData struct {
		TotalOut, TotalIn int64
	}
	var total resultData
	result := a.db.ORM.Model(saldohistory).Select("sum(out_amount) as TotalOut, sum(in_amount) as TotalIn").
		Where("company_id=? AND branch_id=? AND created_at BETWEEN ? AND ?", companyId, branchId, dateFrom, dateTo).First(&total)
	if result.Error != nil {
		return 0, 0, errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return total.TotalIn, total.TotalOut, nil
}

func (a SaldoHistoryRepository) GetTotalOut(companyId string, branchId string, desc string) (int64, error) {
	saldohistory := new(models.SaldoHistory)
	var total int64 = 0
	result := a.db.ORM.Model(saldohistory).Select("sum(out_amount) as total").
		Where("company_id=? AND branch_id=? AND desc=?", companyId, branchId, desc).Group("desc").First(&total)
	if result.Error != nil {
		return 0, errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return total, nil
}

func (a SaldoHistoryRepository) Create(saldohistory *models.SaldoHistory) error {
	result := a.db.ORM.Model(saldohistory).Create(saldohistory)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoHistoryRepository) Update(id string, saldohistory *models.SaldoHistory) error {
	result := a.db.ORM.Model(saldohistory).Where("id=?", id).Updates(saldohistory)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoHistoryRepository) Delete(id string) error {
	saldohistory := new(models.SaldoHistory)

	result := a.db.ORM.Model(saldohistory).Where("id=?", id).Delete(saldohistory)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoHistoryRepository) DeleteByDesc(desc string) error {
	saldohistory := new(models.SaldoHistory)

	result := a.db.ORM.Model(saldohistory).Where("desc=?", desc).Delete(saldohistory)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a SaldoHistoryRepository) UpdateStatus(id string, status int) error {
	saldohistory := new(models.SaldoHistory)

	result := a.db.ORM.Model(saldohistory).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
