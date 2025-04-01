package body

import (
	"regexp"

	"github.com/abasalt-yar/college-food-reserver/config"

	"github.com/gin-gonic/gin"
)

func (s *StudentLoginBody) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	matches, err := regexp.MatchString(`(?i)^\d{14}$`, s.Username)
	if err != nil || !matches {
		return &config.ResponseError{Status: false, Message: "DATA_IS_NOT_VALID"}
	}
	return nil
}

func (s *AdminLoginBody) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	matches, err := regexp.MatchString(`(?i)^([A-Z]|[0-9]){8,20}$`, s.Username)
	if err != nil || !matches {
		return &config.ResponseError{Status: false, Message: "DATA_IS_NOT_VALID"}
	}
	return nil
}

func (s *StudentTBalanceBody) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	matches, err := regexp.MatchString(`(?i)^\d{14}$`, s.Target)
	if err != nil || !matches {
		return &config.ResponseError{Status: false, Message: "TARGET_USERNAME_NOT_VALID"}
	}
	if s.Balance < 1000 || s.Balance > 100_000 || s.Balance%1000 != 0 {
		return &config.ResponseError{Status: false, Message: "TRANSFER_BALANCE_NOT_VALID"}
	}
	return nil
}

func (s *StudentCPassBody) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	new_length := len(s.New)
	cur_length := len(s.Current)
	if new_length < 8 {
		return &config.ResponseError{Status: false, Message: "NEWPASS_LENGTH_LOWER_8"}
	}
	if new_length > 60 {
		return &config.ResponseError{Status: false, Message: "NEWPASS_LENGTH_MORE_60"}
	}
	if cur_length < 8 || cur_length > 60 {
		return &config.ResponseError{Status: false, Message: "CURRENT_PASSWORD_INCORRECT"}
	}
	return nil
}

func (s *StudentC2Pass) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	new_len := len(s.New)
	if new_len < 4 || new_len > 8 {
		return &config.ResponseError{Status: false, Message: "NEWPASS_NOT_BETWEEN_4_8"}
	}
	matches, err := regexp.MatchString(`(?i)^\d{4,8}$`, s.New)
	if err != nil || !matches {
		return &config.ResponseError{Status: false, Message: "NEWPASS_NOT_MATCH"}
	}
	return nil
}

func (s *StudentVerifyPayment) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	matches, err := regexp.MatchString(`(?i)^([a-z]|[0-9]){36,36}$`, s.Authority)
	if err != nil || !matches {
		return &config.ResponseError{Status: false, Message: "INVALID_AUTHORITY"}
	}
	return nil

}

func (s *AdminChangePassword) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	new_length := len(s.New)
	cur_length := len(s.Current)
	if new_length < 8 {
		return &config.ResponseError{Status: false, Message: "NEWPASS_LENGTH_LOWER_8"}
	}
	if new_length > 60 {
		return &config.ResponseError{Status: false, Message: "NEWPASS_LENGTH_MORE_60"}
	}
	if cur_length < 8 || cur_length > 60 {
		return &config.ResponseError{Status: false, Message: "CURRENT_PASSWORD_INCORRECT"}
	}
	return nil
}

func (s *AdminAddFoodBody) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	food_length := len(s.FoodName)
	if food_length < 4 || food_length > 60 {
		return &config.ResponseError{Status: false, Message: "FOOD_NAME_BETWEEN_4_60"}
	}

	return nil
}

func (s *AdminAddStudent) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	if len(s.Password) < 5 || len(s.FirstName) < 3 || len(s.LastName) < 3 {
		return &config.ResponseError{Status: false, Message: "DATA_IS_NOT_VALID"}
	}
	if s.Position != "NORMAL_STUDENT" && s.Position != "DORM_STUDENT" {
		return &config.ResponseError{Status: false, Message: "DATA_IS_NOT_VALID"}
	}
	return nil
}

func (s *AdminUpdateStudentPosition) Validate(c *gin.Context) *config.ResponseError {
	if err := c.BindJSON(&s); err != nil {
		return &config.ResponseError{Status: false, Message: "UNPROCESSABLE_ENTITY"}
	}
	if s.Position != "NORMAL_STUDENT" && s.Position != "DORM_STUDENT" {
		return &config.ResponseError{Status: false, Message: "DATA_IS_NOT_VALID"}
	}
	return nil
}
