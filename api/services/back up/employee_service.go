package services

import (
	"gorm.io/gorm"

	"github.com/Aguztinus/petty-cash-backend/api/repository"
	"github.com/Aguztinus/petty-cash-backend/errors"
	"github.com/Aguztinus/petty-cash-backend/lib"
	"github.com/Aguztinus/petty-cash-backend/models"
	"github.com/Aguztinus/petty-cash-backend/pkg/uuid"
)

// EmployeeService service layer
type EmployeeService struct {
	logger               lib.Logger
	casbinService        CasbinService
	userRepository       repository.UserRepository
	employeeRepository   repository.EmployeeRepository
	menuRepository       repository.MenuRepository
	menuActionRepository repository.MenuActionRepository
}

// NewEmployeeService creates a new employeeservice
func NewEmployeeService(
	logger lib.Logger,
	casbinService CasbinService,
	userRepository repository.UserRepository,
	employeeRepository repository.EmployeeRepository,
	menuRepository repository.MenuRepository,
	menuActionRepository repository.MenuActionRepository,
) EmployeeService {
	return EmployeeService{
		logger:               logger,
		casbinService:        casbinService,
		userRepository:       userRepository,
		employeeRepository:   employeeRepository,
		menuRepository:       menuRepository,
		menuActionRepository: menuActionRepository,
	}
}

// WithTrx delegates transaction to repository database
func (a EmployeeService) WithTrx(trxHandle *gorm.DB) EmployeeService {
	a.employeeRepository = a.employeeRepository.WithTrx(trxHandle)
	a.userRepository = a.userRepository.WithTrx(trxHandle)

	return a
}

func (a EmployeeService) Query(param *models.EmployeeQueryParam) (employeeQR *models.EmployeeQueryResult, err error) {
	return a.employeeRepository.Query(param)
}

func (a EmployeeService) Get(id string) (*models.Employee, error) {
	employee, err := a.employeeRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (a EmployeeService) Check(item *models.Employee) error {
	qr, err := a.employeeRepository.Query(&models.EmployeeQueryParam{
		Name:      item.Name,
		CompanyID: item.CompanyID,
		BranchID:  item.BranchID,
	})

	if err != nil {
		return err
	} else if len(qr.List) > 0 {
		return errors.EmployeeAlreadyExists
	}

	return nil
}

func (a EmployeeService) Create(employee *models.Employee) (id string, err error) {
	if err = a.Check(employee); err != nil {
		return
	}

	employee.ID = uuid.MustString()

	if err = a.employeeRepository.Create(employee); err != nil {
		return
	}
	return employee.ID, nil
}

func (a EmployeeService) Update(id string, employee *models.Employee) error {
	oEmployee, err := a.Get(id)
	if err != nil {
		return err
	}
	employee.ID = oEmployee.ID

	if err := a.employeeRepository.Update(id, employee); err != nil {
		return err
	}

	return nil
}

func (a EmployeeService) Delete(id string) error {
	_, err := a.employeeRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.employeeRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (a EmployeeService) UpdateStatus(id string, status int) error {
	_, err := a.employeeRepository.Get(id)
	if err != nil {
		return err
	}

	if err := a.employeeRepository.UpdateStatus(id, status); err != nil {
		return err
	}

	return nil
}
