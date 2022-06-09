package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type BKKDetail struct {
	database.Model
	database.ModelTrans
	BKKHeaderID string            `gorm:"column:bkk_header_id;size:36;index;not null;" json:"bkk_header_id"`
	TrxID       string            `gorm:"column:trx_id;size:36;index;not null;" json:"trx_id"`
	LinesDesc   string            `gorm:"column:lines_desc;not null;" json:"lines_desc" validate:"required"`
	LinesDate   database.Datetime `gorm:"column:lines_date;" json:"lines_date"`
	LinesAmount int64             `gorm:"column:lines_amount;default:0;" json:"lines_amount"`
	LinesFile   string            `gorm:"column:lines_file;not null;" json:"lines_file" validate:"required"`
	Status      string            `gorm:"column:status;size:1;index;not null;" json:"status"`
	BKKHeader   BKKHeader         `gorm:"-" json:"bkk_header" yaml:"bkk_header"`
	Trx         Trx               `gorm:"foreignKey:TrxID;references:ID" json:"trx" yaml:"trx"`
}

type BKKDetails []*BKKDetail

type BKKDetailQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs          []string          `query:"ids"`
	BKKHeaderIDs []string          `query:"bkk_header_ids"`
	BKKHeaderID  string            `query:"bkk_header_id"`
	TrxID        string            `query:"trx_id"`
	LinesDesc    string            `query:"lines_desc"`
	LinesDate    database.Datetime `query:"paid_date"`
	LinesAmount  int64             `query:"limit_amount"`
	LinesFile    string            `query:"lines_file"`
	Status       string            `query:"status"`
	QueryValue   string            `query:"query_value"`
}

type BKKDetailQueryResult struct {
	List       BKKDetails      `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
