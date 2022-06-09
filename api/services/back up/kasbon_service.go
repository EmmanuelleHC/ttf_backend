package services

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// KasbonService service layer
type KasbonService struct {
	logger             lib.Logger
	casbinService      CasbinService
	counterService     CounterService
	employeeRepository repository.EmployeeRepository
	branchRepository   repository.BranchRepository
	kasbonRepository   repository.KasbonRepository
	counterRepository  repository.CounterRepository
}

// NewKasbonService creates a new kasbonservice
func NewKasbonService(
	logger lib.Logger,
	casbinService CasbinService,
	counterService CounterService,
	employeeRepository repository.EmployeeRepository,
	branchRepository repository.BranchRepository,
	kasbonRepository repository.KasbonRepository,
	counterRepository repository.CounterRepository,
) KasbonService {
	return KasbonService{
		logger:             logger,
		casbinService:      casbinService,
		counterService:     counterService,
		employeeRepository: employeeRepository,
		branchRepository:   branchRepository,
		kasbonRepository:   kasbonRepository,
		counterRepository:  counterRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a KasbonService) WithTrx(trxHandle *gorm.DB) KasbonService {
	a.kasbonRepository = a.kasbonRepository.WithTrx(trxHandle)
	a.employeeRepository = a.employeeRepository.WithTrx(trxHandle)
	a.branchRepository = a.branchRepository.WithTrx(trxHandle)
	a.counterRepository = a.counterRepository.WithTrx(trxHandle)

	return a
}

func (a KasbonService) Query(param *models.KasbonQueryParam) (kasbonQR *models.KasbonQueryResult, err error) {
	return a.kasbonRepository.Query(param)
}

func (a KasbonService) Get(id string) (*models.Kasbon, error) {
	kasbon, err := a.kasbonRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return kasbon, nil
}

func (a KasbonService) Check(item *models.Kasbon) error {
	qr, err := a.kasbonRepository.Query(&models.KasbonQueryParam{Num: item.Num})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.KasbonAlreadyExists
	}

	return nil
}

func (a KasbonService) Create(kasbon *models.Kasbon) (id string, err error) {
	usr, err := a.employeeRepository.Get(kasbon.EmployeeID)
	if err != nil {
		return "", err
	}
	cbn, err := a.branchRepository.Get(usr.BranchID)
	if err != nil {
		return "", err
	}
	key := "KBS" + cbn.Shorter
	count, err := a.counterService.GetAndIncrement(key)
	if err != nil {
		return "", err
	}
	padleft := fmt.Sprintf("%04d", count)
	kasbon.ID = uuid.MustString()
	kasbon.CompanyID = usr.CompanyID
	kasbon.BranchID = usr.BranchID
	kasbon.Num = key + padleft

	if err = a.kasbonRepository.Create(kasbon); err != nil {
		return
	}
	return kasbon.ID, nil
}

func (a KasbonService) Update(id string, kasbon *models.Kasbon) error {
	oKasbon, err := a.Get(id)
	if err != nil {
		return err
	} else if kasbon.Description != oKasbon.Description {
		if err = a.Check(kasbon); err != nil {
			return err
		}
	}
	kasbon.ID = oKasbon.ID

	if err := a.kasbonRepository.Update(id, kasbon); err != nil {
		return err
	}

	return nil
}

func (a KasbonService) Delete(id string) error {
	_, err := a.kasbonRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.kasbonRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a KasbonService) UpdateStatus(id string, status int) error {
	_, err := a.kasbonRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.kasbonRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
