package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// TrxService service layer
type TrxService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	trxRepository        repository.TrxRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewTrxService creates a new trxservice
func NewTrxService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	trxRepository repository.TrxRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) TrxService {
	return TrxService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		trxRepository:        trxRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a TrxService) WithTrx(trxHandle *gorm.DB) TrxService {
	a.trxRepository = a.trxRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a TrxService) Query(param *models.TrxQueryParam) (trxQR *models.TrxQueryResult, err error) {
	return a.trxRepository.Query(param)
}

func (a TrxService) Get(id string) (*models.Trx, error) {
	trx, err := a.trxRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return trx, nil
}

func (a TrxService) Check(item *models.Trx) error {
	qr, err := a.trxRepository.Query(&models.TrxQueryParam{Name: item.Name})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.TrxAlreadyExists
	}

	return nil
}

func (a TrxService) Create(trx *models.Trx) (id string, err error) {
	if err = a.Check(trx); err != nil {
		return
	}

	trx.ID = uuid.MustString()

	if err = a.trxRepository.Create(trx); err != nil {
		return
	}
	return trx.ID, nil
}

func (a TrxService) Update(id string, trx *models.Trx) error {
	oTrx, err := a.Get(id)
	if err != nil {
		return err
	} else if trx.Name != oTrx.Name {
		if err = a.Check(trx); err != nil {
			return err
		}
	}
	trx.ID = oTrx.ID

	if err := a.trxRepository.Update(id, trx); err != nil {
		return err
	}

	return nil
}

func (a TrxService) Delete(id string) error {
	_, err := a.trxRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.trxRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a TrxService) UpdateStatus(id string, status int) error {
	_, err := a.trxRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.trxRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
