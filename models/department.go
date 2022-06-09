package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Department struct {
	database.Model
	database.ModelMaster
	ID        string  `gorm:"column:id;size:36;not null;index:idx_id_dep,unique;" json:"id"`
	Num       string  `gorm:"column:num;size:5;not null;index;" json:"num"`
	Name      string  `gorm:"column:name;not null;" json:"name" validate:"required"`
	Desc      string  `gorm:"column:desc;not null;" json:"desc"`
	CompanyID string  `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	Company   Company `gorm:"-" json:"company" yaml:"company"`
}

type Departments []*Department

type DepartmentQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	Num        string   `query:"num"`
	Name       string   `query:"name"`
	CompanyID  string   `query:"company_id"`
	Desc       string   `query:"desc"`
	ActiveFlag bool     `query:"active_flag"`
	QueryValue string   `query:"query_value"`
}

type DepartmentQueryResult struct {
	List       Departments     `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Departments) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}

	return names
}

func (a Departments) ToMap() map[string]*Department {
	m := make(map[string]*Department)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
