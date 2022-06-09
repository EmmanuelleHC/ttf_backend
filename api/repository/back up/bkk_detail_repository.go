package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// BKKDetailRepository database structure
type BKKDetailRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewBKKDetailRepository creates a new bkkdetail repository
func NewBKKDetailRepository(db lib.Database, logger lib.Logger) BKKDetailRepository {
	return BKKDetailRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a BKKDetailRepository) WithTrx(trxHandle *gorm.DB) BKKDetailRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a BKKDetailRepository) Query(param *models.BKKDetailQueryParam) (*models.BKKDetailQueryResult, error) {
	db := a.db.ORM.Model(&models.BKKDetail{}).Preload("Trx")

	if v := param.IDs; len(v) > 0 {
		db = db.Where("id IN (?)", v)
	}

	if v := param.BKKHeaderIDs; len(v) > 0 {
		db = db.Where("bkk_header_id IN (?)", v)
	}

	if v := param.LinesDesc; v != "" {
		v = "%" + v + "%"
		db = db.Where("lines_desc LIKE ?", v)
	}

	if v := param.Status; v != "" {
		db = db.Where("status=?", v)
	}

	if v := param.BKKHeaderID; v != "" {
		db = db.Where("bkk_header_id=?", v)
	}

	if v := param.QueryValue; v != "" {
		v = "%" + v + "%"
		db = db.Where("lines_desc LIKE ?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.BKKDetails, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.BKKDetailQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a BKKDetailRepository) Get(id string) (*models.BKKDetail, error) {
	bkkdetail := new(models.BKKDetail)

	if ok, err := QueryOne(a.db.ORM.Model(bkkdetail).Where("id=?", id), bkkdetail); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return bkkdetail, nil
}

func (a BKKDetailRepository) Create(bkkdetail *models.BKKDetail) error {
	result := a.db.ORM.Model(bkkdetail).Create(bkkdetail)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKDetailRepository) CreateBatch(bkkdetails models.BKKDetails) error {
	result := a.db.ORM.Model(bkkdetails).Create(bkkdetails)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKDetailRepository) Update(id string, bkkdetail *models.BKKDetail) error {
	result := a.db.ORM.Model(bkkdetail).Where("id=?", id).Updates(bkkdetail)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKDetailRepository) Delete(id string) error {
	bkkdetail := new(models.BKKDetail)

	result := a.db.ORM.Model(bkkdetail).Where("id=?", id).Delete(bkkdetail)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKDetailRepository) DeleteByHeaderID(id string) error {
	bkkdetail := new(models.BKKDetail)

	result := a.db.ORM.Model(bkkdetail).Unscoped().Where("bkk_header_id=?", id).Delete(bkkdetail)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a BKKDetailRepository) UpdateStatus(id string, status int) error {
	bkkdetail := new(models.BKKDetail)

	result := a.db.ORM.Model(bkkdetail).Where("id=?", id).Update("status", status)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
