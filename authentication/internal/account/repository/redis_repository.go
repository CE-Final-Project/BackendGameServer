package repository

import (
	"context"
	"encoding/json"
	"github.com/ce-final-project/backend_game_server/authentication/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const (
	redisAccountPrefixKey = "auth:account"
)

type redisRepository struct {
	redisClient redis.UniversalClient
	log         logger.Logger
}

func (r *redisRepository) getRedisAccountPrefixKey() string {
	return redisAccountPrefixKey
}

func (r *redisRepository) PutAccount(ctx context.Context, key string, account *models.Account) {
	accountBytes, err := json.Marshal(account)
	if err != nil {
		r.log.WarnMsg("json.Marshal", err)
		return
	}

	if err := r.redisClient.HSetNX(ctx, r.getRedisAccountPrefixKey(), key, accountBytes).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisAccountPrefixKey(), key)
}

func (r *redisRepository) PutKeyReference(ctx context.Context, key string, targetKey string) {
	if err := r.redisClient.HSetNX(ctx, r.getRedisAccountPrefixKey(), key, []byte(targetKey)).Err(); err != nil {
		r.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	r.log.Debugf("HSetNX prefix: %s, key: %s", r.getRedisAccountPrefixKey(), key)
}

func (r *redisRepository) GetAccount(ctx context.Context, key string) (*models.Account, error) {
	accountBytes, err := r.redisClient.HGet(ctx, r.getRedisAccountPrefixKey(), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}

	var account models.Account
	if err := json.Unmarshal(accountBytes, &account); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisAccountPrefixKey(), key)
	return &account, nil
}

func (r *redisRepository) GetAccountByKeyReference(ctx context.Context, key string) (*models.Account, error) {
	accountIDBytes, err := r.redisClient.HGet(ctx, r.getRedisAccountPrefixKey(), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var accountBytes []byte
	accountBytes, err = r.redisClient.HGet(ctx, r.getRedisAccountPrefixKey(), string(accountIDBytes)).Bytes()
	if err != nil {
		if err != redis.Nil {
			r.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}
	var account models.Account
	if err := json.Unmarshal(accountBytes, &account); err != nil {
		return nil, err
	}
	r.log.Debugf("HGet prefix: %s, key: %s", r.getRedisAccountPrefixKey(), key)
	return &account, nil
}

func (r *redisRepository) DelAccount(ctx context.Context, key string) {
	if err := r.redisClient.HDel(ctx, r.getRedisAccountPrefixKey(), key).Err(); err != nil {
		r.log.WarnMsg("redisClient.HDel", err)
		return
	}
	r.log.Debugf("HDel prefix: %s, key: %s", r.getRedisAccountPrefixKey(), key)
}

func (r *redisRepository) DelAllAccounts(ctx context.Context) {
	if err := r.redisClient.Del(ctx, r.getRedisAccountPrefixKey()).Err(); err != nil {
		r.log.WarnMsg("redisClient.Del", err)
		return
	}
	r.log.Debugf("Del key: %s", r.getRedisAccountPrefixKey())
}

func NewCacheRepository(redisClient redis.UniversalClient, log logger.Logger) CacheRepository {
	return &redisRepository{redisClient: redisClient, log: log}
}
