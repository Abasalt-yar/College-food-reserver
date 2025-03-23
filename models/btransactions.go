package models

import "time"

type Btransaction struct {
	ID          uint      `gorm:"primary key;autoIncrement" json:"id"`
	Authority   string    `json:"authority"`
	RefID       uint      `json:"ref_id"`
	PayedBy     uint      `json:"payed_by"`
	Price       uint      `json:"price"`
	CreatedDate time.Time `json:"created_date"`
}
