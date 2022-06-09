package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type BKKHeader struct {
	database.Model
	database.ModelTrans
	ID            string            `gorm:"column:id;size:36;not null;index;" json:"id"`
	Num           string            `gorm:"column:num;size:10;not null;index:idx_num,unique;" json:"num"`
	NumberSeq     int64             `gorm:"column:number_seq;default:0;" json:"number_seq"`
	CompanyID     string            `gorm:"column:company_id;size:36;index;not null;" json:"company_id"`
	BranchID      string            `gorm:"column:branch_id;size:36;index;not null;" json:"branch_id"`
	ReleaseDate   database.Datetime `gorm:"column:release_date;" json:"release_date"`
	PaidDate      database.Datetime `gorm:"column:paid_date;" json:"paid_date"`
	TotalAmount   int64             `gorm:"column:total_amount;default:0;" json:"total_amount"`
	KasbonID      string            `gorm:"column:kasbon_id;size:36;index;not null;" json:"kasbon_id"`
	InvoiceID     string            `gorm:"column:invoice_id;size:36;index;not null;" json:"invoice_id"`
	Status        string            `gorm:"column:status;size:15;index;not null;" json:"status"`
	StatusApprove int8              `gorm:"column:status_approve;default:0;" json:"status_approve"`
	BKKDetails    BKKDetails        `gorm:"foreignKey:BKKHeaderID;references:ID" json:"bkk_detail" yaml:"bkk_detail"`
	Company       Company           `gorm:"references:ID" json:"company" yaml:"company"`
	Branch        Branch            `gorm:"references:ID" json:"branch" yaml:"branch"`
}

type BKKHeaders []*BKKHeader

type BKKHeaderQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs           []string          `query:"ids"`
	Num           string            `query:"num"`
	NumberSeq     int64             `query:"number_seq"`
	CompanyID     string            `query:"company_id"`
	BranchID      string            `query:"branch_id"`
	ReleaseDate   database.Datetime `query:"release_date"`
	PaidDate      database.Datetime `query:"paid_date"`
	KasbonID      string            `query:"kasbon_id"`
	InvoiceID     string            `query:"invoice_id"`
	TotalAmount   int64             `query:"total_amount"`
	Status        string            `query:"status"`
	StatusApprove string            `query:"status_approve"`
	DateQuery     []string          `query:"date_query"`
	UserId        string            `query:"user_id"`
	QueryValue    string            `query:"query_value"`
}

type BKKHeaderApproveQueryParam struct {
	IDs        []string `json:"ids"`
	QueryValue string   `json:"query_value"`
}

type BKKHeaderQueryResult struct {
	List       BKKHeaders      `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}

func (a BKKHeaders) ToMap() map[string]*BKKHeader {
	m := make(map[string]*BKKHeader)
	for _, item := range a {
		m[item.ID] = item
	}

	return m
}
