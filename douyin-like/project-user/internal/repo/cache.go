package repo

import (
	"context"
	"time"
)

//repo 包是接口  Cache 存储数据库的接口

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}
