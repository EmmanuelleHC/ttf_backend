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

// BKKHeaderService service layer
type BKKHeaderService struct {
	logger                 lib.Logger
	casbinService          CasbinService
	counterService         CounterService
	saldoService           SaldoService
	userRepository         repository.UserRepository
	branchRepository       repository.BranchRepository
	bkkheaderRepository    repository.BKKHeaderRepository
	bkkdetailRepository    repository.BKKDetailRepository
	counterRepository      repository.CounterRepository
	saldoRepository        repository.SaldoRepository
	saldohistoryRepository repository.SaldoHistoryRepository
}

// NewBKKHeaderService creates a new bkkheaderservice
func NewBKKHeaderService(
	logger lib.Logger,
	casbinService CasbinService,
	counterService CounterService,
	saldoService SaldoService,
	userRepository repository.UserRepository,
	branchRepository repository.BranchRepository,
	bkkheaderRepository repository.BKKHeaderRepository,
	bkkdetailRepository repository.BKKDetailRepository,
	counterRepository repository.CounterRepository,
	saldoRepository repository.SaldoRepository,
	saldohistoryRepository repository.SaldoHistoryRepository,
) BKKHeaderService {
	return BKKHeaderService{
		logger:                 logger,
		casbinService:          casbinService,
		counterService:         counterService,
		saldoService:           saldoService,
		userRepository:         userRepository,
		branchRepository:       branchRepository,
		bkkheaderRepository:    bkkheaderRepository,
		bkkdetailRepository:    bkkdetailRepository,
		counterRepository:      counterRepository,
		saldoRepository:        saldoRepository,
		saldohistoryRepository: saldohistoryRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a BKKHeaderService) WithTrx(trxHandle *gorm.DB) BKKHeaderService {
	a.bkkheaderRepository = a.bkkheaderRepository.WithTrx(trxHandle)
	a.bkkdetailRepository = a.bkkdetailRepository.WithTrx(trxHandle)
	a.counterRepository = a.counterRepository.WithTrx(trxHandle)
	a.saldoRepository = a.saldoRepository.WithTrx(trxHandle)
	a.saldohistoryRepository = a.saldohistoryRepository.WithTrx(trxHandle)

	return a
}

func (a BKKHeaderService) Query(param *models.BKKHeaderQueryParam) (bkkheaderQR *models.BKKHeaderQueryResult, err error) {
	return a.bkkheaderRepository.Query(param)
}

func (a BKKHeaderService) QueryByUser(param *models.BKKHeaderQueryParam) (bkkheaderQR *models.BKKHeaderQueryResult, err error) {
	usr, err := a.userRepository.Get(param.UserId)
	param.CompanyID = usr.CompanyID
	param.BranchID = usr.BranchID
	if err != nil {
		return nil, err
	}
	return a.bkkheaderRepository.Query(param)
}

func (a BKKHeaderService) Get(id string) (*models.BKKHeader, error) {
	bkkheader, err := a.bkkheaderRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return bkkheader, nil
}

func (a BKKHeaderService) Check(item *models.BKKHeader) error {
	qr, err := a.bkkheaderRepository.Query(&models.BKKHeaderQueryParam{Num: item.Num})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.BKKHeaderAlreadyExists
	}

	return nil
}

func (a BKKHeaderService) Create(bkkheader *models.BKKHeader) (id string, err error) {
	cbn, err := a.branchRepository.Get(bkkheader.BranchID)
	if err != nil {
		return "", err
	}
	key := "BKK" + cbn.Shorter
	count, err := a.counterService.GetAndIncrement(key)
	padleft := fmt.Sprintf("%04d", count)
	bkkNum := key + padleft
	if err != nil {
		return "", err
	}
	var total int64 = 0
	for _, item := range bkkheader.BKKDetails {
		total += item.LinesAmount
	}

	saldo, err := a.saldoRepository.GetbyCompanyAndBranch(bkkheader.CompanyID, bkkheader.BranchID)
	if err != nil {
		return "", err
	}

	if total > saldo.LimitBKK {
		errtotal := fmt.Errorf("not allowed more than limit: %d", saldo.LimitBKK)
		return "", errtotal
	}

	if total > saldo.SaldoAkhir {
		errtotal := fmt.Errorf("not enough saldo: %d", saldo.SaldoAkhir)
		return "", errtotal
	}

	bkkheader.ID = uuid.MustString()
	bkkheader.Num = bkkNum
	bkkheader.TotalAmount = total

	if err = a.bkkheaderRepository.Create(bkkheader); err != nil {
		return "", err
	}

	saldoHisCreate := new(models.SaldoHistory)
	saldoHisCreate.Desc = bkkNum
	saldoHisCreate.CompanyID = bkkheader.CompanyID
	saldoHisCreate.BranchID = bkkheader.BranchID
	saldoHisCreate.OutAmount = total
	saldoHisCreate.SaldoAkhir = 0
	if err = a.saldohistoryRepository.Create(saldoHisCreate); err != nil {
		return "", err
	}

	return bkkheader.ID, nil
}

func (a BKKHeaderService) Update(id string, bkkheader *models.BKKHeader) error {
	oBKKHeader, err := a.Get(id)
	if err != nil {
		return err
	} else if bkkheader.Num != oBKKHeader.Num {
		if err = a.Check(bkkheader); err != nil {
			return err
		}
	}
	bkkheader.ID = oBKKHeader.ID

	if err := a.bkkdetailRepository.DeleteByHeaderID(id); err != nil {
		return err
	}

	if err := a.bkkheaderRepository.Update(id, bkkheader); err != nil {
		return err
	}

	return nil
}

func (a BKKHeaderService) Delete(id string) error {
	_, err := a.bkkheaderRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.bkkheaderRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a BKKHeaderService) UpdateStatus(id string, status string, statusReject int) error {
	bkk, err := a.bkkheaderRepository.Get(id)
	if err != nil {
		return err
	}

	if status == "Paid" && statusReject == 0 {
		if _, err = a.saldoService.CreateNewSaldoOrUpdate(bkk.CompanyID, bkk.BranchID, bkk.TotalAmount, 0); err != nil {
			return err
		}
	}

	if err := a.bkkheaderRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}

func (a BKKHeaderService) UpdateApprove(ids []string, status int) error {
	reject := -1
	if status == reject {
		for _, id := range ids {
			bkk, err := a.bkkheaderRepository.Get(id)
			if err != nil {
				return err
			}

			qr, err := a.saldohistoryRepository.Query(&models.SaldoHistoryQueryParam{
				Desc:      bkk.Num,
				CompanyID: bkk.CompanyID,
				BranchID:  bkk.BranchID,
			})

			if err != nil {
				return err
			} else if len(qr.List) > 0 {
				if err := a.saldohistoryRepository.DeleteByDesc(bkk.Num); err != nil {
					return err
				}
			}
		}
	}

	if err := a.bkkheaderRepository.UpdateApprove(ids, status); err != nil {
		return err
	}

	return nil
}
