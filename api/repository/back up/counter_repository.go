package repository

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// CounterRepository database structure
type CounterRepository struct {
	db     lib.Database
	logger lib.Logger
}

// NewCounterRepository creates a new counter repository
func NewCounterRepository(db lib.Database, logger lib.Logger) CounterRepository {
	return CounterRepository{
		db:     db,
		logger: logger,
	}
}

// WithTrx enables repository with transaction
func (a CounterRepository) WithTrx(trxHandle *gorm.DB) CounterRepository {
	if trxHandle == nil {
		a.logger.Zap.Error("Transaction Database not found in echo context. ")
		return a
	}

	a.db.ORM = trxHandle
	return a
}

func (a CounterRepository) Query(param *models.CounterQueryParam) (*models.CounterQueryResult, error) {
	db := a.db.ORM.Model(&models.Counter{})

	if v := param.Key; v != "" {
		db = db.Where("key=?", v)
	}

	if v := param.QueryValue; v != "" {
		db = db.Where("key=?", v)
	}

	db = db.Order(param.OrderParam.ParseOrder())

	list := make(models.Counters, 0)
	pagination, err := QueryPagination(db, param.PaginationParam, &list)
	if err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	}

	qr := &models.CounterQueryResult{
		Pagination: pagination,
		List:       list,
	}

	return qr, nil
}

func (a CounterRepository) Get(id string) (*models.Counter, error) {
	counter := new(models.Counter)

	if ok, err := QueryOne(a.db.ORM.Model(counter).Where("id=?", id), counter); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return counter, nil
}

func (a CounterRepository) GetbyKey(key string) (*models.Counter, error) {
	counter := new(models.Counter)

	if ok, err := QueryOne(a.db.ORM.Model(counter).Where("key_counter = ?", key), counter); err != nil {
		return nil, errors.Wrap(errors.DatabaseInternalError, err.Error())
	} else if !ok {
		return nil, errors.DatabaseRecordNotFound
	}

	return counter, nil
}

func (a CounterRepository) Create(counter *models.Counter) error {
	result := a.db.ORM.Model(counter).Create(counter)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CounterRepository) Update(id string, counter *models.Counter) error {
	result := a.db.ORM.Model(counter).Where("id=?", id).Updates(counter)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}

func (a CounterRepository) Delete(id string) error {
	counter := new(models.Counter)

	result := a.db.ORM.Model(counter).Where("id=?", id).Delete(counter)
	if result.Error != nil {
		return errors.Wrap(errors.DatabaseInternalError, result.Error.Error())
	}

	return nil
}
