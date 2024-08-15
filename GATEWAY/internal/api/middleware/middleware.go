package middleware

import (
	"fmt"
	"gateway/config"
	"gateway/internal/api/casbin"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Check(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")

		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		token, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{},
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(cfg.JWT.SecretKey), nil
			})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token could not be parsed",
			})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token provided",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid claims' type",
			})
			return
		}

		userRole, roleOk := claims["role"].(string)
		userID, idOk := claims["user_id"].(string)
		if !roleOk || !idOk {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		c.Set("user_id", userID)
		c.Set("user_role", userRole)

		e, err := casbin.CasbinEnforcer(cfg)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Policy could not be loaded",
			})
			return
		}

		ok, err = e.Enforce(userRole, c.Request.URL.Path, c.Request.Method)
		if !ok || err != nil {
			msg := fmt.Sprintf("Access denied: %s cannot %s %s",
				userRole, c.Request.Method, c.Request.URL.Path,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": msg,
			})
			return
		}

		c.Next()
	}
}
