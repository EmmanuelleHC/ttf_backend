package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type TarikDana struct {
	database.Model
	database.ModelTrans
	ID          string            `gorm:"column:id;size:36;not null;index;" json:"id"`
	Type        string            `gorm:"column:type;size:15;index;not null;" json:"type"`
	Amount      int64             `gorm:"column:amount;default:0;" json:"amount"`
	Description string            `gorm:"column:description;not null;" json:"description"`
	Date        database.Datetime `gorm:"column:date;" json:"date"`
	File        string            `gorm:"column:file;not null;" json:"file"`
	CompanyID   string            `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID    string            `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`

	Company Company `gorm:"references:ID" json:"company" yaml:"company"`
	Branch  Branch  `gorm:"references:ID" json:"branch" yaml:"branch"`
}

type TarikDanas []*TarikDana

type TarikDanaQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs         []string `query:"ids"`
	Type        string   `query:"type"`
	Amount      int64    `query:"amount"`
	Description string   `query:"description"`
	Date        string   `query:"date"`
	File        string   `query:"file"`
	CompanyID   string   `query:"company_id"`
	BranchID    string   `query:"branch_id"`
	QueryValue  string   `query:"query_value"`
}

type TarikDanaQueryResult struct {
	List       TarikDanas      `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a TarikDanas) ToMap() map[string]*TarikDana {
	m := make(map[string]*TarikDana)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
