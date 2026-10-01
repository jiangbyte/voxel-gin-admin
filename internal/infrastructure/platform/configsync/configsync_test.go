// Package configsync 测试。
//
// Author: Charlie
package configsync

import (
	"context"
	"encoding/json"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"voxel-gin-admin/internal/infrastructure/platform/cache"
)

func TestPublishAndReceive(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer rdb.Close()

	var reloads atomic.Int32
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stop := Start(ctx, rdb, func() { reloads.Add(1) })
	defer stop()

	// 等待订阅就绪
	time.Sleep(50 * time.Millisecond)

	// 本实例 Publish 应被忽略
	if !Publish(ctx, rdb, "self") {
		t.Fatal("Publish self should succeed")
	}
	time.Sleep(50 * time.Millisecond)
	if reloads.Load() != 0 {
		t.Fatalf("self event must be ignored, reloads=%d", reloads.Load())
	}

	// 模拟其他实例
	payload, _ := json.Marshal(map[string]string{
		"source": "other-instance",
		"reason": "config.update",
		"at":     time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err := rdb.Publish(ctx, cache.ConfigChangedChannel, payload).Err(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if reloads.Load() >= 1 {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if reloads.Load() < 1 {
		t.Fatal("expected reload from peer event")
	}
}

func TestDecodeEvent(t *testing.T) {
	ev, ok := decodeEvent(`{"source":"a","reason":"x","at":"t"}`)
	if !ok || ev.Source != "a" || ev.Reason != "x" {
		t.Fatalf("got %#v ok=%v", ev, ok)
	}
	if _, ok := decodeEvent(`{}`); ok {
		t.Fatal("empty source should fail")
	}
	if _, ok := decodeEvent(`not-json`); ok {
		t.Fatal("invalid json should fail")
	}
}

func TestPublishNilRedis(t *testing.T) {
	if Publish(context.Background(), nil, "x") {
		t.Fatal("nil redis should return false")
	}
}
