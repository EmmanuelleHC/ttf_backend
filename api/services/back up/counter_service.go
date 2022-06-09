package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// CounterService service layer
type CounterService struct {
	logger            lib.Logger
	casbinService     CasbinService
	counterRepository repository.CounterRepository
}

// NewCounterService creates a new counterservice
func NewCounterService(
	logger lib.Logger,
	casbinService CasbinService,
	counterRepository repository.CounterRepository,
) CounterService {
	return CounterService{
		logger:            logger,
		casbinService:     casbinService,
		counterRepository: counterRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a CounterService) WithTrx(trxHandle *gorm.DB) CounterService {
	a.counterRepository = a.counterRepository.WithTrx(trxHandle)

	return a
}

func (a CounterService) Query(param *models.CounterQueryParam) (counterQR *models.CounterQueryResult, err error) {
	return a.counterRepository.Query(param)
}

func (a CounterService) Get(id string) (*models.Counter, error) {
	counter, err := a.counterRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

func (a CounterService) GetAndIncrement(key string) (count int, err error) {
	counterdb, err := a.counterRepository.GetbyKey(key)
	if err != nil {
		counter := new(models.Counter)
		counter.ID = uuid.MustString()
		counter.KeyCounter = key
		counter.CounterMe = 2
		if err = a.counterRepository.Create(counter); err != nil {
			return 0, err
		}
		return 1, nil
	}

	if err != nil {
		return 0, err
	}

	counterdb.CounterMe = counterdb.CounterMe + 1

	if err = a.counterRepository.Update(counterdb.ID, counterdb); err != nil {
		return 0, err
	}

	return counterdb.CounterMe + 1, nil
}
