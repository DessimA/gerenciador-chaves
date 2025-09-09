package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dessima/gerenciador-chaves-api/infrastructure/config"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/database"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/http/router"
	"github.com/dessima/gerenciador-chaves-api/infrastructure/repository"
	"github.com/dessima/gerenciador-chaves-api/usecase"
)

func main() {
	cfg := config.Load()

	dbClient, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Falha ao conectar com o banco de dados:", err)
	}
	defer func() {
		if err = dbClient.Disconnect(context.Background()); err != nil {
			log.Fatal("Falha ao desconectar do banco de dados:", err)
		}
	}()

	// Initialize Repositories
	db := dbClient.Database(cfg.DatabaseName)
	keyRepo := repository.NewKeyRepository(db)
	userRepo := repository.NewUserRepository(db)
	reservationRepo := repository.NewReservationRepository(db, dbClient) // Passa o client para suportar transações

	// Initialize Use Cases
	keyUseCase := usecase.NewKeyUseCase(keyRepo, reservationRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, reservationRepo)
	reservationUseCase := usecase.NewReservationUseCase(reservationRepo, keyRepo, userRepo)

	// Setup Router
	r := router.Setup(cfg, userUseCase, keyUseCase, reservationUseCase)

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Falha ao iniciar servidor:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Falha no graceful shutdown:", err)
	}
}
