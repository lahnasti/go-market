package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lahnasti/go-market/lib/rabbitmq"
	"github.com/lahnasti/go-market/products/internal/config"
	"github.com/lahnasti/go-market/products/internal/logger"
	"github.com/lahnasti/go-market/products/internal/repository"
	"github.com/lahnasti/go-market/products/internal/server"
	"github.com/lahnasti/go-market/products/internal/server/routes"

	"golang.org/x/sync/errgroup"
)

// @title Go-market
// @version 1.0
// @description This is a sample server for market.
// @host localhost:8080
// @BasePath /
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		fmt.Println("Received shutdown signal")
		<-c
		cancel()
	}()
	fmt.Println("Server starting")
	cfg := config.ReadConfig()
	zlog := logger.SetupLogger(cfg.DebugFlag)
	zlog.Debug().Any("config", cfg).Msg("Check cfg value")

	rabbit, err := rabbitmq.InitRabbit()
	if err != nil {
		zlog.Fatal().Err(err).Msg("RabbitMQ connection failed")
		return
	}
	defer rabbit.CloseRabbit()
	fmt.Println("Connected to RabbitMQ")
	// Создаём пул соединений
	pool, err := pgxpool.New(ctx, cfg.DBAddr) // используем cfg.DBAddr для подключения через пул
	if err != nil {
		fmt.Println("Failed to connect to PostgreSQL:", err)
		return
	}
	defer pool.Close()

	// Получаем соединение из пула
	conn, err := pool.Acquire(ctx) // Используем Acquire для получения соединения
	if err != nil {
		fmt.Println("Failed to acquire connection:", err)
		return
	}
	defer conn.Release() // Освобождаем соединение, когда оно не нужно

	// Проверяем и создаём базу данных для сервиса auth
	err = repository.EnsureMarketDatabaseExists(conn)
	if err != nil {
		fmt.Println("Failed to ensure auth database exists:", err)
		return
	}

	pool, err = initDB(cfg.DBAddr)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Connection DB failed")
	}
	defer pool.Close()

	err = repository.Migrations(cfg.DBAddr, cfg.MPath, zlog)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Init migrations failed")
	}

	dbStorage, err := repository.NewDB(pool)
	if err != nil {
		panic(err)
	}
	defer dbStorage.Close()

	group, gCtx := errgroup.WithContext(ctx)
	srv := server.NewServer(gCtx, dbStorage, zlog)
	group.Go(func() error {
		r := routes.SetupMarketRoutes(srv)
		zlog.Info().Msg("Server was started")

		if err := r.Run(cfg.Addr); err != nil {
			return err
		}
		return nil
	})

	group.Go(func() error {
		err := <-srv.ErrorChan
		return err
	})
	group.Go(func() error {
		<-gCtx.Done()
		return gCtx.Err()
	})

	if err := group.Wait(); err != nil {
		zlog.Fatal().Err(err).Msg("Error during server shutdown")
	} else {
		zlog.Info().Msg("Server excited gracefully")
	}
}

func initDB(addr string) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 7; i++ {
		time.Sleep(2 * time.Second)
		pool, err = pgxpool.New(context.Background(), addr)
		if err == nil {
			return pool, nil
		}
	}
	pool, err = pgxpool.New(context.Background(), addr)
	if err != nil {
		return nil, fmt.Errorf("database initialization error: %w", err)
	}
	return pool, nil
}
