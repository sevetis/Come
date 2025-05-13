package client

import (
	"chat-service/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserClient struct {
	Host string
}

func NewUserClient(host string) *UserClient {
	return &UserClient{
		Host: host,
	}
}

func (uc *UserClient) GetUserByID(userID uint) (model.User, error) {
	url := fmt.Sprintf("%s/api/user/%d", uc.Host, userID)
	resp, err := http.Get(url)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.User{}, fmt.Errorf("failed to get user, status code: %d", resp.StatusCode)
	}

	var user model.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return model.User{}, fmt.Errorf("failed to decode user response: %w", err)
	}

	return user, nil
}
