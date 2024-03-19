package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"dam/config"
	"dam/handlers"
	"dam/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("create logger error: %s", err.Error()))
	}

	config.LoadConfig()
	db, err := gorm.Open(postgres.Open(config.Cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		logger.Sugar().Errorf("connect database error: %s", err.Error())
		return
	}
	rdClient := redis.NewClient(&redis.Options{Addr: config.Cfg.Redis.Addr()})

	// gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	userHandler := handlers.NewUserHandler(db, rdClient)

	router := gin.Default()

	router.POST("/users/login", userHandler.Login)
	router.POST("/users/logout", userHandler.Logout)
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/me", middlewares.Authentication(rdClient), userHandler.GetCurrentUser)
	router.PUT("/users/me", middlewares.Authentication(rdClient), userHandler.UpdateUser)

	// TODO: add ping and health

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Cfg.Application.Port),
		Handler: router,
	}
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Sugar().Errorf("shutdown http.Server error: %s", err.Error())
			return
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Sugar().Errorf("ListenAndServe error: %s", err.Error())
		return
	}
}
