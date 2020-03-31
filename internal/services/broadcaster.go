package services

import (
	"context"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/ivankatalenic/web-chat/internal/interfaces"
	"github.com/ivankatalenic/web-chat/internal/models"
	"sync"
	"time"
)

type Broadcaster struct {
	sendQueue chan models.Message
	connsLock sync.Mutex
	conns     map[string]*websocket.Conn

	ctx      context.Context
	stopFunc context.CancelFunc

	log interfaces.Logger
}

func NewBroadcaster(log interfaces.Logger) *Broadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	return &Broadcaster{
		sendQueue: make(chan models.Message, 32),
		connsLock: sync.Mutex{},
		conns:     make(map[string]*websocket.Conn),
		log:       log,
		ctx:       ctx,
		stopFunc:  cancel,
	}
}

func (b *Broadcaster) Start() {
	for {
		select {
		case msg := <-b.sendQueue:
			b.broadcast(msg)
		case <-b.ctx.Done():
			break
		}
	}

}

func (b *Broadcaster) broadcast(msg models.Message) {
	var removeConns []*websocket.Conn

	b.connsLock.Lock()
	defer b.connsLock.Unlock()

	for addr, conn := range b.conns {
		err := conn.WriteJSON(msg)
		if _, isCloseError := err.(*websocket.CloseError); isCloseError {
			b.log.Info("[" + addr + "] has disconnected")
			removeConns = append(removeConns, conn)
		}
		if err != nil {
			b.log.Error("Failed to write a JSON to [" + addr + "]:\n\t" + err.Error())
			removeConns = append(removeConns, conn)
		}
	}

	// Removing disconnected connections
	for _, conn := range removeConns {
		_ = conn.Close()
		delete(b.conns, conn.RemoteAddr().String())
	}
}

func (b *Broadcaster) RemoveConn(conn *websocket.Conn) {
	b.connsLock.Lock()
	defer b.connsLock.Unlock()
	delete(b.conns, conn.RemoteAddr().String())
	_ = conn.Close()
}

func (b *Broadcaster) AddConn(conn *websocket.Conn) error {
	if conn == nil {
		return nil
	}

	addr := conn.RemoteAddr().String()

	b.connsLock.Lock()
	defer b.connsLock.Unlock()
	_, ok := b.conns[addr]
	if ok {
		return errors.New("connection is already added")
	}

	b.conns[addr] = conn
	return nil
}

func (b *Broadcaster) SendMessage(msg *models.Message) {
	b.sendQueue <- *msg
}

func (b *Broadcaster) Stop() {
	b.connsLock.Lock()
	defer b.connsLock.Unlock()

	b.stopFunc()

	msg := models.Message{
		Author:    "SERVER",
		Content:   "Closing the connection: Server is shutting down!",
		Timestamp: time.Now(),
	}

	for _, conn := range b.conns {
		_ = conn.WriteJSON(msg)
		_ = conn.WriteControl(websocket.CloseGoingAway, nil, time.Now().Add(50*time.Millisecond))
		_ = conn.Close()
	}
}
