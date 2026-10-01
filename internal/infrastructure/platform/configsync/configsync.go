// Package configsync 通过 Redis Pub/Sub 同步 sys_config 变更（对齐 voxel-boot ConfigChangeNotifier / voxel-fastapi sync）。
//
// Author: Charlie
package configsync

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"voxel-gin-admin/internal/infrastructure/core/logger"
	"voxel-gin-admin/internal/infrastructure/platform/cache"
)

// 本进程实例 ID，用于忽略自身发布的事件回环。
var instanceID = uuid.NewString()

var errPubSubClosed = errors.New("config sync pubsub channel closed")

// InstanceID 返回本进程配置同步实例 ID。
func InstanceID() string { return instanceID }

type eventPayload struct {
	Source string `json:"source"`
	Reason string `json:"reason"`
	At     string `json:"at"`
}

// Publish 向其他实例广播配置变更；Redis 不可用时打日志并返回 false，不阻断本地刷新。
func Publish(ctx context.Context, rdb *redis.Client, reason string) bool {
	if rdb == nil {
		logger.L.Warn("Config changed locally but Redis is unavailable; peers were not notified")
		return false
	}
	if reason == "" {
		reason = "updated"
	}
	payload, err := json.Marshal(eventPayload{
		Source: instanceID,
		Reason: reason,
		At:     time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		logger.L.Warn("Failed to encode config change event", zap.Error(err))
		return false
	}
	if err := rdb.Publish(ctx, cache.ConfigChangedChannel, payload).Err(); err != nil {
		logger.L.Warn("Failed to publish config change event", zap.Error(err))
		return false
	}
	return true
}

// Start 订阅配置变更频道；收到非本机事件时调用 onReload。返回 cancel 用于停止监听。
func Start(parent context.Context, rdb *redis.Client, onReload func()) context.CancelFunc {
	ctx, cancel := context.WithCancel(parent)
	if rdb == nil {
		logger.L.Warn("Config sync listener not started because Redis is unavailable")
		return cancel
	}
	if onReload == nil {
		onReload = func() {}
	}
	go listen(ctx, rdb, onReload)
	return cancel
}

func listen(ctx context.Context, rdb *redis.Client, onReload func()) {
	for {
		if ctx.Err() != nil {
			return
		}
		if err := runSubscribe(ctx, rdb, onReload); err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.L.Warn("Config sync listener reconnecting after failure", zap.Error(err))
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
			}
		}
	}
}

func runSubscribe(ctx context.Context, rdb *redis.Client, onReload func()) error {
	pubsub := rdb.Subscribe(ctx, cache.ConfigChangedChannel)
	defer func() { _ = pubsub.Close() }()

	if _, err := pubsub.Receive(ctx); err != nil {
		return err
	}
	logger.L.Info("Config sync listener started", zap.String("channel", cache.ConfigChangedChannel))

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-ch:
			if !ok {
				return errPubSubClosed
			}
			if msg == nil {
				continue
			}
			event, ok := decodeEvent(msg.Payload)
			if !ok {
				continue
			}
			if event.Source == instanceID {
				continue
			}
			onReload()
			logger.L.Info(
				"Reloaded config from distributed event",
				zap.String("reason", event.Reason),
				zap.String("source", event.Source),
			)
		}
	}
}

func decodeEvent(raw string) (eventPayload, bool) {
	var event eventPayload
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		return eventPayload{}, false
	}
	if event.Source == "" {
		return eventPayload{}, false
	}
	return event, true
}
