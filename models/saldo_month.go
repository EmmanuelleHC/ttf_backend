package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type SaldoMonth struct {
	database.Model
	database.ModelMaster
	ID         string  `gorm:"column:id;size:36;not null;index;" json:"id"`
	CompanyID  string  `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID   string  `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	SaldoAwal  int64   `gorm:"column:saldo_awal;default:0;" json:"saldo_awal"`
	SaldoIn    int64   `gorm:"column:saldo_in;default:0;" json:"saldo_in"`
	UsedBKK    int64   `gorm:"column:used_bkk;default:0;" json:"used_bkk"`
	UsedKBS    int64   `gorm:"column:used_kbs;default:0;" json:"used_kbs"`
	SaldoAkhir int64   `gorm:"column:saldo_akhir;default:0;" json:"saldo_akhir"`
	MonthYear  string  `gorm:"column:month_year;size:10;index;not null;" json:"month_year"`
	Month      int     `gorm:"column:month;index;not null;" json:"month"`
	Year       int     `gorm:"column:year;index;not null;" json:"year"`
	Company    Company `gorm:"-" json:"company" yaml:"company"`
	Branch     Branch  `gorm:"-" json:"branch" yaml:"branch"`
}

type SaldoMonths []*SaldoMonth

type SaldoMonthQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	CompanyID  string   `query:"company_id"`
	BranchID   string   `query:"branch_id"`
	SaldoAwal  int64    `query:"saldo_awal"`
	SaldoIn    int64    `query:"saldo_in"`
	UsedBKK    int64    `query:"used_bkk"`
	UsedKBS    int64    `query:"used_kbs"`
	ActiveFlag bool     `query:"active_flag"`
	SaldoAkhir int64    `query:"saldo_akhir"`
	MonthYear  string   `query:"month_year"`
	Month      int      `query:"month"`
	Year       int      `query:"year"`
	QueryValue string   `query:"query_value"`
}

type SaldoMonthQueryResult struct {
	List       SaldoMonths     `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a SaldoMonths) ToMap() map[string]*SaldoMonth {
	m := make(map[string]*SaldoMonth)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
