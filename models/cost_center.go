package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type CostCentre struct {
	database.Model
	database.ModelMaster
	ID        string  `gorm:"column:id;size:36;not null;index:idx_id_cc,unique;" json:"id"`
	Code      string  `gorm:"column:code;size:10;not null;index;" json:"code"`
	Name      string  `gorm:"column:name;not null;" json:"name" validate:"required"`
	CompanyID string  `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	Company   Company `gorm:"references:ID" json:"company" yaml:"company"`
}

type CostCentres []*CostCentre

type CostCentreQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	Code       string   `query:"code"`
	Name       string   `query:"name"`
	CompanyID  string   `query:"company_id"`
	ActiveFlag bool     `query:"active_flag"`
	QueryValue string   `query:"query_value"`
}

type CostCentreQueryResult struct {
	List       CostCentres     `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a CostCentres) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}

	return names
}

func (a CostCentres) ToMap() map[string]*CostCentre {
	m := make(map[string]*CostCentre)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
