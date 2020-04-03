package services

import (
	"context"
	"errors"
	"github.com/ivankatalenic/web-chat/internal/interfaces"
	"github.com/ivankatalenic/web-chat/internal/models"
	"sync"
)

type Broadcaster struct {
	sendQueue chan models.Message

	clientMapLock sync.Mutex
	clientMap     map[string]interfaces.Client

	ctx      context.Context
	stopFunc context.CancelFunc

	log interfaces.Logger
}

func NewBroadcaster(log interfaces.Logger) *Broadcaster {
	ctx, cancel := context.WithCancel(context.Background())
	return &Broadcaster{
		sendQueue:     make(chan models.Message, 512),
		clientMapLock: sync.Mutex{},
		clientMap:     make(map[string]interfaces.Client),
		log:           log,
		ctx:           ctx,
		stopFunc:      cancel,
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
	var disconnectedClients []interfaces.Client

	b.clientMapLock.Lock()
	defer b.clientMapLock.Unlock()

	for _, client := range b.clientMap {
		if client.IsDisconnected() {
			disconnectedClients = append(disconnectedClients, client)
			continue
		}

		err := client.SendMessage(&msg)
		if err != nil {
			b.log.Error(err.Error())
		}
	}

	// Removing disconnected connections
	for _, client := range disconnectedClients {
		delete(b.clientMap, client.GetAddress())
	}
}

func (b *Broadcaster) AddClient(client interfaces.Client) error {
	if client == nil {
		return nil
	}

	addr := client.GetAddress()

	b.clientMapLock.Lock()
	defer b.clientMapLock.Unlock()

	_, ok := b.clientMap[addr]
	if ok {
		return errors.New("client is already added to broadcaster")
	}

	b.clientMap[addr] = client
	return nil
}

func (b *Broadcaster) BroadcastMessage(msg *models.Message) {
	b.sendQueue <- *msg
}

func (b *Broadcaster) Stop() {
	b.clientMapLock.Lock()
	defer b.clientMapLock.Unlock()

	b.stopFunc()

	for _, client := range b.clientMap {
		delete(b.clientMap, client.GetAddress())
	}
}
