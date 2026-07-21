package api_redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"repo-stat/api/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	redisDB *redis.Client
	log     *slog.Logger
	ttl     time.Duration
	ch      chan *domain.RepoInfo
}

func NewRedisRepo(address string, log *slog.Logger, ttlSec int) *RedisRepo {
	rRepo := &RedisRepo{
		redisDB: redis.NewClient(&redis.Options{
			Addr: address,
		}),
		log: log,
		ttl: time.Duration(ttlSec) * time.Second,
		ch:  make(chan *domain.RepoInfo, 100),
	}

	for i := 1; i <= 5; i++ {
		go rRepo.CacheWork()
	}

	return rRepo
}

func (rr *RedisRepo) CacheWork() {

	for task := range rr.ch {
		_ = rr.SetRepoInfo(context.Background(), task)
	}
}

func (rr *RedisRepo) AddToChan(repo *domain.RepoInfo) {
	select {
	case rr.ch <- repo:
	default:
		rr.log.Error("cache queue is overloaded")
	}
}

func (rr *RedisRepo) GetRepoInfo(ctx context.Context, req *domain.GetRepoInfoReq) (*domain.RepoInfo, error) {

	RepoInfoString, err := rr.redisDB.Get(ctx, strings.ToLower(req.Owner+"/"+req.Repo)).Result()

	if err != nil {

		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		rr.log.Error("failed to get repoinfo from cache", "error", err)
		return nil, domain.ErrInternalError
	}

	var RepoInfo domain.RepoInfo
	err = json.Unmarshal([]byte(RepoInfoString), &RepoInfo)

	if err != nil {
		rr.log.Error("failed to parse string", "error", err)
		return nil, nil
	}

	return &RepoInfo, nil

}

func (rr *RedisRepo) SetRepoInfo(ctx context.Context, repoInfo *domain.RepoInfo) error {

	repoInfoSlice, err := json.Marshal(repoInfo)

	if err != nil {
		rr.log.Error("failed to serialize repoinfo", "error", err)
		return domain.ErrInternalError
	}

	err = rr.redisDB.Set(ctx, strings.ToLower(repoInfo.FullName), repoInfoSlice, rr.ttl).Err()

	if err != nil {
		rr.log.Error("failed to set repoinfo in cache", "error", err)
		return domain.ErrInternalError
	}

	return nil

}

func (rr *RedisRepo) CheckRateLimit(ctx context.Context, ip string, limit int, rps float64) (bool, error) {

	key := fmt.Sprintf("limiter:%s", ip)
	now := time.Now().UnixNano()
	windowStart := now - int64(limit/int(rps))*int64(time.Second)

	_, err := rr.redisDB.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart)).Result()

	if err != nil {
		return false, domain.ErrInternalError
	}

	count, err := rr.redisDB.ZCard(ctx, key).Result()
	if err != nil {
		return false, domain.ErrInternalError
	}

	if count >= int64(limit) {
		return false, nil
	}

	err = rr.redisDB.ZAdd(ctx, key, redis.Z{
		Score:  float64(now),
		Member: fmt.Sprintf("%d-%s", now, uuid.New().String()),
	}).Err()

	if err != nil {
		return false, domain.ErrInternalError
	}

	rr.redisDB.Expire(ctx, key, time.Second)

	return true, nil
}
