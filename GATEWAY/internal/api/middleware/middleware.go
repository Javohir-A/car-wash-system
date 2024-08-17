package middleware

import (
	"gateway/config"
	"gateway/genproto/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func OnlySudo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strToken, ok := ctx.Get("Authorization")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "making unauthorized request"})
			ctx.Abort()
			return
		}
		token, err := ExtractToken(strToken.(string))
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

	}

}

type TokenData struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	// Expiration  int    `json:"exp"`
	jwt.RegisteredClaims
}

type Claims struct {
	User *auth.UserResponse `json:"user"`
	jwt.RegisteredClaims
}

func ExtractToken(str string) (*jwt.Token, error) {
	var tokenClaim Claims
	cnf := config.NewConfig()
	if err := cnf.Load(); err != nil {
		log.Println(err)
		return nil, err
	}

	secretKey := cnf.JWT.SecretKey

	token, err := jwt.ParseWithClaims(str, &tokenClaim, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return token, nil
}
