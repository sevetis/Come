package controller

import (
	"net/http"
	"os"
	"time"
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var regReq struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&regReq); err != nil {
		c.JSON(http.StatusBadRequest, "wrong request format")
		return
	}
	if regReq.Email == "" || regReq.Username == "" || regReq.Password == "" {
		c.JSON(http.StatusBadRequest, "missing required registration information")
		return
	}

	if _, err := repository.QueryUserByEmail(regReq.Email); err == nil {
		c.JSON(http.StatusConflict, "email already in use")
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, "database error: "+err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(regReq.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to hash password: "+err.Error())
		return
	}

	user := model.User{
		Email:    regReq.Email,
		Username: regReq.Username,
		Password: string(hashedPassword),
	}

	if repository.CreateUser(&user) != nil {
		c.JSON(http.StatusInternalServerError, "failed to save user")
		return
	}

	c.JSON(http.StatusCreated, "register successful")
}

func Login(c *gin.Context) {
	var loginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	c.ShouldBindJSON(&loginReq)

	user, err := repository.QueryUserByEmail(loginReq.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, "account not exist")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusBadRequest, "wrong password")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECERT")))

	c.JSON(http.StatusAccepted, tokenString)
}
