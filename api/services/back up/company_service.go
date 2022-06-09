package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// CompanyService service layer
type CompanyService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	companyRepository    repository.CompanyRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewCompanyService creates a new companyservice
func NewCompanyService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	companyRepository repository.CompanyRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) CompanyService {
	return CompanyService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		companyRepository:    companyRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a CompanyService) WithTrx(trxHandle *gorm.DB) CompanyService {
	a.companyRepository = a.companyRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a CompanyService) Query(param *models.CompanyQueryParam) (companyQR *models.CompanyQueryResult, err error) {
	return a.companyRepository.Query(param)
}

func (a CompanyService) Get(id string) (*models.Company, error) {
	company, err := a.companyRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return company, nil
}

func (a CompanyService) Check(item *models.Company) error {
	qr, err := a.companyRepository.Query(&models.CompanyQueryParam{Name: item.Name})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.CompanyAlreadyExists
	}

	return nil
}

func (a CompanyService) Create(company *models.Company) (id string, err error) {
	if err = a.Check(company); err != nil {
		return
	}

	company.ID = uuid.MustString()

	if err = a.companyRepository.Create(company); err != nil {
		return
	}
	return company.ID, nil
}

func (a CompanyService) Update(id string, company *models.Company) error {
	oCompany, err := a.Get(id)
	if err != nil {
		return err
	} else if company.Name != oCompany.Name {
		if err = a.Check(company); err != nil {
			return err
		}
	}
	company.ID = oCompany.ID

	if err := a.companyRepository.Update(id, company); err != nil {
		return err
	}

	return nil
}

func (a CompanyService) Delete(id string) error {
	_, err := a.companyRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.companyRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a CompanyService) UpdateStatus(id string, status int) error {
	_, err := a.companyRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.companyRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
