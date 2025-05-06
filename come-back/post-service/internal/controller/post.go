package controller

import (
	"log"
	"net/http"
	"post-service/internal/model"
	"post-service/internal/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllPost(c *gin.Context) {
	users, err := repository.QueryAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusAccepted, users)
	}
}

func GetPostsPaginated(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	posts, total, err := repository.QueryPostsPaginated(offset, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Failed to fetch posts")
	}

	c.JSON(
		http.StatusOK,
		map[string]any{
			"posts":    posts,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		})

}

func GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid post ID")
		return
	}

	post, err := repository.QueryPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, "post not found")
		return
	}

	c.JSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, "user not authenticated")
		return
	}

	// banned, err := repository.UserIsBanned(authorID.(uint))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "failed to validate author")
	// 	return
	// }
	// if banned {
	// 	c.JSON(http.StatusForbidden, "user banned")
	// 	return
	// }

	var postInput struct {
		Title   string `json:"title" binding:"required,max=50"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&postInput); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request format: "+err.Error())
		return
	}

	post := model.Post{
		Title:    postInput.Title,
		Content:  postInput.Content,
		AuthorID: authorID.(uint),
	}
	if repository.CreatePost(&post) != nil {
		c.JSON(http.StatusInternalServerError, "failed to save post")
		return
	}

	c.JSON(http.StatusCreated, "post successful")
}

func UpdatePost(c *gin.Context) {
	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, "user not authenticated")
		return
	}

	// banned, err := repository.UserIsBanned(authorID.(uint))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "failed to validate author")
	// 	return
	// }
	// if banned {
	// 	c.JSON(http.StatusForbidden, "user banned")
	// 	return
	// }

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid post ID")
		return
	}

	post, err := repository.QueryPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, "post not found")
		return
	}

	if post.AuthorID != authorID.(uint) {
		c.JSON(http.StatusForbidden, "you can only update your own posts")
		return
	}

	var postInput struct {
		Title   string `json:"title" binding:"required,max=50"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&postInput); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request format: "+err.Error())
		return
	}

	updates := map[string]any{
		"title":   postInput.Title,
		"content": postInput.Content,
	}

	if err := repository.UpdatePost(uint(postID), updates); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to update post")
		return
	}

	updatedPost, err := repository.QueryPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to fetch updated post")
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func DeletePost(c *gin.Context) {
	authorID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, "user not authenticated")
		return
	}

	// banned, err := repository.UserIsBanned(authorID.(uint))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "failed to validate author")
	// 	return
	// }
	// if banned {
	// 	c.JSON(http.StatusForbidden, "user banned")
	// 	return
	// }

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := repository.QueryPost(uint(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, "post not found")
		return
	}

	if post.AuthorID != authorID.(uint) {
		c.JSON(http.StatusForbidden, "you can only delete your own posts")
		return
	}

	if err := repository.DeletePost(uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to delete post")
		return
	}

	c.JSON(http.StatusOK, "post deleted successfully")
}

func GetPostComments(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("Invalid post ID: %v", err)
		c.JSON(http.StatusBadRequest, "invalid post ID")
		return
	}

	comments, err := repository.QueryAllComments(uint(postID))
	if err != nil {
		log.Printf("Failed to fetch comments for post %d: %v", postID, err)
		c.JSON(http.StatusInternalServerError, "failed to fetch comments")
		return
	}

	c.JSON(http.StatusOK, comments)
}

func CreateComment(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		log.Printf("Invalid post ID: %v", err)
		c.JSON(http.StatusBadRequest, "invalid post ID")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, "user not authenticated")
		return
	}

	// banned, err := repository.UserIsBanned(userID.(uint))
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, "failed to validate author")
	// 	return
	// }
	// if banned {
	// 	c.JSON(http.StatusForbidden, "user banned")
	// 	return
	// }

	var input struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	comment := model.Comment{
		PostID:   uint(postID),
		AuthorID: userID.(uint),
		Content:  input.Content,
	}

	if repository.CreateComment(&comment) != nil {
		log.Printf("Failed to create comment: %v", err)
		c.JSON(http.StatusInternalServerError, "failed to create comment")
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func DeletePostAdmin(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid post ID")
		return
	}

	if err := repository.DeletePost(uint(postID)); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to delete post")
		return
	}
	c.JSON(http.StatusOK, "post deleted successfully")
}

func DeleteCommentAdmin(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid comment ID")
		return
	}

	if err := repository.DeleteComment(uint(commentID)); err != nil {
		c.JSON(http.StatusInternalServerError, "failed to delete comment")
		return
	}
	c.JSON(http.StatusOK, "comment deleted successfully")
}
