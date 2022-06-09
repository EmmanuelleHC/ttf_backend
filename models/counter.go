package models

import (
	"github.com/Aguztinus/petty-cash-backend/models/database"
	"github.com/Aguztinus/petty-cash-backend/models/dto"
)

// Status - 1: Enable -1: Disable
type Counter struct {
	database.Model
	database.ModelMaster
	ID         string `gorm:"column:id;size:36;not null;index;" json:"id"`
	KeyCounter string `gorm:"column:key_counter;size:10;not null;index:idx_key_count,unique;" json:"key_counter"`
	CounterMe  int    `gorm:"column:counter_me;not null;" json:"counter_me" validate:"required"`
}

type Counters []*Counter

type CounterQueryParam struct {
	dto.PaginationParam
	dto.OrderParam

	IDs        []string `query:"ids"`
	KeyCounter string   `query:"key_counter"`
	CounterMe  string   `query:"counter_me"`
	QueryValue string   `query:"query_value"`
}

type CounterQueryResult struct {
	List       Counters        `json:"list"`
	Pagination *dto.Pagination `json:"pagination"`
}
