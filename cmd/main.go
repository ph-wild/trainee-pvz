package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"trainee-pvz/config"
	"trainee-pvz/internal/auth"
	"trainee-pvz/internal/database"
	proto_pvz "trainee-pvz/internal/grpc"
	"trainee-pvz/internal/handler"
	"trainee-pvz/internal/metrics"
	"trainee-pvz/internal/repository"
	"trainee-pvz/internal/service"
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

	db, err := database.ConnectDB(ctx, cfg.DB.Connection)
	if err != nil {
		slog.Error("failed to connect", slog.Any("err", err))
		<-ctx.Done()
		return
	}
	defer db.Close()

	m := metrics.InitMetrics()

	userRepo := repository.NewUserRepository(db)
	receptionRepo := repository.NewReceptionRepository(db)
	productRepo := repository.NewProductRepository(db)
	PVZRepo := repository.NewPVZRepository(db)

	userService := service.NewUserService(userRepo)
	receptionService := service.NewReceptionService(receptionRepo, m)
	productService := service.NewProductService(productRepo, m)
	PVZService := service.NewPVZService(PVZRepo, m)

	services := handler.Services{
		User:      userService,
		Product:   productService,
		Reception: receptionService,
		PVZ:       PVZService,
	}

	expiration := time.Duration(cfg.Auth.JWTExpirationMinutes) * time.Minute
	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, expiration)

	server := handler.NewServer(services, jwtManager, cfg, m)
	router := server.Routes()

	go metrics.RunMetricServer(cfg.Prometheus.Port)

	go func() {
		err = proto_pvz.StartGRPCServer(PVZRepo, fmt.Sprintf(":%s", cfg.GRPC.Port))
		if err != nil {
			slog.Error("Can't start GRPC:", slog.Any("error", err))
		}
	}()
	slog.Info("Starting server on", slog.String("port", cfg.GRPC.Port))

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
