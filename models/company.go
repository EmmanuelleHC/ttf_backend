package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Company struct {
	database.Model
	database.ModelMaster
	ID           string `gorm:"column:id;size:36;not null;index:idx_id_company,unique;" json:"id"`
	Num          string `gorm:"column:num;size:5;not null;index;" json:"num"`
	Name         string `gorm:"column:name;not null;" json:"name" validate:"required"`
	Address      string `gorm:"column:address;not null;" json:"address" validate:"required"`
	PaymentFlag  bool   `gorm:"column:payment_flag;" json:"payment_flag"`
	BDCFlag      bool   `gorm:"column:bdc_flag;" json:"bdc_flag"`
	ApprovalFlag bool   `gorm:"column:approval_flag;" json:"approval_flag"`
}

type Companies []*Company

type CompanyQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs          []string `query:"ids"`
	Num          string   `query:"num"`
	Name         string   `query:"name"`
	Address      string   `query:"address"`
	PaymentFlag  bool     `query:"payment_flag"`
	BDCFlag      bool     `query:"bdc_flag"`
	ApprovalFlag bool     `query:"approval_flag"`
	ActiveFlag   bool     `query:"active_flag"`
	QueryValue   string   `query:"query_value"`
}

type CompanyQueryResult struct {
	List       Companies       `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a Companies) ToNames() []string {
	names := make([]string, len(a))
	for i, item := range a {
		names[i] = item.Name
	}

	return names
}

func (a Companies) ToMap() map[string]*Company {
	m := make(map[string]*Company)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
