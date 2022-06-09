package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// DepartmentService service layer
type DepartmentService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	departmentRepository repository.DepartmentRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewDepartmentService creates a new departmentservice
func NewDepartmentService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	departmentRepository repository.DepartmentRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) DepartmentService {
	return DepartmentService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		departmentRepository: departmentRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a DepartmentService) WithTrx(trxHandle *gorm.DB) DepartmentService {
	a.departmentRepository = a.departmentRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a DepartmentService) Query(param *models.DepartmentQueryParam) (departmentQR *models.DepartmentQueryResult, err error) {
	return a.departmentRepository.Query(param)
}

func (a DepartmentService) Get(id string) (*models.Department, error) {
	department, err := a.departmentRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return department, nil
}

func (a DepartmentService) Check(item *models.Department) error {
	qr, err := a.departmentRepository.Query(&models.DepartmentQueryParam{Name: item.Name})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.DepartmentAlreadyExists
	}

	return nil
}

func (a DepartmentService) Create(department *models.Department) (id string, err error) {
	if err = a.Check(department); err != nil {
		return
	}

	department.ID = uuid.MustString()

	if err = a.departmentRepository.Create(department); err != nil {
		return
	}
	return department.ID, nil
}

func (a DepartmentService) Update(id string, department *models.Department) error {
	oDepartment, err := a.Get(id)
	if err != nil {
		return err
	} else if department.Name != oDepartment.Name {
		if err = a.Check(department); err != nil {
			return err
		}
	}
	department.ID = oDepartment.ID

	if err := a.departmentRepository.Update(id, department); err != nil {
		return err
	}

	return nil
}

func (a DepartmentService) Delete(id string) error {
	_, err := a.departmentRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.departmentRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a DepartmentService) UpdateStatus(id string, status int) error {
	_, err := a.departmentRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.departmentRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
