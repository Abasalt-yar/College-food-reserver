package models

import "time"

type Student struct {
	ID          uint64    `gorm:"primary key;autoIncrement" json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Position    string    `json:"position"`
	CollegeId   uint64    `json:"college_id"`
	AddedBy     uint64    `json:"added_by"`
	Counter     int8      `json:"counter"`
	Balance     int16     `json:"balance"`
	SecondPass  *string   `json:"second_pass"`
	CreatedDate time.Time `json:"created_date" gorm:"default:CURRENT_TIMESTAMP()"`
}
