package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ArataEM/message-service/config"
	"github.com/ArataEM/message-service/model"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	Client *redis.Client
}

const messagesSetName string = "messages"

func messageIdKey(id uuid.UUID) string {
	return fmt.Sprintf("message:%s", id)
}

func NewRedisRepo(config config.Config) *RedisRepo {
	return &RedisRepo{
		Client: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		})}
}

func (r *RedisRepo) Insert(ctx context.Context, message model.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}
	key := messageIdKey(message.Id)

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("error creating message %s: %w", message.Id, err)
	}

	err = txn.SAdd(ctx, messagesSetName, key).Err()
	if err != nil {
		txn.Discard()
		return fmt.Errorf("error adding message %s to set: %w", message.Id, err)
	}

	_, err = txn.Exec(ctx)
	if err != nil {
		return fmt.Errorf("exec error for message %s: %w", message.Id, err)
	}

	return nil
}

func (r *RedisRepo) Get(ctx context.Context, id uuid.UUID) (model.Message, error) {
	key := messageIdKey(id)

	value, err := r.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Message{}, ErrNotExist
	} else if err != nil {
		return model.Message{}, fmt.Errorf("error getting message %s: %w", id, err)
	}

	var message model.Message
	err = json.Unmarshal([]byte(value), &message)
	if err != nil {
		return model.Message{}, fmt.Errorf("error decoding message %s: %w", id, err)
	}

	return message, nil
}

func (r *RedisRepo) Delete(ctx context.Context, id uuid.UUID) error {
	key := messageIdKey(id)

	txn := r.Client.TxPipeline()

	err := txn.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("error deleting message %s: %w", id, err)
	}

	err = txn.SRem(ctx, messagesSetName, key).Err()
	if err != nil {
		txn.Discard()
		return fmt.Errorf("error deleting message %s from set: %w", id, err)
	}

	_, err = txn.Exec(ctx)
	if err != nil {
		return fmt.Errorf("exec error for message %s: %w", id, err)
	}

	return nil
}

func (r *RedisRepo) Update(ctx context.Context, message model.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to decode message: %w", err)
	}
	key := messageIdKey(message.Id)

	err = r.Client.SetXX(ctx, key, string(data), 0).Err()
	if errors.Is(err, redis.Nil) {
		return ErrNotExist
	} else if err != nil {
		return fmt.Errorf("error updating message %s: %w", message.Id, err)
	}

	return nil
}

func (r *RedisRepo) FindAll(ctx context.Context, page FindAllPage) (FindResult, error) {
	res := r.Client.SScan(ctx, messagesSetName, page.Offset, "*", int64(page.Size))

	keys, cursor, err := res.Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to list messages: %w", err)
	}

	if len(keys) == 0 {
		return FindResult{
			Messages: []model.Message{},
		}, nil
	}

	messagesData, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return FindResult{}, fmt.Errorf("failed to get list: %w", err)
	}

	messages := make([]model.Message, len(messagesData))
	for k, v := range messagesData {
		v := v.(string)
		var message model.Message

		err := json.Unmarshal([]byte(v), &message)
		if err != nil {
			return FindResult{}, fmt.Errorf("failed to decode message: %w", err)
		}
		messages[k] = message
	}

	return FindResult{
		Messages: messages,
		Cursor:   cursor,
	}, nil
}

func (r *RedisRepo) Ping(ctx context.Context) error {
	err := r.Client.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}
	return nil
}

func (r *RedisRepo) Close() error {
	err := r.Client.Close()
	if err != nil {
		return fmt.Errorf("failed to close redis: %w", err)
	}
	return nil
}
