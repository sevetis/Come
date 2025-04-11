package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/util"

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

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := repository.QueryUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, "user not authenticated")
	}

	var input struct {
		Username string `json:"username" binding:"required,max=50"`
		Email    string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request format: "+err.Error())
	}

	updates := map[string]any{
		"username": input.Username,
		"email":    input.Email,
	}
	if err := repository.UpdateUser(userID.(uint), updates); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update profile")
	}

	c.JSON(http.StatusOK, "profile updated successfully")
}

func UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, "user not authenticated")
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, "failed to upload avatar")
	}

	if file.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, "avatar file too large (max 5MB)")
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, "only JPG and PNG files are allowed")
	}

	user, err := repository.QueryUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to fetch user profile")
	}
	oldAvatarPath := user.Avatar

	filename := fmt.Sprintf("%d_%s%s", userID, util.GenerateRandomString(8), ext)
	savePath := filepath.Join("uploads/avatars", filename)

	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		c.JSON(http.StatusInternalServerError, "server error")
	}

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		os.Remove(savePath)
		c.JSON(http.StatusInternalServerError, "failed to save avatar")
	}

	if err := repository.UpdateUser(userID.(uint), map[string]any{"avatar": savePath}); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update avatar")
	}

	if oldAvatarPath != "" && oldAvatarPath != savePath {
		if err := os.Remove(oldAvatarPath); err != nil {
			fmt.Println("Failed to delete old avatar:", err)
		}
	}

	c.JSON(http.StatusOK, savePath)
}
