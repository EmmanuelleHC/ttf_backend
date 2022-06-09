package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// TarikDanaService service layer
type TarikDanaService struct {
	logger                 lib.Logger
	casbinService          CasbinService
	counterService         CounterService
	saldoService           SaldoService
	branchRepository       repository.BranchRepository
	tarikdanaRepository    repository.TarikDanaRepository
	counterRepository      repository.CounterRepository
	saldohistoryRepository repository.SaldoHistoryRepository
	saldoRepository        repository.SaldoRepository
	saldomonthRepository   repository.SaldoMonthRepository
}

// NewTarikDanaService creates a new tarikdanaservice
func NewTarikDanaService(
	logger lib.Logger,
	casbinService CasbinService,
	counterService CounterService,
	saldoService SaldoService,
	branchRepository repository.BranchRepository,
	tarikdanaRepository repository.TarikDanaRepository,
	counterRepository repository.CounterRepository,
	saldohistoryRepository repository.SaldoHistoryRepository,
	saldoRepository repository.SaldoRepository,
	saldomonthRepository repository.SaldoMonthRepository,
) TarikDanaService {
	return TarikDanaService{
		logger:                 logger,
		casbinService:          casbinService,
		counterService:         counterService,
		saldoService:           saldoService,
		branchRepository:       branchRepository,
		tarikdanaRepository:    tarikdanaRepository,
		counterRepository:      counterRepository,
		saldohistoryRepository: saldohistoryRepository,
		saldoRepository:        saldoRepository,
		saldomonthRepository:   saldomonthRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a TarikDanaService) WithTrx(trxHandle *gorm.DB) TarikDanaService {
	a.tarikdanaRepository = a.tarikdanaRepository.WithTrx(trxHandle)
	a.branchRepository = a.branchRepository.WithTrx(trxHandle)
	a.counterRepository = a.counterRepository.WithTrx(trxHandle)
	a.saldohistoryRepository = a.saldohistoryRepository.WithTrx(trxHandle)
	a.saldomonthRepository = a.saldomonthRepository.WithTrx(trxHandle)
	a.saldoRepository = a.saldoRepository.WithTrx(trxHandle)

	return a
}

func (a TarikDanaService) Query(param *models.TarikDanaQueryParam) (tarikdanaQR *models.TarikDanaQueryResult, err error) {
	return a.tarikdanaRepository.Query(param)
}

func (a TarikDanaService) Get(id string) (*models.TarikDana, error) {
	tarikdana, err := a.tarikdanaRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return tarikdana, nil
}

func (a TarikDanaService) Check(item *models.TarikDana) error {
	qr, err := a.tarikdanaRepository.Query(&models.TarikDanaQueryParam{
		Description: item.Description,
		CompanyID:   item.CompanyID,
		BranchID:    item.BranchID,
	})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.TarikDanaAlreadyExists
	}

	return nil
}

func (a TarikDanaService) Create(tarikdana *models.TarikDana) (id string, err error) {
	if err = a.Check(tarikdana); err != nil {
		return
	}

	saldoNow, err := a.saldoService.CreateNewSaldoOrUpdate(tarikdana.CompanyID, tarikdana.BranchID, 0, tarikdana.Amount)
	if err != nil {
		return
	}
	// saldoHis, _ := a.saldohistoryRepository.GetbyCompanyAndBranch(tarikdana.CompanyID, tarikdana.BranchID)

	saldoHisCreate := new(models.SaldoHistory)
	saldoHisCreate.Desc = "Penerimaan Dana"
	saldoHisCreate.CompanyID = tarikdana.CompanyID
	saldoHisCreate.BranchID = tarikdana.BranchID
	saldoHisCreate.InAmount = tarikdana.Amount
	saldoHisCreate.SaldoAkhir = saldoNow
	if err = a.saldohistoryRepository.Create(saldoHisCreate); err != nil {
		return "", err
	}

	tarikdana.ID = uuid.MustString()

	if err = a.tarikdanaRepository.Create(tarikdana); err != nil {
		return
	}

	return tarikdana.ID, nil
}

func (a TarikDanaService) Update(id string, tarikdana *models.TarikDana) error {
	oTarikDana, err := a.Get(id)
	if err != nil {
		return err
	} else if tarikdana.Description != oTarikDana.Description {
		if err = a.Check(tarikdana); err != nil {
			return err
		}
	}
	tarikdana.ID = oTarikDana.ID

	if err := a.tarikdanaRepository.Update(id, tarikdana); err != nil {
		return err
	}

	return nil
}

func (a TarikDanaService) Delete(id string) error {
	_, err := a.tarikdanaRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.tarikdanaRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a TarikDanaService) UpdateStatus(id string, status int) error {
	_, err := a.tarikdanaRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.tarikdanaRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
