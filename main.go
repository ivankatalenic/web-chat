package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ivankatalenic/web-chat/interfaces"
	"github.com/ivankatalenic/web-chat/models"
	"github.com/ivankatalenic/web-chat/services"
	"net/http"
	"time"
)

func main() {
	gin.DisableConsoleColor()

	router := gin.Default()

	log := services.GetLogger()
	repo := services.GetMessageRepository()
	broadcaster := services.NewBroadcaster(log)

	broadcastCtx, broadcastCancel := context.WithCancel(context.Background())
	defer broadcastCancel()
	go broadcaster.Start(broadcastCtx)

	websocketUpgrader := websocket.Upgrader{
		ReadBufferSize:  32,
		WriteBufferSize: 32,
	}

	router.GET("", func(context *gin.Context) {
		context.File("index.html")
	})

	router.GET("/chat", func(context *gin.Context) {
		if context.IsWebsocket() {
			conn, err := websocketUpgrader.Upgrade(context.Writer, context.Request, nil)
			if err != nil {
				log.Error("Failed to upgrade to a WebSocket:\n\t" + err.Error())
				return
			}

			processWebSocket(conn, log, repo, broadcaster)
			return
		}

		context.Status(http.StatusBadRequest)
	})

	_ = router.Run(":80")
}

func processWebSocket(conn *websocket.Conn, log interfaces.Logger, repo interfaces.MessageRepository, broadcaster *services.Broadcaster) {
	var err error

	// Send last n messages
	msgs, err := repo.GetLast(10)
	for _, msg := range msgs {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Error("Cannot send a JSON\n\t" + err.Error())
			_ = conn.Close()
			return
		}
	}

	err = broadcaster.AddConn(conn)
	if err != nil {
		log.Info(err.Error())
		return
	}

	go func(conn *websocket.Conn, log interfaces.Logger, repo interfaces.MessageRepository, broadcaster *services.Broadcaster) {
		addr := conn.RemoteAddr().String()
		for {
			var msg models.Message
			err := conn.ReadJSON(&msg)

			if _, isCloseError := err.(*websocket.CloseError); isCloseError {
				broadcaster.RemoveConn(conn)
				break
			}

			if err != nil {
				log.Error("Cannot read a JSON\n\t" + err.Error())
				broadcaster.RemoveConn(conn)
				break
			}

			log.Info("[" + addr + "] " + msg.Author + ": " + msg.Content)

			msg.Timestamp = time.Now()
			err = repo.Put(&msg)
			if err != nil {
				log.Error("Cannot put message in a repository\n\t" + err.Error())
			}

			broadcaster.SendMessage(&msg)
		}
	}(conn, log, repo, broadcaster)

}
