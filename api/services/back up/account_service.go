package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// AccountService service layer
type AccountService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	accountRepository    repository.AccountRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewAccountService creates a new accountservice
func NewAccountService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	accountRepository repository.AccountRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) AccountService {
	return AccountService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		accountRepository:    accountRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a AccountService) WithTrx(trxHandle *gorm.DB) AccountService {
	a.accountRepository = a.accountRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a AccountService) Query(param *models.AccountQueryParam) (accountQR *models.AccountQueryResult, err error) {
	return a.accountRepository.Query(param)
}

func (a AccountService) Get(id string) (*models.Account, error) {
	account, err := a.accountRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a AccountService) Check(item *models.Account) error {
	qr, err := a.accountRepository.Query(&models.AccountQueryParam{Name: item.Name})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.AccountAlreadyExists
	}

	return nil
}

func (a AccountService) Create(account *models.Account) (id string, err error) {
	if err = a.Check(account); err != nil {
		return
	}

	account.ID = uuid.MustString()

	if err = a.accountRepository.Create(account); err != nil {
		return
	}
	return account.ID, nil
}

func (a AccountService) Update(id string, account *models.Account) error {
	oAccount, err := a.Get(id)
	if err != nil {
		return err
	} else if account.Name != oAccount.Name {
		if err = a.Check(account); err != nil {
			return err
		}
	}
	account.ID = oAccount.ID

	if err := a.accountRepository.Update(id, account); err != nil {
		return err
	}

	return nil
}

func (a AccountService) Delete(id string) error {
	_, err := a.accountRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.accountRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a AccountService) UpdateStatus(id string, status int) error {
	_, err := a.accountRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.accountRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
