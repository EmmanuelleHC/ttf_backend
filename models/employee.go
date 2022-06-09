package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Employee struct {
	database.Model
	database.ModelMaster
	ID        string  `gorm:"column:id;size:36;not null;index:idx_id_emp,unique;" json:"id"`
	CompanyID string  `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID  string  `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	Name      string  `gorm:"column:name;not null;" json:"name" validate:"required"`
	Company   Company `gorm:"-" json:"company" yaml:"company"`
	Branch    Branch  `gorm:"-" json:"branch" yaml:"branch"`
}

type Employees []*Employee

type EmployeeQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	CompanyID  string   `query:"company_id"`
	BranchID   string   `query:"branch_id"`
	Name       string   `query:"name"`
	ActiveFlag bool     `query:"active_flag"`
	QueryValue string   `query:"query_value"`
}

type EmployeeQueryResult struct {
	List       Employees       `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Employees) ToMap() map[string]*Employee {
	m := make(map[string]*Employee)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
