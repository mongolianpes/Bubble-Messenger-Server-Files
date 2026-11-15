package rdb

import (
	"context"
	"time"
)

const (
	prefixForTempPath = "temp_path_"
	TTLTempPath       = time.Minute * 10
)

func (s *RedisStorage) SetStoragePath(ctx context.Context, tempPath, storagePath string) error {
	err := s.rdb.Set(ctx, prefixForTempPath+tempPath, storagePath, TTLTempPath).Err()
	return err
}

func (s *RedisStorage) GetStoragePath(ctx context.Context, tempPath string) (string, error) {
	storagePath, err := s.rdb.Get(ctx, prefixForTempPath+tempPath).Result()
	if err != nil {
		return "", err
	}

	return storagePath, nil
}
