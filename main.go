package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"location-history-server/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var appAddr string
	appAddr = os.Getenv("HISTORY_SERVER_LISTEN_ADDR")
	if appAddr == "" {
		appAddr = ":8080"
	}

	r := gin.Default()
	h := handler.NewService()

	r.POST("/location/:order_id/now", h.Create)
	r.GET("/location/:order_id", h.Get)
	r.DELETE("/location/:order_id", h.Delete)

	srv := &http.Server{
		Addr:    appAddr,
		Handler: r,
	}

	go func() {
		//service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 15 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {

		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
