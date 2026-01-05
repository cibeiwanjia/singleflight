package main

import (
	"context"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
)

var (
	rdb     = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	sfGroup singleflight.Group
	ctx     = context.Background()
)

func getHotData(key string) string {
	if v, _ := rdb.Get(ctx, key).Result(); v != "" {
		return v
	}
	val, _, _ := sfGroup.Do(key, func() (interface{}, error) {
		data := "永辉超市: 牛肉 | 价格49.9"
		rdb.Set(ctx, key, data, 5*60)
		return data, nil
	})
	return val.(string)
}

func main() {
	r := gin.Default()
	// Gin接口+Header设置
	r.GET("/hot", func(c *gin.Context) {
		// 强制缓存Header设置（核心）
		maxAge := 300 // 缓存5分钟（和Redis缓存时间一致）
		c.Header("Cache-Control", "public, max-age="+string(rune(maxAge)))
		c.Header("Expires", time.Now().Add(time.Second*time.Duration(maxAge)).Format(time.RFC1123)) // 兼容HTTP/1.0
		c.Header("Access-Control-Allow-Origin", "*")

		// 返回数据
		c.String(http.StatusOK, getHotData("goods:660390230106"))
	})
	r.Run(":8080")
}
