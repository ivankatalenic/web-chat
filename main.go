package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ivankatalenic/web-chat/internal/models"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ivankatalenic/web-chat/internal/config"
	"github.com/ivankatalenic/web-chat/internal/impl/client"
	"github.com/ivankatalenic/web-chat/internal/impl/logger"
	"github.com/ivankatalenic/web-chat/internal/interfaces"
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
	broadcaster := services.NewBroadcaster(logger.NewPrefix(log, "BROADCASTER"))

	go broadcaster.Start()

	tlsRouter := gin.Default()

	tlsRouter.GET("/favicon.ico", func(c *gin.Context) {
		c.File("web/favicon.ico")
	})

	authorized := tlsRouter.Group("/", gin.BasicAuth(gin.Accounts{
		config.Auth.Username: config.Auth.Password,
	}))

	authorized.GET("", func(c *gin.Context) {
		c.File("web/index.html")
	})

	websocketUpgrader := websocket.Upgrader{
		ReadBufferSize:  32,
		WriteBufferSize: 32,
	}
	authorized.GET("/chat", func(context *gin.Context) {
		if !context.IsWebsocket() {
			context.Status(http.StatusBadRequest)
			return
		}

		conn, err := websocketUpgrader.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			log.Error("Failed to upgrade to a WebSocket:\n\t" + err.Error())
			context.Status(http.StatusBadRequest)
			return
		}

		c := client.NewWebSocket(conn)
		initChatClient(c, log, repo, broadcaster)
		context.Status(http.StatusOK)
	})

	tlsServer := &http.Server{
		Addr:    ":https",
		Handler: tlsRouter,
	}
	go func() {
		tlsManager := services.NewCertificateManager(config.TLS)
		if err := tlsServer.ListenAndServeTLS(
			tlsManager.CertFilePath,
			tlsManager.KeyFilePath,
		); err != nil && err != http.ErrServerClosed {
			log.Error(err.Error())
		}
	}()

	redirectRouter := gin.Default()

	redirectRouter.GET("*catchAll", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://" + config.Server.Host + c.Param("catchAll"))
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

func initChatClient(
	client interfaces.Client,
	log interfaces.Logger,
	repo interfaces.MessageRepository,
	broadcaster *services.Broadcaster) {

	var err error

	// Send last n messages
	msgs, err := repo.GetLast(10)
	for _, msg := range msgs {
		err := client.SendMessage(msg)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	err = broadcaster.AddClient(client)
	if err != nil {
		log.Info(err.Error())
		return
	}

	go serveChatClient(client, log, repo, broadcaster)
}

func serveChatClient(
	client interfaces.Client,
	log interfaces.Logger,
	repo interfaces.MessageRepository,
	broadcaster *services.Broadcaster) {

	addr := client.GetAddress()
	for {
		if client.IsDisconnected() {
			break
		}

		msg, err := client.GetMessage()
		if err != nil {
			log.Error(err.Error())
			break
		}

		if len(msg.Author) == 0 {
			err := client.SendMessage(&models.Message{
				Author:    "SERVER",
				Content:   "Your message is missing the author",
				Timestamp: time.Now(),
			})
			if err != nil {
				log.Error(err.Error())
				break
			}
			continue
		}

		msg.Timestamp = time.Now()

		log.Info("New message: [" + addr + "] " + msg.Author + ": " + msg.Content)

		err = repo.Put(msg)
		if err != nil {
			log.Error(err.Error())
		}

		err = broadcaster.BroadcastMessage(msg)
		if err != nil {
			log.Error(err.Error())
		}
	}
}
