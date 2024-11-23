package repository

import (
	"context"
	"errors"

	"github.com/ArataEM/message-service/model"
	"github.com/google/uuid"
)

type Repo interface {
	Insert(ctx context.Context, message model.Message) error
	Get(ctx context.Context, id uuid.UUID) (model.Message, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, message model.Message) error
	FindAll(ctx context.Context, page FindAllPage) (FindResult, error)
	Ping(ctx context.Context) error
	Close() error
}

type FindAllPage struct {
	Size   uint64
	Offset uint64
}

type FindResult struct {
	Messages []model.Message
	Cursor   uint64
}

var ErrNotExist = errors.New("message does not exist")
