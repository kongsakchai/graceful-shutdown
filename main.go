package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		fmt.Println("Ping -> 10s")
		time.Sleep(10 * time.Second)
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// r.Run(":8080")

	serv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	idle := make(chan struct{})
	go gracefulShutdown(&serv, idle)

	if err := serv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idle
	fmt.Println("Server shutdown")
}

func gracefulShutdown(serv *http.Server, idle chan<- struct{}) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server shutdown gracefully")
	close(idle)
}
