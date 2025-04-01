package middleware

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/abasalt-yar/college-food-reserver/config"
	"github.com/abasalt-yar/college-food-reserver/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddlewareStudent(c *gin.Context) {
	header := strings.Split(c.GetHeader("authorization"), " ")
	if len(header) != 2 || header[0] != "Bearer" || len(header[1]) < 10 {
		c.JSON(401, config.ResponseError{
			Status:  false,
			Message: "AUTHORIZATION HEADER NOT FOUND",
		})
		return
	}
	token, err := jwt.Parse(header[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.JSON(401, config.ResponseError{
			Status:  false,
			Message: "Authentication Failed.",
		})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		admin_token, ok := claims["admin"].(bool)
		if !ok || admin_token {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed",
			})
			return
		}
		var username string = claims["username"].(string)
		matches, err := regexp.MatchString(`(?i)^\d{14}$`, username)
		if err != nil || !matches {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed",
			})
			return
		}
		db := (c.Keys["PSQLDB"].(*gorm.DB))
		student := models.Student{}
		sub, er := claims.GetSubject()
		if er != nil {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed.",
			})
			return
		}
		newsub, er := strconv.ParseUint(sub, 10, 64)
		if er != nil {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed",
			})
			return
		}
		if err := db.First(&student, &models.Student{ID: newsub}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(401, config.ResponseError{
					Status:  false,
					Message: "Authentication Failed.",
				})
				return
			}
			go config.CustomError(config.CErrorOptions{
				Err:   err.Error(),
				Level: sentry.LevelFatal,
			})
			c.JSON(500, config.ResponseError{
				Status:  false,
				Message: "SERVER_INTERNAL_ERROR",
			})
			return
		}
		exp, err := claims.GetExpirationTime()
		var inc int8 = int8(claims["incremental"].(float64))
		if student.Counter != inc || err != nil || exp.Before(time.Now()) {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Token Expired.",
			})
			return
		}
		c.Set("Student", &student)

		c.Next()
	} else {
		c.JSON(401, config.ResponseError{
			Status:  false,
			Message: "Authentication Failed.",
		})
		return
	}
}

func AuthMiddlewareAdmin(c *gin.Context) {
	header := strings.Split(c.GetHeader("authorization"), " ")
	if len(header) != 2 || header[0] != "Bearer" || len(header[1]) < 10 {
		c.JSON(401, config.ResponseError{
			Status:  false,
			Message: "AUTHORIZATION HEADER NOT FOUND",
		})
		return
	}
	token, err := jwt.Parse(header[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ADMIN_SECRET")), nil
	})
	if err != nil {
		c.JSON(401, config.ResponseError{
			Status:  false,
			Message: "Authentication Failed.",
		})
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		admin_token, ok := claims["admin"].(bool)
		if !ok || !admin_token {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed",
			})
			return
		}
		var username string = claims["username"].(string)
		matches, err := regexp.MatchString(`(?i)^([A-Z]|[0-9]){8,20}$`, username)
		if err != nil || !matches {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed",
			})
			return
		}
		db := (c.Keys["PSQLDB"].(*gorm.DB))
		admin := models.Admin{}
		sub, er := claims.GetSubject()
		if er != nil {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed.",
			})
			return
		}
		newsub, er := strconv.ParseUint(sub, 10, 64)
		if er != nil {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Authentication Failed",
			})
			return
		}
		if err := db.First(&admin, &models.Admin{ID: newsub}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(401, config.ResponseError{
					Status:  false,
					Message: "Authentication Failed.",
				})
				return
			}
			go config.CustomError(config.CErrorOptions{
				Err:   err.Error(),
				Level: sentry.LevelFatal,
			})
			c.JSON(500, config.ResponseError{
				Status:  false,
				Message: "SERVER_INTERNAL_ERROR",
			})
			return
		}
		exp, err := claims.GetExpirationTime()
		var inc int8 = int8(claims["incremental"].(float64))
		if admin.Counter != inc || err != nil || exp.Before(time.Now()) {
			c.JSON(401, config.ResponseError{
				Status:  false,
				Message: "Token Expired.",
			})
			return
		}
		c.Set("Admin", &admin)

		c.Next()
	} else {
		c.JSON(401, config.ResponseError{
			Status:  false,
			Message: "Authentication Failed.",
		})
		return
	}
}
