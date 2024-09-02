package server

import (
	"context"

	"github.com/go-playground/validator"
	"github.com/lahnasti/go-market/internal/repository"
	"github.com/rs/zerolog"
)

type Server struct {
	Db repository.Repository
	ErrorChan  chan error
	deleteChan chan int
	Valid      *validator.Validate
	log        zerolog.Logger
}

func NewServer(ctx context.Context, db repository.Repository, zlog *zerolog.Logger) *Server {
	dChan := make(chan int, 5)
	errChan := make(chan error)
	srv := Server{
		Db:         db,
		deleteChan: dChan,
		ErrorChan:  errChan,
		log:        *zlog,
	}
	go srv.deleter(ctx)
	return &Server{
		Db:         db,
		deleteChan: dChan,
		ErrorChan:  errChan,
	}
}
