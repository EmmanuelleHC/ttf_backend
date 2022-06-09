package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type InvoiceHeader struct {
	database.Model
	database.ModelTrans
	ID             string            `gorm:"column:id;size:36;not null;index;" json:"id"`
	Num            string            `gorm:"column:num;size:10;not null;index:idx_num,unique;" json:"num"`
	Type           string            `gorm:"column:type;size:15;index;not null;" json:"type"`
	Amount         int64             `gorm:"column:amount;default:0;" json:"amount"`
	Description    string            `gorm:"column:description;not null;" json:"description"`
	Date           database.Datetime `gorm:"column:date;" json:"date"`
	File           string            `gorm:"column:file;not null;" json:"file"`
	SisaAmount     int64             `gorm:"column:sisa_amount;default:0;" json:"sisa_amount"`
	Status         string            `gorm:"column:status;size:15;index;not null;" json:"status"`
	StatusApprove  int8              `gorm:"column:status_approve;default:0;" json:"status_approve"`
	CompanyID      string            `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID       string            `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	InvoiceDetails InvoiceDetails    `gorm:"foreignKey:InvoiceHeaderID;references:ID" json:"invoice_detail" yaml:"invoice_detail"`
	Company        Company           `gorm:"references:ID" json:"company" yaml:"company"`
	Branch         Branch            `gorm:"references:ID" json:"branch" yaml:"branch"`
}

type InvoiceHeaders []*InvoiceHeader

type InvoiceHeaderQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs           []string `query:"ids"`
	Num           string   `query:"num"`
	Type          string   `query:"type"`
	Amount        int64    `query:"amount"`
	Description   string   `query:"description"`
	Date          string   `query:"date"`
	File          string   `query:"file"`
	Status        string   `query:"status"`
	StatusApprove string   `query:"status_approve"`
	CompanyID     string   `query:"company_id"`
	BranchID      string   `query:"branch_id"`
	QueryValue    string   `query:"query_value"`
}

type InvoiceHeaderApproveQueryParam struct {
	IDs        []string `json:"ids"`
	QueryValue string   `json:"query_value"`
}

type InvoiceHeaderQueryResult struct {
	List       InvoiceHeaders  `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a InvoiceHeaders) ToMap() map[string]*InvoiceHeader {
	m := make(map[string]*InvoiceHeader)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
