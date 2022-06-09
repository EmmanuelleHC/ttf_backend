package services

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// InvoiceHeaderService service layer
type InvoiceHeaderService struct {
	logger                  lib.Logger
	casbinService           CasbinService
	counterService          CounterService
	saldoService            SaldoService
	userService             UserService
	bkkheaderService        BKKHeaderService
	userRepository          repository.UserRepository
	branchRepository        repository.BranchRepository
	invoiceheaderRepository repository.InvoiceHeaderRepository
	invociedetailRepository repository.InvoiceDetailRepository
	bkkheaderRepository     repository.BKKHeaderRepository
	counterRepository       repository.CounterRepository
	saldoRepository         repository.SaldoRepository
	saldohistoryRepository  repository.SaldoHistoryRepository
	saldomonthRepository    repository.SaldoMonthRepository
}

// NewInvoiceHeaderService creates a new invoiceheaderservice
func NewInvoiceHeaderService(
	logger lib.Logger,
	casbinService CasbinService,
	counterService CounterService,
	saldoService SaldoService,
	userService UserService,
	bkkheaderService BKKHeaderService,
	userRepository repository.UserRepository,
	branchRepository repository.BranchRepository,
	invoiceheaderRepository repository.InvoiceHeaderRepository,
	invociedetailRepository repository.InvoiceDetailRepository,
	bkkheaderRepository repository.BKKHeaderRepository,
	counterRepository repository.CounterRepository,
	saldoRepository repository.SaldoRepository,
	saldohistoryRepository repository.SaldoHistoryRepository,
	saldomonthRepository repository.SaldoMonthRepository,
) InvoiceHeaderService {
	return InvoiceHeaderService{
		logger:                  logger,
		casbinService:           casbinService,
		counterService:          counterService,
		saldoService:            saldoService,
		userService:             userService,
		bkkheaderService:        bkkheaderService,
		userRepository:          userRepository,
		branchRepository:        branchRepository,
		invoiceheaderRepository: invoiceheaderRepository,
		invociedetailRepository: invociedetailRepository,
		bkkheaderRepository:     bkkheaderRepository,
		counterRepository:       counterRepository,
		saldoRepository:         saldoRepository,
		saldohistoryRepository:  saldohistoryRepository,
		saldomonthRepository:    saldomonthRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a InvoiceHeaderService) WithTrx(trxHandle *gorm.DB) InvoiceHeaderService {
	a.invoiceheaderRepository = a.invoiceheaderRepository.WithTrx(trxHandle)
	a.bkkheaderRepository = a.bkkheaderRepository.WithTrx(trxHandle)
	a.saldoRepository = a.saldoRepository.WithTrx(trxHandle)
	a.saldohistoryRepository = a.saldohistoryRepository.WithTrx(trxHandle)
	a.saldomonthRepository = a.saldomonthRepository.WithTrx(trxHandle)
	a.counterRepository = a.counterRepository.WithTrx(trxHandle)

	return a
}

func (a InvoiceHeaderService) Query(param *models.InvoiceHeaderQueryParam) (invoiceheaderQR *models.InvoiceHeaderQueryResult, err error) {
	return a.invoiceheaderRepository.Query(param)
}

func (a InvoiceHeaderService) Get(id string) (*models.InvoiceHeader, error) {
	invoiceheader, err := a.invoiceheaderRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return invoiceheader, nil
}

func (a InvoiceHeaderService) Check(item *models.InvoiceHeader) error {
	qr, err := a.invoiceheaderRepository.Query(&models.InvoiceHeaderQueryParam{Num: item.Num})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.InvoiceHeaderAlreadyExists
	}

	return nil
}

func (a InvoiceHeaderService) Create(invoiceheader *models.InvoiceHeader) (id string, err error) {
	cbn, err := a.branchRepository.Get(invoiceheader.BranchID)
	if err != nil {
		return "", err
	}

	key := "INV" + cbn.Shorter
	count, err := a.counterService.GetAndIncrement(key)
	padleft := fmt.Sprintf("%04d", count)
	if err != nil {
		return "", err
	}

	invoiceheader.ID = uuid.MustString()
	invoiceheader.Num = key + padleft

	if err = a.invoiceheaderRepository.Create(invoiceheader); err != nil {
		return "", err
	}
	// update bkk status to invoice
	for _, item := range invoiceheader.InvoiceDetails {
		if err = a.bkkheaderService.UpdateStatus(item.BKKHeaderID, "Invoice", 0); err != nil {
			return "", err
		}
	}

	return invoiceheader.Num, nil
}

func (a InvoiceHeaderService) Update(id string, invoiceheader *models.InvoiceHeader) error {
	oInvoiceHeader, err := a.Get(id)
	if err != nil {
		return err
	} else if invoiceheader.Num != oInvoiceHeader.Num {
		if err = a.Check(invoiceheader); err != nil {
			return err
		}
	}
	invoiceheader.ID = oInvoiceHeader.ID

	if err := a.invoiceheaderRepository.Update(id, invoiceheader); err != nil {
		return err
	}

	if err := a.invoiceheaderRepository.UpdateApprove(id, 0); err != nil {
		return err
	}

	// update bkk status to invoice
	for _, item := range invoiceheader.InvoiceDetails {
		if err = a.bkkheaderService.UpdateStatus(item.BKKHeaderID, "Invoice", 0); err != nil {
			return err
		}
	}

	return nil
}

func (a InvoiceHeaderService) Delete(id string) error {
	_, err := a.invoiceheaderRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.invoiceheaderRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a InvoiceHeaderService) UpdateStatus(id string, status int) error {
	_, err := a.invoiceheaderRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.invoiceheaderRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}

func (a InvoiceHeaderService) UpdateApprove(ids []string, status int) error {
	reject := "2,4"
	approveFinal := 3
	if strings.Contains(reject, strconv.Itoa(status)) {
		for _, id := range ids {
			_, err := a.invociedetailRepository.Query(&models.InvoiceDetailQueryParam{BKKHeaderID: id})
			if err != nil {
				return err
			}
			// for _, item := range dtl.List {
			// 	if err := a.bkkheaderService.UpdateStatus(item.BKKHeaderID, "Paid", 1); err != nil {
			// 		return err
			// 	}
			// }
		}
	} else if status == approveFinal {
		hdrs, err := a.invoiceheaderRepository.Query(&models.InvoiceHeaderQueryParam{IDs: ids})
		if err != nil {
			return err
		}
		for _, hdr := range hdrs.List {
			if hdr.SisaAmount > 0 {
				saldoNow, err := a.saldoService.CreateNewSaldoOrUpdate(hdr.CompanyID, hdr.BranchID, 0, hdr.SisaAmount)

				if err != nil {
					return err
				}

				saldoHisCreate := new(models.SaldoHistory)
				saldoHisCreate.Desc = "Pengembalian Dana Invoice: " + hdr.Num
				saldoHisCreate.CompanyID = hdr.CompanyID
				saldoHisCreate.BranchID = hdr.BranchID
				saldoHisCreate.InAmount = hdr.SisaAmount
				saldoHisCreate.SaldoAkhir = saldoNow
				if err = a.saldohistoryRepository.Create(saldoHisCreate); err != nil {
					return err
				}
			}
		}
	}

	if err := a.invoiceheaderRepository.UpdateApproves(ids, status); err != nil {
		return err
	}

	return nil
}
