package repository

import "post-service/internal/model"

func QueryAllComments(postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := dB.Where("post_id = ?", postID).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func CreateComment(comment *model.Comment) error {
	return dB.Create(comment).Error
}

func DeleteComment(commentID uint) error {
	return dB.Where("id = ?", commentID).Delete(&model.Comment{}).Error
}

func CountComments() (int64, error) {
	var count int64
	err := dB.Model(&model.Comment{}).Count(&count).Error
	return count, err
}
