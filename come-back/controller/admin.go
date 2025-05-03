package controller

import (
	"come-back/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DeletePostAdmin(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Error(http.StatusBadRequest, "invalid post ID"))
		return
	}

	if err := repository.DeletePost(uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, Error(http.StatusInternalServerError, "failed to delete post"))
		return
	}
	c.JSON(http.StatusOK, Success(http.StatusOK, "post deleted successfully"))
}

func DeleteCommentAdmin(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, Error(http.StatusBadRequest, "invalid comment ID"))
		return
	}

	if err := repository.DeleteComment(uint(commentID)); err != nil {
		c.JSON(http.StatusInternalServerError, Error(http.StatusInternalServerError, "failed to delete comment"))
		return
	}
	c.JSON(http.StatusOK, Success(http.StatusOK, "comment deleted successfully"))
}

func AdminDashboard(c *gin.Context) {
	usersCount, _ := repository.CountUsers()
	postsCount, _ := repository.CountPosts()
	commentsCount, _ := repository.CountComments()

	data := map[string]any{
		"users_count":    usersCount,
		"posts_count":    postsCount,
		"comments_count": commentsCount,
	}
	c.JSON(http.StatusOK, Success(http.StatusOK, data))
}
