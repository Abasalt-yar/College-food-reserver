package models

import "time"

type Admin struct {
	ID          uint      `gorm:"primary key;autoincrement" json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"Password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Position    string    `json:"position"`
	CollegeId   uint      `json:"college_id"`
	Counter     int8      `json:"counter"`
	CreatedDate time.Time `json:"created_date"`
}

type Food struct {
	ID          uint      `gorm:"primary key;autoincrement" json:"id"`
	FoodName    string    `json:"food_name"`
	Expensive   bool      `json:"expensive"`
	CollegeId   uint      `json:"college_id"`
	AddedBy     uint      `json:"added_by"`
	CreatedDate time.Time `json:"created_date"`
}
