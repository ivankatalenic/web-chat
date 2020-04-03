package client

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/ivankatalenic/web-chat/internal/models"
	"time"
)

// WebSocket holds infromation about a WebSocket chat client
type WebSocket struct {
	conn           *websocket.Conn
	addr           string
	isDisconnected bool
}

// NewWebSocket creates a new WebSocket chat client
func NewWebSocket(conn *websocket.Conn) *WebSocket {
	if conn == nil {
		return nil
	}
	return &WebSocket{
		conn: conn,
		addr: conn.RemoteAddr().String(),
	}
}

// GetAddress getter
func (client *WebSocket) GetAddress() string {
	return client.addr
}

// SendMessage sends message to a client
func (client *WebSocket) SendMessage(message *models.Message) error {
	if client.isDisconnected {
		return errors.New("client is disconnected")
	}
	return client.conn.WriteJSON(message)
}

// GetMessage returns the last unread message from a client
func (client *WebSocket) GetMessage() (*models.Message, error) {
	if client.isDisconnected {
		return nil, errors.New("client is disconnected")
	}
	msg := new(models.Message)
	err := client.conn.ReadJSON(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Disconnect disconnects client
func (client *WebSocket) Disconnect() error {
	_ = client.conn.WriteControl(
		websocket.CloseNormalClosure,
		[]byte("Disconnecting the client"),
		time.Now().Add(100*time.Millisecond),
	)
	_ = client.conn.Close()
	client.isDisconnected = true
	return nil
}

// IsDisconnected getter
func (client *WebSocket) IsDisconnected() bool {
	return client.isDisconnected
}
