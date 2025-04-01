package models

import "time"

type Btransferlog struct {
	ID          uint64    `gorm:"primary key;autoIncrement" json:"id"`
	FromWho     uint64    `json:"from_who"`
	ToWho       uint64    `json:"to_who"`
	Price       uint      `json:"price"`
	CreatedDate time.Time `json:"created_date" gorm:"default:CURRENT_TIMESTAMP()"`
}
