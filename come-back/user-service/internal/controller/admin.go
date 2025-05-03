package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"user-service/internal/model"
	"user-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	users, err := repository.QueryAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

func BanUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user ID")
		return
	}

	var input struct {
		Banned bool `json:"banned" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request format")
		return
	}

	if err := repository.UpdateUser(uint(userID), map[string]interface{}{"banned": input.Banned}); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update user")
		return
	}
	status := "banned"
	if !input.Banned {
		status = "unbanned"
	}
	c.JSON(http.StatusOK, fmt.Sprintf("user %s successfully", status))
}

func PromoteToAdmin(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid user ID")
		return
	}

	if err := repository.UpdateUser(uint(userID), map[string]interface{}{"role": model.RoleAdmin}); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to promote user")
		return
	}
	c.JSON(http.StatusOK, "user promoted to admin")
}
