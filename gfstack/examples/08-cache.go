// ================================================================
// 示例: Cache 缓存适配器 (internal/library/cache/cache.go)
// 支持 redis/file/memory 三种缓存驱动
// ================================================================

package cache

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gfile"

	"hotgo/internal/library/cache/file"
)

var cache *gcache.Cache

// Instance 缓存实例
func Instance() *gcache.Cache {
	if cache == nil {
		panic("cache uninitialized.")
	}
	return cache
}

// SetAdapter 设置缓存适配器，根据 config.yaml 中 cache.adapter 选择驱动
func SetAdapter(ctx context.Context) {
	var adapter gcache.Adapter

	switch g.Cfg().MustGet(ctx, "cache.adapter").String() {
	case "redis":
		adapter = gcache.NewAdapterRedis(g.Redis())
	case "file":
		fileDir := g.Cfg().MustGet(ctx, "cache.fileDir").String()
		if fileDir == "" {
			g.Log().Fatal(ctx, "file path must be configured for file caching.")
			return
		}
		if !gfile.Exists(fileDir) {
			if err := gfile.Mkdir(fileDir); err != nil {
				g.Log().Fatalf(ctx, "failed to create the cache directory. err:%+v", err)
				return
			}
		}
		adapter = file.NewAdapterFile(fileDir)
	default:
		adapter = gcache.NewAdapterMemory()
	}

	// 同时设置数据库缓存适配器
	g.DB().GetCache().SetAdapter(adapter)

	// 通用缓存
	cache = gcache.New()
	cache.SetAdapter(adapter)
}
