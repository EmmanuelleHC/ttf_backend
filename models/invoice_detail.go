package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type InvoiceDetail struct {
	database.Model
	database.ModelTrans
	BKKHeaderID     string        `gorm:"column:bkk_header_id;size:36;index;not null;" json:"bkk_header_id"`
	InvoiceHeaderID string        `gorm:"column:invoice_header_id;size:36;index;not null;" json:"invoice_header_id"`
	Status          string        `gorm:"column:status;size:1;index;" json:"status"`
	TotalAmount     int64         `gorm:"column:total_amount;default:0;" json:"total_amount"`
	BKKHeader       BKKHeader     `gorm:"-" json:"bkk_header" yaml:"bkk_header"`
	InvoiceHeader   InvoiceHeader `gorm:"foreignKey:InvoiceHeaderID;references:ID" json:"invoice_header" yaml:"invoice_header"`
}

type InvoiceDetails []*InvoiceDetail

type InvoiceDetailQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs             []string `query:"ids"`
	BKKHeaderID     string   `query:"bkk_header_id"`
	InvoiceHeaderID string   `query:"invoice_header_id"`
	Status          string   `query:"status"`
	QueryValue      string   `query:"query_value"`
}

type InvoiceDetailQueryResult struct {
	List       InvoiceDetails  `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
