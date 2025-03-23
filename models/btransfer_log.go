package models

import "time"

type Btransferlog struct {
	ID          uint      `gorm:"primary key;autoIncrement" json:"id"`
	FromWho     uint      `json:"from_who"`
	ToWho       uint      `json:"to_who"`
	Price       uint      `json:"price"`
	CreatedDate time.Time `json:"created_date"`
}
