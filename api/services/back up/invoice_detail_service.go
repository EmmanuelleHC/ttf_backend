package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
)

// InvoiceDetailService service layer
type InvoiceDetailService struct {
	logger                  lib.Logger
	casbinService           CasbinService
	invociedetailRepository repository.InvoiceDetailRepository
	bkkheaderRepository     repository.BKKHeaderRepository
	bkkdetailRepository     repository.BKKDetailRepository
	redis                   lib.Redis
}

// NewInvoiceDetailService creates a new invociedetailservice
func NewInvoiceDetailService(
	redis lib.Redis,
	logger lib.Logger,
	casbinService CasbinService,
	invociedetailRepository repository.InvoiceDetailRepository,
	bkkheaderRepository repository.BKKHeaderRepository,
	bkkdetailRepository repository.BKKDetailRepository,
) InvoiceDetailService {
	return InvoiceDetailService{
		logger:                  logger,
		casbinService:           casbinService,
		invociedetailRepository: invociedetailRepository,
		bkkheaderRepository:     bkkheaderRepository,
		bkkdetailRepository:     bkkdetailRepository,
		redis:                   redis,
	}
}

// WithTrx delegates transaction to repository database
func (a InvoiceDetailService) WithTrx(trxHandle *gorm.DB) InvoiceDetailService {
	a.invociedetailRepository = a.invociedetailRepository.WithTrx(trxHandle)

	return a
}

func (a InvoiceDetailService) Query(param *models.InvoiceDetailQueryParam) (invociedetailQR *models.InvoiceDetailQueryResult, err error) {
	return a.invociedetailRepository.Query(param)
}

func (a InvoiceDetailService) QueryBkk(param *models.InvoiceDetailQueryParam) (bkkheaderQR *models.BKKHeaderQueryResult, err error) {
	res, err := a.invociedetailRepository.Query(param)
	if err != nil {
		return nil, err
	}
	var bkkheader []string
	for _, item := range res.List {
		bkkheader = append(bkkheader, item.BKKHeaderID)
	}

	return a.bkkheaderRepository.Query(&models.BKKHeaderQueryParam{IDs: bkkheader})
}

func (a InvoiceDetailService) Get(id string) (*models.InvoiceDetail, error) {
	invociedetail, err := a.invociedetailRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return invociedetail, nil
}

func (a InvoiceDetailService) Check(item *models.InvoiceDetail) error {
	qr, err := a.invociedetailRepository.Query(&models.InvoiceDetailQueryParam{BKKHeaderID: item.BKKHeaderID})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.InvoiceDetailAlreadyExists
	}

	return nil
}

func (a InvoiceDetailService) Create(invociedetail *models.InvoiceDetail) (id uint, err error) {
	if err = a.Check(invociedetail); err != nil {
		return
	}

	if err = a.invociedetailRepository.Create(invociedetail); err != nil {
		return
	}
	return invociedetail.RecordID, nil
}

func (a InvoiceDetailService) Update(id string, invociedetail *models.InvoiceDetail) error {
	oInvoiceDetail, err := a.Get(id)
	if err != nil {
		return err
	} else if invociedetail.BKKHeaderID != oInvoiceDetail.BKKHeaderID {
		if err = a.Check(invociedetail); err != nil {
			return err
		}
	}

	if err := a.invociedetailRepository.Update(id, invociedetail); err != nil {
		return err
	}

	return nil
}

func (a InvoiceDetailService) Delete(id string) error {
	_, err := a.invociedetailRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.invociedetailRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a InvoiceDetailService) UpdateStatus(id string, status int) error {
	_, err := a.invociedetailRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.invociedetailRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
