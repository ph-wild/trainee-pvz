package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"log/slog"
	"net/http"

	"trainee-pvz/config"
	"trainee-pvz/internal/database"
	"trainee-pvz/internal/handler"
)

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()

	cfg, err := config.GetConfig("config.yaml")
	if err != nil {
		slog.Error("Can't read config.yaml", slog.Any("error", err))
		<-ctx.Done()
		slog.Info("Shutting down after config failure")
		return
	}

	db := database.ConnectDB(cfg.DB.Connection)
	defer db.Close()

	server := handler.NewServer()
	router := server.Routes()

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTP.Port), router)
		if err != nil {
			slog.Error("Can't start service:", slog.Any("error", err))
		}
	}()
	slog.Info("Starting server on", slog.String("port", cfg.HTTP.Port))

	<-ctx.Done()
	slog.Info("Got shutdown signal, exit program")
}
