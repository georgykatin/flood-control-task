package floodcontroller

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type flood struct {
	client       *redis.Client
	RequestLimit int
	Timeout      time.Duration
}

func NewFloodConfig(client *redis.Client, RequestLimit int, Timeout time.Duration) *flood {
	return &flood{
		client:       client,
		RequestLimit: RequestLimit,
		Timeout:      Timeout,
	}
}
func (fl *flood) Check(ctx context.Context, userID int64) (bool, error) {

	key := strconv.FormatInt(userID, 10)

	rec, err := fl.client.Get(ctx, key).Int64()
	if err != nil {
		return false, err
	}

	num, err := fl.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	now := time.Now()
	if now.Sub(time.Unix(rec, 0)) >= fl.Timeout {
		return int(num) <= fl.RequestLimit, nil
	}
	return true, nil
}
