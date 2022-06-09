package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// SaldoService service layer
type SaldoService struct {
	logger                 lib.Logger
	casbinService          CasbinService
	userRepository         repository.UserRepository
	saldoRepository        repository.SaldoRepository
	saldomonthRepository   repository.SaldoMonthRepository
	saldohistoryRepository repository.SaldoHistoryRepository
	menuRepository         repository.MenuRepository
	menuActionRepository   repository.MenuActionRepository
}

// NewSaldoService creates a new saldoservice
func NewSaldoService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	saldoRepository repository.SaldoRepository,
	saldomonthRepository repository.SaldoMonthRepository,
	saldohistoryRepository repository.SaldoHistoryRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) SaldoService {
	return SaldoService{
		logger:                 logger,
		casbinService:          casbinService,
		userRepository:         userRepository,
		saldoRepository:        saldoRepository,
		saldomonthRepository:   saldomonthRepository,
		saldohistoryRepository: saldohistoryRepository,
		menuRepository:         menuRepository,
		menuActionRepository:   menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a SaldoService) WithTrx(trxHandle *gorm.DB) SaldoService {
	a.saldoRepository = a.saldoRepository.WithTrx(trxHandle)
	a.saldomonthRepository = a.saldomonthRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a SaldoService) Query(param *models.SaldoQueryParam) (saldoQR *models.SaldoQueryResult, err error) {
	return a.saldoRepository.Query(param)
}

func (a SaldoService) Get(id string) (*models.Saldo, error) {
	saldo, err := a.saldoRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return saldo, nil
}

func (a SaldoService) GetByUser(username string) (*models.Saldo, error) {
	if username == "root" {
		saldo, err := a.saldoRepository.GetRoot()
		if err != nil {
			return nil, err
		}
		return saldo, nil
	}
	user, err := a.userRepository.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	saldo, err := a.saldoRepository.GetbyCompanyAndBranch(user.CompanyID, user.BranchID)
	if err != nil {
		return nil, err
	}
	return saldo, nil
}

func (a SaldoService) GetMonthbyID(id string) (*models.SaldoMonth, error) {
	saldo, err := a.saldomonthRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return saldo, nil
}

func (a SaldoService) GetMonth(companyId string, branchId string, monthYear string) (*models.SaldoMonth, error) {
	saldo, err := a.saldomonthRepository.GetbyCompanyAndBranchMonth(companyId, branchId, monthYear)
	if err != nil {
		return nil, err
	}
	return saldo, nil
}

func (a SaldoService) GetSaldoAwal(companyId string, branchId string, dateTo string, firstday string) (int64, error) {
	t := time.Now()
	var saldoAwal int64 = 0

	saldo, err := a.GetMonth(companyId, branchId, t.Format("2006-01"))
	if err != nil {
		return 0, err
	}

	in, out, err := a.saldohistoryRepository.GetSaldoForAwal(companyId, branchId, firstday, dateTo)
	if err != nil {
		return 0, err
	}

	saldoAwal = saldo.SaldoAwal + in - out

	return saldoAwal, nil
}

func (a SaldoService) Check(item *models.Saldo) error {
	qr, err := a.saldoRepository.Query(&models.SaldoQueryParam{
		CompanyID: item.CompanyID,
		BranchID:  item.BranchID,
	})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.SaldoAlreadyExists
	}

	return nil
}

func (a SaldoService) CheckIsNew(item *models.SaldoMonth) error {
	qr, err := a.saldomonthRepository.Query(&models.SaldoMonthQueryParam{
		CompanyID:  item.CompanyID,
		BranchID:   item.BranchID,
		OrderParam: dto.OrderParam{Key: "created_at", Direction: dto.OrderByDESC},
	})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		if qr.List[0].MonthYear != item.MonthYear {
			return errors.SaldoNeedsNew
		}
	}

	return nil
}

func (a SaldoService) CheckMonth(item *models.SaldoMonth) error {
	qr, err := a.saldomonthRepository.Query(&models.SaldoMonthQueryParam{
		CompanyID: item.CompanyID,
		BranchID:  item.BranchID,
		MonthYear: item.MonthYear,
	})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.SaldoAlreadyExists
	}

	return nil
}

func (a SaldoService) Create(saldo *models.Saldo) (id string, err error) {
	if err = a.Check(saldo); err != nil {
		return
	}
	begin := time.Now()

	saldo.ID = uuid.MustString()
	saldo.MonthYear = begin.Format("2006-01")

	if err = a.saldoRepository.Create(saldo); err != nil {
		return
	}
	return saldo.ID, nil
}

