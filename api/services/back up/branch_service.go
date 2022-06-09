package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// BranchService service layer
type BranchService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	branchRepository     repository.BranchRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewBranchService creates a new branchservice
func NewBranchService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	branchRepository repository.BranchRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) BranchService {
	return BranchService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		branchRepository:     branchRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a BranchService) WithTrx(trxHandle *gorm.DB) BranchService {
	a.branchRepository = a.branchRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a BranchService) Query(param *models.BranchQueryParam) (branchQR *models.BranchQueryResult, err error) {
	return a.branchRepository.Query(param)
}

func (a BranchService) Get(id string) (*models.Branch, error) {
	branch, err := a.branchRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return branch, nil
}

func (a BranchService) Check(item *models.Branch) error {
	qr, err := a.branchRepository.Query(&models.BranchQueryParam{Name: item.Name})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.BranchAlreadyExists
	}

	return nil
}

func (a BranchService) Create(branch *models.Branch) (id string, err error) {
	if err = a.Check(branch); err != nil {
		return
	}

	branch.ID = uuid.MustString()

	if err = a.branchRepository.Create(branch); err != nil {
		return
	}
	return branch.ID, nil
}

func (a BranchService) Update(id string, branch *models.Branch) error {
	oBranch, err := a.Get(id)
	if err != nil {
		return err
	} else if branch.Name != oBranch.Name {
		if err = a.Check(branch); err != nil {
			return err
		}
	}
	branch.ID = oBranch.ID

	if err := a.branchRepository.Update(id, branch); err != nil {
		return err
	}

	return nil
}

func (a BranchService) Delete(id string) error {
	_, err := a.branchRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.branchRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a BranchService) UpdateStatus(id string, status int) error {
	_, err := a.branchRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.branchRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
