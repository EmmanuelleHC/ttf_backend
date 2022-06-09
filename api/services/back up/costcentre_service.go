package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// CostCentreService service layer
type CostCentreService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	costcentreRepository repository.CostCentreRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewCostCentreService creates a new costcentreservice
func NewCostCentreService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	costcentreRepository repository.CostCentreRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) CostCentreService {
	return CostCentreService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		costcentreRepository: costcentreRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a CostCentreService) WithTrx(trxHandle *gorm.DB) CostCentreService {
	a.costcentreRepository = a.costcentreRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a CostCentreService) Query(param *models.CostCentreQueryParam) (costcentreQR *models.CostCentreQueryResult, err error) {
	return a.costcentreRepository.Query(param)
}

func (a CostCentreService) Get(id string) (*models.CostCentre, error) {
	costcentre, err := a.costcentreRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return costcentre, nil
}

func (a CostCentreService) Check(item *models.CostCentre) error {
	qr, err := a.costcentreRepository.Query(&models.CostCentreQueryParam{Name: item.Name})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.CostCentreAlreadyExists
	}

	return nil
}

func (a CostCentreService) Create(costcentre *models.CostCentre) (id string, err error) {
	if err = a.Check(costcentre); err != nil {
		return
	}

	costcentre.ID = uuid.MustString()

	if err = a.costcentreRepository.Create(costcentre); err != nil {
		return
	}
	return costcentre.ID, nil
}

func (a CostCentreService) Update(id string, costcentre *models.CostCentre) error {
	oCostCentre, err := a.Get(id)
	if err != nil {
		return err
	} else if costcentre.Name != oCostCentre.Name {
		if err = a.Check(costcentre); err != nil {
			return err
		}
	}
	costcentre.ID = oCostCentre.ID

	if err := a.costcentreRepository.Update(id, costcentre); err != nil {
		return err
	}

	return nil
}

func (a CostCentreService) Delete(id string) error {
	_, err := a.costcentreRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.costcentreRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a CostCentreService) UpdateStatus(id string, status int) error {
	_, err := a.costcentreRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.costcentreRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
