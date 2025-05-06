package repository

import "post-service/internal/model"

func QueryAllPosts() ([]model.Post, error) {
	var posts []model.Post
	err := dB.Find(&posts).Error
	return posts, err
}

func QueryPost(postId uint) (model.Post, error) {
	var post model.Post
	err := dB.Where("id = ?", postId).First(&post).Error
	return post, err
}

func QueryPostsPaginated(offset, limit int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64
	if err := dB.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := dB.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func CreatePost(post *model.Post) error {
	return dB.Create(post).Error
}

func UpdatePost(postId uint, updates map[string]any) error {
	return dB.Model(&model.Post{}).Where("id = ?", postId).Updates(updates).Error
}

func DeletePost(postId uint) error {
	return dB.Where("id = ?", postId).Delete(&model.Post{}).Error
}

func CountPosts() (int64, error) {
	var count int64
	err := dB.Model(&model.Comment{}).Count(&count).Error
	return count, err
}
