package controller

import (
	"come-back/user-service/internal/model"
	"come-back/user-service/internal/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid user id")
	}
	user, err := repository.QueryUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetUsersBatch(c *gin.Context) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		c.JSON(http.StatusBadRequest, "ids parameter is required")
		return
	}

	ids := strings.Split(idsStr, ",")
	var userIDs []uint
	for _, id := range ids {
		if uid, err := strconv.ParseUint(id, 10, 32); err == nil {
			userIDs = append(userIDs, uint(uid))
		}
	}

	users, err := repository.QueryUsers(userIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to fetch users")
		return
	}

	userMap := make(map[uint]model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	c.JSON(http.StatusOK, userMap)
}
