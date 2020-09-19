package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"mail-sending/emails/welcome_mail"
	"mail-sending/handlers/welcome"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	appAddr := ":" + os.Getenv("APP_PORT")

	r := gin.Default()

	sendWelcomeMail := welcome_mail.NewService()

	welcomeMail := welcome.NewWelcome(sendWelcomeMail)

	r.POST("/welcome", welcomeMail.WelcomeMail)

	//Starting and Shutting down Server

	srv := &http.Server{
		Addr:    appAddr,
		Handler: r,
	}
	fmt.Println("App port: ", appAddr)
	go func() {
		//service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

}
