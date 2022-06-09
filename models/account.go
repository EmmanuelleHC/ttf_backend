package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Account struct {
	database.Model
	database.ModelMaster
	ID        string  `gorm:"column:id;size:36;not null;index:idx_id_acc,unique;" json:"id"`
	Num       string  `gorm:"column:num;size:10;not null;index;" json:"num"`
	Name      string  `gorm:"column:name;not null;" json:"name" validate:"required"`
	CompanyID string  `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	Company   Company `gorm:"-" json:"company" yaml:"company"`
}

type Accounts []*Account

type AccountQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	Num        string   `query:"num"`
	Name       string   `query:"name"`
	CompanyID  string   `query:"company_id"`
	ActiveFlag bool     `query:"active_flag"`
	QueryValue string   `query:"query_value"`
}

type AccountQueryResult struct {
	List       Accounts        `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Accounts) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}

	return names
}

func (a Accounts) ToMap() map[string]*Account {
	m := make(map[string]*Account)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
