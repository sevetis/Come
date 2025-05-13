package controller

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"chat-service/internal/client"
	"chat-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatController struct {
	UserClient client.UserClient
}

func NewChatController(userClinet client.UserClient) *ChatController {
	return &ChatController{UserClient: userClinet}
}

func init() {
	StartWebSocketManager()
}

func (cc *ChatController) HandleChat(c *gin.Context) {
	var user model.User
	userID, exists := c.Get("user_id")
	if exists {
		var err error
		user, err = cc.UserClient.GetUserByID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			return
		}
	} else {
		user = model.User{
			ID:       0,
			Username: fmt.Sprintf("anonymous_%d", rand.Intn(100000)),
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("websocket upgrade failed:", err)
		return
	}
	log.Println("New WebSocket connection established for user:", user.Username)

	client := &Client{conn: conn, user: user, closed: false}
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()

	joinMsg := model.ChatMessage{
		UserID:    user.ID,
		Username:  user.Username,
		Content:   fmt.Sprintf("%s come", user.Username),
		Timestamp: time.Now().Unix(),
		Type:      model.JoinType,
	}
	BroadcastMessage(joinMsg)

	done := make(chan struct{})
	defer func() {
		close(done)
		client.Close()
	}()

	go func() {
		<-done
		clientsMu.Lock()
		delete(clients, client)
		log.Println("Client disconnected:", user.Username)
		leaveMsg := model.ChatMessage{
			UserID:    user.ID,
			Username:  user.Username,
			Content:   fmt.Sprintf("%s leave", user.Username),
			Timestamp: time.Now().Unix(),
			Type:      model.LeaveType,
		}
		BroadcastMessage(leaveMsg)
		clientsMu.Unlock()
	}()

	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	conn.SetPongHandler(func(string) error {
		log.Println("Pong received from", user.Username)
		conn.SetReadDeadline(time.Now().Add(15 * time.Second))
		return nil
	})

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := client.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Println("Ping failed, closing connection:", err)
					client.Close()
					return
				}
			case <-done:
				return
			}
		}
	}()

	type Content struct {
		Content string `json:"content"`
	}

	for {
		var content Content
		err := conn.ReadJSON(&content)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Printf("Client %s closed connection normally", user.Username)
			} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Read timeout for %s, closing connection", user.Username)
			} else {
				log.Printf("read error for %s: %v", user.Username, err)
			}
			break
		}

		msg := model.ChatMessage{
			UserID:    user.ID,
			Username:  user.Username,
			Content:   content.Content,
			Timestamp: time.Now().Unix(),
			Type:      model.MessageType,
		}

		BroadcastMessage(msg)
		StoreMessage(c.Request.Context(), msg)
	}
}

func (cc *ChatController) GetChatHistory(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 50
	}

	messages, err := GetChatHistory(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chat history"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": messages})
}

func (cc *ChatController) GetOnlineCount(c *gin.Context) {
	clientsMu.Lock()
	count := len(clients)
	clientsMu.Unlock()
	c.JSON(http.StatusOK, gin.H{"count": count})
}
