package models

import "time"

type Admin struct {
	ID          uint64    `gorm:"primary key;autoincrement" json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"Password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Position    string    `json:"position"`
	CollegeId   uint64    `json:"college_id"`
	Counter     int8      `json:"counter"`
	CreatedDate time.Time `json:"created_date" gorm:"default:CURRENT_TIMESTAMP()"`
}

type Food struct {
	ID          uint64    `gorm:"primary key;autoincrement" json:"id"`
	FoodName    string    `json:"food_name"`
	Expensive   bool      `json:"expensive"`
	CollegeId   uint64    `json:"college_id"`
	AddedBy     uint64    `json:"added_by"`
	CreatedDate time.Time `json:"created_date" gorm:"default:CURRENT_TIMESTAMP()"`
}
