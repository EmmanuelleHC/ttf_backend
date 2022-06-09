package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Branch struct {
	database.Model
	database.ModelMaster
	ID        string  `gorm:"column:id;size:36;not null;index:idx_id_branch,unique;" json:"id"`
	Code      string  `gorm:"column:code;size:5;not null;index;" json:"code"`
	Name      string  `gorm:"column:name;not null;" json:"name" validate:"required"`
	Shorter   string  `gorm:"column:shorter;size:5;not null;" json:"shorter" validate:"required"`
	CompanyID string  `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	RegOrgID  int     `gorm:"column:reg_org_id;" json:"reg_org_id"`
	Company   Company `gorm:"-" json:"company" yaml:"company"`
}

type Branchs []*Branch

type BranchQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	Code       string   `query:"code"`
	Name       string   `query:"name"`
	Shorter    string   `query:"shorter"`
	CompanyID  string   `query:"company_id"`
	RegOrgID   int      `query:"reg_org_id"`
	ActiveFlag bool     `query:"active_flag"`
	QueryValue string   `query:"query_value"`
}

type BranchQueryResult struct {
	List       Branchs         `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Branchs) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}

	return names
}

func (a Branchs) ToMap() map[string]*Branch {
	m := make(map[string]*Branch)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
