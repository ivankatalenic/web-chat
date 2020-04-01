package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ivankatalenic/web-chat/internal/interfaces"
	"github.com/ivankatalenic/web-chat/internal/models"
	"github.com/ivankatalenic/web-chat/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	gin.DisableConsoleColor()

	log := services.GetLogger()
	repo := services.GetMessageRepository()
	broadcaster := services.NewBroadcaster(log)

	go broadcaster.Start()

	websocketUpgrader := websocket.Upgrader{
		ReadBufferSize:  32,
		WriteBufferSize: 32,
	}

	tlsRouter := gin.Default()

	tlsRouter.GET("/favicon.ico", func(c *gin.Context) {
		c.File("web/favicon.ico")
	})

	authorized := tlsRouter.Group("/", gin.BasicAuth(gin.Accounts{
		"nyx": "jezvalilmuskepodiskacu",
	}))

	authorized.GET("", func(c *gin.Context) {
		c.File("web/index.html")
	})

	authorized.GET("/chat", func(context *gin.Context) {
		if !context.IsWebsocket() {
			context.Status(http.StatusBadRequest)
			return
		}

		conn, err := websocketUpgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade to a WebSocket:\n\t" + err.Error())
			return
		}

		processWebSocket(conn, log, repo, broadcaster)
	})

	tlsServer := &http.Server{
		Addr:    ":https",
		Handler: tlsRouter,
	}
	go func() {
		if err := tlsServer.ListenAndServeTLS(
			"/etc/letsencrypt/live/northcroatia.org/fullchain.pem",
			"/etc/letsencrypt/live/northcroatia.org/privkey.pem",
		); err != nil && err != http.ErrServerClosed {
			log.Error(err.Error())
		}
	}()

	redirectRouter := gin.Default()

	redirectRouter.GET("", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://northcroatia.org")
	})

	redirectRouter.GET("/favicon.ico", func(c *gin.Context) {
		c.File("web/favicon.ico")
	})

	redirectServer := &http.Server{
		Addr:    ":http",
		Handler: redirectRouter,
	}
	go func() {
		if err := redirectServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Wait for selected signals

	log.Warning("Shutting down the server!")
	broadcaster.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redirectServer.Shutdown(ctx); err != nil {
		log.Error("The server forced to shutdown: " + err.Error())
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := tlsServer.Shutdown(ctx); err != nil {
		log.Error("The server forced to shutdown: " + err.Error())
	}

	log.Info("The server shutdown is complete!")
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
