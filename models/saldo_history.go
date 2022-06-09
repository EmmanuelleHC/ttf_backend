package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type SaldoHistory struct {
	database.Model
	database.ModelMaster
	ID         string  `gorm:"column:id;size:36;not null;index;" json:"id"`
	Desc       string  `gorm:"column:desc;not null;index:idx_saldo_his_desc" json:"desc"`
	CompanyID  string  `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID   string  `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	SaldoAwal  int64   `gorm:"column:saldo_awal;default:0;" json:"saldo_awal"`
	InAmount   int64   `gorm:"column:in_amount;default:0;" json:"in_amount"`
	OutAmount  int64   `gorm:"column:out_amount;default:0;" json:"out_amount"`
	SaldoAkhir int64   `gorm:"column:saldo_akhir;default:0;" json:"saldo_akhir"`
	Company    Company `gorm:"-" json:"company" yaml:"company"`
	Branch     Branch  `gorm:"-" json:"branch" yaml:"branch"`
}

type SaldoHistories []*SaldoHistory

type SaldoHistoryQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	Desc       string   `query:"desc"`
	CompanyID  string   `query:"company_id"`
	BranchID   string   `query:"branch_id"`
	SaldoAwal  int64    `query:"saldo_awal"`
	InAmount   int64    `query:"in_amount"`
	OutAmount  int64    `query:"out_amount"`
	ActiveFlag bool     `query:"active_flag"`
	SaldoAkhir int64    `query:"saldo_akhir"`
	DateQuery  []string `query:"date_query"`
	QueryValue string   `query:"query_value"`
}

type SaldoHistoryQueryResult struct {
	List       SaldoHistories  `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a SaldoHistories) ToMap() map[string]*SaldoHistory {
	m := make(map[string]*SaldoHistory)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
