package repository

import (
	"context"
	"encoding/json"
	"github.com/ce-final-project/backend_game_server/account/internal/models"
	"github.com/ce-final-project/backend_game_server/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const (
	redisAccountPrefixKey = "auth:account"
)

type accountCacheRepository struct {
	log         logger.Logger
	redisClient redis.UniversalClient
}

func NewAccountCacheRepository(log logger.Logger, redis redis.UniversalClient) *accountCacheRepository {
	return &accountCacheRepository{
		log:         log,
		redisClient: redis,
	}
}

func (a *accountCacheRepository) PutAccount(ctx context.Context, key string, account *models.Account) {
	accountBytes, err := json.Marshal(account)
	if err != nil {
		a.log.WarnMsg("json.Marshal", err)
		return
	}

	if err := a.redisClient.HSetNX(ctx, a.getRedisAccountPrefixKey(), key, accountBytes).Err(); err != nil {
		a.log.WarnMsg("redisClient.HSetNX", err)
		return
	}
	a.log.Debugf("HSetNX prefix: %s, key: %s", a.getRedisAccountPrefixKey(), key)
}

func (a *accountCacheRepository) GetAccount(ctx context.Context, key string) (*models.Account, error) {
	accountBytes, err := a.redisClient.HGet(ctx, a.getRedisAccountPrefixKey(), key).Bytes()
	if err != nil {
		if err != redis.Nil {
			a.log.WarnMsg("redisClient.HGet", err)
		}
		return nil, errors.Wrap(err, "redisClient.HGet")
	}

	var account models.Account
	if err := json.Unmarshal(accountBytes, &account); err != nil {
		return nil, err
	}
	a.log.Debugf("HGet prefix: %s, key: %s", a.getRedisAccountPrefixKey(), key)
	return &account, nil
}

func (a *accountCacheRepository) DelAccount(ctx context.Context, key string) {
	if err := a.redisClient.HDel(ctx, a.getRedisAccountPrefixKey(), key).Err(); err != nil {
		a.log.WarnMsg("redisClient.HDel", err)
		return
	}
	a.log.Debugf("HDel prefix: %s, key: %s", a.getRedisAccountPrefixKey(), key)
}

func (a *accountCacheRepository) DelAllAccounts(ctx context.Context) {
	if err := a.redisClient.Del(ctx, a.getRedisAccountPrefixKey()).Err(); err != nil {
		a.log.WarnMsg("redisClient.Del", err)
		return
	}
	a.log.Debugf("Del key: %s", a.getRedisAccountPrefixKey())
}

func (a *accountCacheRepository) getRedisAccountPrefixKey() string {
	return redisAccountPrefixKey
}
