package cache

import "context"

//缓存驱动接口，所以缓存驱动都需要实现以下接口
type Driver interface {
    Get(ctx context.Context, key string) (interface{}, error)
    GetMulti(ctx context.Context, keys ... string) (map[string]interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl int) (bool, error)
    SetMulti(ctx context.Context, items map[string]interface{}, ttl int) (bool, error)
    Delete(ctx context.Context, key string) (bool, error)
    DeleteMulti(ctx context.Context, key ... string) (bool, error)
    Expire(ctx context.Context, key string, ttl int) (bool, error)
}
