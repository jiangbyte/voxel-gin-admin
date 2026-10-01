// internal/modules/sys/config/service.go 业务服务。
//
// Author: Charlie

package config

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"voxel-gin-admin/internal/infrastructure/core/crypto"
	"voxel-gin-admin/internal/types/flagutil"
	"voxel-gin-admin/internal/infrastructure/platform/configsync"
	"voxel-gin-admin/internal/infrastructure/platform/idgen"
	"voxel-gin-admin/internal/infrastructure/platform/module"
	"voxel-gin-admin/internal/infrastructure/platform/notify"
	"voxel-gin-admin/internal/infrastructure/platform/runtimecfg"
	"voxel-gin-admin/internal/infrastructure/platform/storage"
)

// 敏感配置键（与 FastAPI SENSITIVE_CONFIG_KEYS 对齐的子集）。
var sensitiveKeys = map[string]struct{}{
	"AUTH_DEFAULT_PASSWORD":         {},
	"AUDIT_ALERT_WEBHOOK_SECRET":    {},
	"MAIL_LOCAL_PASSWORD":           {},
	"MAIL_ALIYUN_ACCESS_KEY_SECRET": {},
	"MAIL_TENCENT_SECRET_KEY":       {},
	"SMS_ALIYUN_ACCESS_KEY_SECRET":  {},
	"SMS_TENCENT_SECRET_KEY":        {},
	"PUSH_DINGTALK_SECRET":          {},
	"PUSH_LARK_SECRET":              {},
	"STORAGE_MINIO_ACCESS_KEY":      {},
	"STORAGE_MINIO_SECRET_KEY":      {},
	"STORAGE_RUSTFS_ACCESS_KEY":     {},
	"STORAGE_RUSTFS_SECRET_KEY":     {},
	"STORAGE_ALIYUN_ACCESS_KEY":     {},
	"STORAGE_ALIYUN_SECRET_KEY":     {},
	"STORAGE_TENCENT_ACCESS_KEY":    {},
	"STORAGE_TENCENT_SECRET_KEY":    {},
}

// Service 系统配置业务服务。
//
// Author: Charlie
type Service struct {
	repo    *Repo
	fernet  *crypto.Codec
	notify  *notify.Facade
	storage *storage.Manager
	runtime *runtimecfg.Settings
	redis   *redis.Client
}

// NewService 构造配置服务。
func NewService(db *gorm.DB, fernet *crypto.Codec, nf *notify.Facade, sto *storage.Manager, rt *runtimecfg.Settings, rdb *redis.Client) *Service {
	return &Service{repo: NewRepo(db), fernet: fernet, notify: nf, storage: sto, runtime: rt, redis: rdb}
}

// New 构建 sys.config 模块。

// ServiceFromDeps 从模块依赖构造领域 Service（供 Case/Trigger 组装）。
func ServiceFromDeps(d *module.Deps) *Service {
	var codec *crypto.Codec
	if d.Cfg != nil {
		codec, _ = crypto.NewFernetFromConfig(d.Cfg.Crypto.FernetKey, d.Cfg.Crypto.VaultAddr)
	}
	return NewService(d.DB, codec, d.Notify, d.Storage, d.Runtime, d.Redis)
}

func New(d *module.Deps) module.Module {
	_ = d
	return module.Module{
		Name:   "sys.config",
		Models: []any{&Config{}},
	}
}

func pickRemark(description, remark *string) *string {
	if remark != nil {
		return remark
	}
	return description
}

func (s *Service) maybeEncrypt(key string, category *string, value *string) *string {
	if value == nil || *value == "" || s.fernet == nil {
		return value
	}
	if crypto.LooksEncrypted(*value) {
		return value
	}
	if !isSensitive(key, category) {
		return value
	}
	enc, err := s.fernet.Encrypt(*value)
	if err != nil {
		return value
	}
	return &enc
}

func (s *Service) maybeDecrypt(key string, category *string, value *string) *string {
	if value == nil || *value == "" || s.fernet == nil {
		return value
	}
	if !crypto.LooksEncrypted(*value) && !isSensitive(key, category) {
		return value
	}
	plain, err := s.fernet.Decrypt(*value)
	if err != nil {
		return value
	}
	return &plain
}

func isSensitive(key string, category *string) bool {
	if _, ok := sensitiveKeys[key]; ok {
		return true
	}
	// AUTH_OAUTH_*_CLIENT_SECRET / AUTH_OAUTH_*_APP_SECRET（对齐 voxel-boot isSensitive）。
	up := strings.ToUpper(strings.TrimSpace(key))
	if strings.HasPrefix(up, "AUTH_OAUTH_") && (strings.HasSuffix(up, "_CLIENT_SECRET") || strings.HasSuffix(up, "_APP_SECRET")) {
		return true
	}
	if category == nil {
		return false
	}
	c := strings.ToUpper(strings.TrimSpace(*category))
	return c == "SECRET" || c == "SENSITIVE" || c == "CREDENTIAL"
}

