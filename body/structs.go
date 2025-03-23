package body

import (
	"github.com/abasalt-yar/college-food-reserver/config"

	"github.com/gin-gonic/gin"
)

type BodyParse interface {
	Validate(c *gin.Context) *config.ResponseError
}

type StudentLoginBody struct {
	Username string
	Password string
}
type LoginTokenResponse struct {
	Status      bool
	AccessToken string `json:"access_token"`
}

type StudentTBalanceBody struct {
	Target  string
	Balance int32
}
type StudentTBalanceResponse struct {
	Status     bool
	To         string
	NewBalance int16
}

type StudentCPassBody struct {
	Current string
	New     string
}

type StudentCPassResponse struct {
	Status  bool
	Message string
	Logout  bool
}

type StudentC2Pass struct {
	New string
}

type StudentAddBalanceResponse struct {
	Status    bool
	Authority string
}

type StudentVerifyPayment struct {
	Authority string
}

type StudentVerifyPaymentResponse struct {
	Status bool
	RefID  uint64
	Price  uint64
}

type AdminLoginBody struct {
	Username string
	Password string
}

type AdminChangePassword struct {
	Current string
	New     string
}

type AdminCPassResponse struct {
	Status  bool
	Message string
	Logout  bool
}

type AdminAddFoodBody struct {
	FoodName    string `json:"food_name"`
	IsExpensive bool   `json:"is_expensive"`
}
