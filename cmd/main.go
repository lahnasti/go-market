package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lahnasti/go-market/internal/config"
	"github.com/lahnasti/go-market/internal/logger"
	"github.com/lahnasti/go-market/internal/repository"
	"github.com/lahnasti/go-market/internal/server"
	"github.com/lahnasti/go-market/internal/server/routes"
	"github.com/lahnasti/go-market/internal/storage"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-c
		cancel()
	}()
	fmt.Println("Server starting")
	cfg := config.ReadConfig()
	fmt.Println(cfg)
	zlog := logger.SetupLogger(cfg.DebugFlag)
	zlog.Debug().Any("config", cfg).Msg("Check cfg value")
	err := repository.Migrations(cfg.DBAddr, cfg.MPath, zlog)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Init migrations failed")
	}
	pool, err := initDB(cfg.DBAddr)
	if err != nil {
		zlog.Fatal().Err(err).Msg("Connection DB failed")
	}
	defer pool.Close()

	dbStorage, err := storage.NewDB(pool)
	if err != nil {
		panic(err)
	}

	group, gCtx := errgroup.WithContext(ctx)
	srv := server.NewServer(gCtx, dbStorage, zlog)
	group.Go(func() error {
		r := gin.Default()
		routes.SetupRoutes(srv)
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
		panic(err)
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