func (s *Service) reveal(row *Config) {
	if row == nil {
		return
	}
	// 1. 合并 config_json / config_value 为对外统一字段
	row.ConfigValue = ExternalValue(row)
	plain := s.maybeDecrypt(row.ConfigKey, row.Category, row.ConfigValue)
	// 敏感键对外不回显明文，仅置 ext_json.is_set（对齐 voxel-boot toPublicResult）。
	if isSensitive(row.ConfigKey, row.Category) {
		if plain != nil && strings.TrimSpace(*plain) != "" {
			var ext map[string]any
			if len(row.ExtJSON) > 0 {
				_ = json.Unmarshal(row.ExtJSON, &ext)
			}
			if ext == nil {
				ext = map[string]any{}
			}
			ext["is_set"] = true
			b, _ := json.Marshal(ext)
			row.ExtJSON = b
		}
		empty := ""
		row.ConfigValue = &empty
		return
	}
	row.ConfigValue = plain
}

// Create 创建配置（config_key 唯一校验；对齐 voxel-boot ConfigServiceImpl.create）。
func (s *Service) Create(ctx context.Context, req AddParam) error {
	if _, err := s.repo.GetByKey(ctx, req.ConfigKey); err == nil {
		return errors.New("配置键已存在")
	}
	// 1. 加密敏感标量并按 value_type 分列落库
	val := s.maybeEncrypt(req.ConfigKey, req.Category, req.ConfigValue)
	scalar, jsonCol, vt := ApplyExternalValue(req.ValueType, val)
	row := Config{
		ID: idgen.Next(), ConfigKey: req.ConfigKey, ConfigValue: scalar, ConfigJSON: jsonCol,
		Category: req.Category,
		Remark: req.Remark, SortCode: req.SortCode, ValueType: vt, Label: req.Label, Scope: req.Scope,
		Scene: req.Scene, IsBuiltin: 0, ExtJSON: datatypes.JSON([]byte("{}")),
	}
	if err := s.repo.Create(ctx, &row); err != nil {
		return err
	}
	s.refreshRuntime(ctx, "config.create")
	return nil
}

// Update 更新配置（config_key 唯一排除自身 + 内置保护；对齐 voxel-boot ConfigServiceImpl.update）。
func (s *Service) Update(ctx context.Context, req EditParam) error {
	cur, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errors.New("配置不存在")
	}
	if existing, err2 := s.repo.GetByKey(ctx, req.ConfigKey); err2 == nil && existing.ID != req.ID {
		return errors.New("配置键已存在")
	}
	val := s.maybeEncrypt(req.ConfigKey, req.Category, req.ConfigValue)
	scalar, jsonCol, vt := ApplyExternalValue(req.ValueType, val)
	updates := map[string]any{
		"config_key": req.ConfigKey, "config_value": scalar, "config_json": jsonCol, "category": req.Category,
		"remark": req.Remark, "sort_code": req.SortCode, "value_type": vt, "label": req.Label,
		"scope": req.Scope, "scene": req.Scene,
	}
	if cur.IsBuiltin == 1 {
		// 内置配置不可改 scene；is_builtin 恒为 1（对齐 voxel-boot）
		delete(updates, "scene")
		updates["is_builtin"] = 1
		if req.Scope == nil || *req.Scope == "" {
			updates["scope"] = cur.Scope
		}
	}
	if err := s.repo.UpdateOptimistic(ctx, req.ID, req.Version, updates); err != nil {
		return err
	}
	s.refreshRuntime(ctx, "config.update")
	return nil
}

// Delete 批量删除（内置配置禁止删除；对齐 voxel-boot ConfigServiceImpl.delete）。
func (s *Service) Delete(ctx context.Context, ids []string) error {
	rows, err := s.repo.ListByIDs(ctx, ids)
	if err != nil {
		return err
	}
	var builtin []string
	for i := range rows {
		if rows[i].IsBuiltin == 1 {
			builtin = append(builtin, rows[i].ConfigKey)
		}
	}
	if len(builtin) > 0 {
		return errors.New("内置配置不可删除: " + strings.Join(builtin, ", "))
	}
	if err := s.repo.DeleteByIDs(ctx, ids); err != nil {
		return err
	}
	s.refreshRuntime(ctx, "config.delete")
	return nil
}

// Detail 详情（敏感值解密后返回）。
func (s *Service) Detail(ctx context.Context, id string) (*Config, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	s.reveal(row)
	return row, nil
}

// Page 分页。
func (s *Service) Page(ctx context.Context, q PageParam) (rows []Config, total int64, current, size int, err error) {
	current, size = q.Normalize()
	rows, total, err = s.repo.Page(ctx, q)
	for i := range rows {
		s.reveal(&rows[i])
	}
	return rows, total, current, size, err
}

