package controller

import (
	"net/http"
	"os"
	"user-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func UserAuth() gin.HandlerFunc {
	return jwtAuth(model.RoleUser)
}

func AdminAuth() gin.HandlerFunc {
	return jwtAuth(model.RoleAdmin)
}

func jwtAuth(allow_role model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "auth token missing"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claim"})
			return
		}

		user_id := uint(claims["user_id"].(float64))
		role := model.UserRole(int(claims["role"].(float64)))
		if role < allow_role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Set("user_id", user_id)
		c.Set("role", role)
		c.Next()
	}
}
