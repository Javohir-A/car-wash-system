package middleware

import (
	"gateway/config"
	"gateway/genproto/auth"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	User *auth.UserResponse `json:"user"`
	jwt.RegisteredClaims
}

func PermmissonChecker() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strToken := ctx.Request.Header.Get("Authorization")
		log.Println("Authorization Header:", strToken)

		token, err := ExtractToken(strToken)
		if err != nil || !token.Valid {
			log.Println("Invalid token:", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			ctx.Abort()
			return
		}

		obj := ctx.Request.URL.Path
		act := ctx.Request.Method

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			role := claims.User.Role

			log.Println("sub: ", role, "obj:", obj, "act:", act)

			enforcer, err := casbin.NewEnforcer("./internal/api/casbin/model.conf",
				"./internal/api/casbin/policy.csv")
			if err != nil {
				log.Println("Failed to create enforcer:", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
				ctx.Abort()
				return
			}

			ok, err := enforcer.Enforce(role, obj, act)
			if err != nil {
				log.Println("Enforcement error:", err)
				ctx.JSON(http.StatusForbidden, gin.H{"error": "not allowed"})
				ctx.Abort()
				return
			}
			if !ok {
				log.Println("Access denied for role:", role)
				ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func ExtractToken(str string) (*jwt.Token, error) {
	var tokenClaim Claims
	cnf := config.NewConfig()
	if err := cnf.Load(); err != nil {
		log.Println("Failed to load config:", err)
		return nil, err
	}

	secretKey := cnf.JWT.SecretKey

	token, err := jwt.ParseWithClaims(str, &tokenClaim, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Println("Failed to parse token:", err)
		return nil, err
	}

	return token, nil
}