// List 列表。
func (s *Service) List(ctx context.Context, q ListParam) ([]Config, error) {
	rows, err := s.repo.List(ctx, q)
	for i := range rows {
		s.reveal(&rows[i])
	}
	return rows, err
}

// BatchSave 批量保存配置：按 config_key 存在则更新，不存在则新建。
func (s *Service) BatchSave(ctx context.Context, req BatchSaveParam) error {
	items := make([]BatchItemParam, 0, len(req.Items))
	for _, it := range req.Items {
		if isSensitive(it.ConfigKey, it.Category) && (it.ConfigValue == nil || *it.ConfigValue == "") {
			// 敏感配置空值跳过，避免覆盖已有密文
			continue
		}
		items = append(items, it)
	}
	if len(items) == 0 {
		return nil
	}
	keys := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, it := range items {
		if _, ok := seen[it.ConfigKey]; ok {
			continue
		}
		seen[it.ConfigKey] = struct{}{}
		keys = append(keys, it.ConfigKey)
	}
	rows, err := s.repo.ListByKeys(ctx, keys)
	if err != nil {
		return err
	}
	existing := make(map[string]Config, len(rows))
	for _, row := range rows {
		existing[row.ConfigKey] = row
	}
	creates := make([]Config, 0)
	for _, it := range items {
		row, ok := existing[it.ConfigKey]
		// 敏感键落库后以 ext_json.is_set=true 标记「已配置」（对齐 web loadByCategory 契约）。
		ext := datatypes.JSON([]byte("{}"))
		if isSensitive(it.ConfigKey, it.Category) && it.ConfigValue != nil && *it.ConfigValue != "" {
			ext, _ = json.Marshal(map[string]bool{"is_set": true})
		}
		enc := s.maybeEncrypt(it.ConfigKey, it.Category, it.ConfigValue)
		vtInput := strOr(it.ValueType, "STRING")
		if ok && it.ValueType == nil {
			vtInput = row.ValueType
		}
		scalar, jsonCol, vt := ApplyExternalValue(vtInput, enc)
		if !ok {
			creates = append(creates, Config{
				ID:          idgen.Next(),
				ConfigKey:   it.ConfigKey,
				ConfigValue: scalar,
				ConfigJSON:  jsonCol,
				Category:    it.Category,
				Remark:      pickRemark(it.Description, it.Remark),
				ValueType:   vt,
				Label:       it.Label,
				Scope:       it.Scope,
				Scene:       it.Scene,
				IsBuiltin:   flagutil.OrInt(it.IsBuiltin),
				SortCode:    intOr(it.SortCode),
				ExtJSON:     ext,
			})
			continue
		}
		updates := map[string]any{
			"config_value": scalar,
			"config_json":  jsonCol,
			"value_type":   vt,
		}
		if string(ext) != "{}" {
			updates["ext_json"] = ext
		}
		if it.Category != nil {
			updates["category"] = it.Category
		}
		if it.Description != nil {
			updates["remark"] = pickRemark(it.Description, it.Remark)
		}
		if it.Label != nil {
			updates["label"] = it.Label
		}
		if it.Scope != nil {
			updates["scope"] = it.Scope
		}
		if it.Scene != nil {
			updates["scene"] = it.Scene
		}
		if it.IsBuiltin != nil {
			updates["is_builtin"] = *it.IsBuiltin
		}
		if it.SortCode != nil {
			updates["sort_code"] = *it.SortCode
		}
		if err := s.repo.UpdateOptimistic(ctx, row.ID, it.Version, updates); err != nil {
			return err
		}
	}
	if err := s.repo.CreateInBatches(ctx, creates); err != nil {
		return err
	}
	s.refreshRuntime(ctx, "config.batch_save")
	return nil
}

func (s *Service) refreshRuntime(ctx context.Context, reason string) {
	if s.runtime != nil {
		s.runtime.Invalidate()
	}
	if s.storage != nil {
		s.storage.Refresh()
	}
	configsync.Publish(ctx, s.redis, reason)
}

func strOr(v *string, def string) string {
	if v == nil || *v == "" {
		return def
	}
	return *v
}

func boolOr(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func intOr(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

// TestWebhook 发送审计告警测试 Webhook。
func (s *Service) TestWebhook(ctx context.Context, url, secret string) error {
	if s.notify == nil {
		return errors.New("notify facade unavailable")
	}
	payload := `{"msg_type":"text","content":{"text":"审计告警系统测试消息"}}`
	return s.notify.SendWebhook(ctx, url, secret, payload)
}

// TestPush 发送审计告警测试推送。
func (s *Service) TestPush(ctx context.Context) error {
	if s.notify == nil {
		return errors.New("notify facade unavailable")
	}
	return s.notify.SendPush(ctx, "审计告警", "审计告警系统测试消息")
}
