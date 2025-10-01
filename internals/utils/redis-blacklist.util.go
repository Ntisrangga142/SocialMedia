package utils

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func IsBlacklisted(ctx context.Context, rdb *redis.Client, token string) (bool, error) {
	res, err := rdb.Exists(ctx, "Blacklist:"+token).Result()
	if err != nil {
		return false, err
	}
	return res == 1, nil
}