func (a SaldoService) CreateAndMonth(saldo *models.Saldo) (id string, err error) {
	if err = a.Check(saldo); err != nil {
		return
	}
	begin := time.Now()

	saldo.ID = uuid.MustString()
	saldo.MonthYear = begin.Format("2006-01")
	saldo.SaldoAkhir = saldo.SaldoAwal

	if err = a.saldoRepository.Create(saldo); err != nil {
		return
	}

	saldoMonth := new(models.SaldoMonth)
	saldoMonth.ID = uuid.MustString()
	saldoMonth.CompanyID = saldo.CompanyID
	saldoMonth.BranchID = saldo.BranchID
	saldoMonth.MonthYear = begin.Format("2006-01")
	saldoMonth.Month = int(begin.Month())
	saldoMonth.Year = begin.Year()
	saldoMonth.SaldoAwal = saldo.SaldoAwal
	saldoMonth.SaldoAkhir = saldo.SaldoAwal
	saldoMonth.CreatedBy = saldo.CreatedBy

	if err = a.saldomonthRepository.Create(saldoMonth); err != nil {
		return
	}

	return saldo.ID, nil
}

func (a SaldoService) CreateNewSaldoOrUpdate(companyID string, branchID string, totalOut int64, totalIn int64) (saldoNow int64, err error) {
	now := time.Now()
	prev := now.AddDate(0, -1, 0)
	prevMonthYear := prev.Format("2006-01")
	nowMonthYear := now.Format("2006-01")
	saldo, err := a.saldoRepository.GetbyCompanyAndBranch(companyID, branchID)
	if err != nil {
		return 0, err
	}

	var saldoAkhir int64 = saldo.SaldoAkhir - totalOut + totalIn
	if saldoAkhir < 0 {
		var err error = fmt.Errorf("the balance is not sufficient")
		return 0, err
	}

	saldo.UsedBKK += totalOut
	saldo.SaldoIn += totalIn
	saldo.SaldoAkhir = saldoAkhir

	if err = a.saldoRepository.Update(saldo.ID, saldo); err != nil {
		return 0, err
	}

	saldomonthnow, err := a.saldomonthRepository.GetbyCompanyAndBranchMonth(companyID, branchID, nowMonthYear)
	if err == errors.DatabaseRecordNotFound {
		saldoprev, errs := a.saldomonthRepository.GetbyCompanyAndBranchMonth(companyID, branchID, prevMonthYear)
		if errs != nil {
			return 0, errs
		}
		saldomonth := new(models.SaldoMonth)
		saldomonth.ID = uuid.MustString()
		saldomonth.CompanyID = companyID
		saldomonth.BranchID = branchID
		saldomonth.MonthYear = nowMonthYear
		saldomonth.Month = int(now.Month())
		saldomonth.Year = now.Year()
		saldomonth.SaldoAwal = saldoprev.SaldoAkhir
		saldomonth.SaldoIn = totalIn
		saldomonth.UsedBKK = totalOut
		saldomonth.SaldoAkhir = saldoprev.SaldoAkhir - totalOut + totalIn
		saldomonth.CreatedBy = saldoprev.CreatedBy

		if err = a.saldomonthRepository.Create(saldomonth); err != nil {
			return 0, err
		}

		return saldoAkhir, nil
	}

	saldomonthnow.UsedBKK += totalOut
	saldomonthnow.SaldoIn += totalIn
	saldomonthnow.SaldoAkhir = saldomonthnow.SaldoAkhir - totalOut + totalIn

	if err = a.UpdateMonth(saldomonthnow.ID, saldomonthnow); err != nil {
		return 0, err
	}

	return saldoAkhir, nil
}

func (a SaldoService) Update(id string, saldo *models.Saldo) error {
	oSaldo, err := a.Get(id)
	if err != nil {
		return err
	}
	saldo.ID = oSaldo.ID

	if err := a.saldoRepository.Update(id, saldo); err != nil {
		return err
	}

	return nil
}

func (a SaldoService) UpdateMonth(id string, saldo *models.SaldoMonth) error {
	oSaldo, err := a.GetMonthbyID(id)
	if err != nil {
		return err
	}
	saldo.ID = oSaldo.ID

	if err := a.saldomonthRepository.Update(id, saldo); err != nil {
		return err
	}

	return nil
}

func (a SaldoService) Delete(id string) error {
	_, err := a.saldoRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.saldoRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a SaldoService) UpdateStatus(id string, status int) error {
	_, err := a.saldoRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.saldoRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
