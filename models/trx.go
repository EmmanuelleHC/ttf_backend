package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Trx struct {
	database.Model
	database.ModelMaster
	ID             string     `gorm:"column:id;size:36;not null;index:idx_id_trx,unique;" json:"id"`
	Name           string     `gorm:"column:name;not null;" json:"name" validate:"required"`
	CompanyID      string     `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID       string     `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	CCID           string     `gorm:"column:cc_id;size:36;index;not null;" json:"cc_id"`
	AccID          string     `gorm:"column:account_id;size:36;index;not null;" json:"account_id"`
	DeptID         string     `gorm:"column:department_id;size:36;index;not null;" json:"department_id"`
	SegmentedValue string     `gorm:"column:segmented_value;not null;" json:"segmented_value" validate:"required"`
	Company        Company    `gorm:"references:ID" json:"company" yaml:"company"`
	Branch         Branch     `gorm:"references:ID" json:"branch" yaml:"branch"`
	CostCentre     CostCentre `gorm:"foreignKey:CCID;references:ID" json:"costcenter" yaml:"costcenter"`
	Account        Account    `gorm:"foreignKey:AccID;references:ID" json:"account" yaml:"account"`
	Department     Department `gorm:"foreignKey:DeptID;references:ID" json:"department" yaml:"department"`
}

type Trxs []*Trx

type TrxQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs            []string `query:"ids"`
	Name           string   `query:"name"`
	CompanyID      string   `query:"company_id"`
	BranchID       string   `query:"branch_id"`
	CCID           string   `query:"cc_id"`
	AccID          string   `query:"account_id"`
	DeptID         string   `query:"department_id"`
	SegmentedValue string   `query:"segmented_value"`
	ActiveFlag     bool     `query:"active_flag"`
	QueryValue     string   `query:"query_value"`
}

type TrxQueryResult struct {
	List       Trxs            `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Trxs) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}

	return names
}

func (a Trxs) ToMap() map[string]*Trx {
	m := make(map[string]*Trx)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
