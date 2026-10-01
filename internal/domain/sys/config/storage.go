// Package config 列存储拆分：HTTP 统一 config_value，落库按 value_type 写入列。
//
// Author: Charlie
package config

import (
	"encoding/json"
	"errors"
	"strings"

	"gorm.io/datatypes"
)

// ErrVersionConflict 乐观锁冲突（对外可识别）。
var ErrVersionConflict = errors.New("配置已被他人修改，请刷新后重试")

// NormalizeValueType 归一化 value_type（INT/INTEGER/LONG/BOOL → NUMBER）。
func NormalizeValueType(valueType string) string {
	u := strings.ToUpper(strings.TrimSpace(valueType))
	switch u {
	case "INT", "INTEGER", "LONG", "BOOL":
		return "NUMBER"
	case "":
		return "STRING"
	default:
		return u
	}
}

// ExternalValue 将 DB 列合并为对外统一的 config_value。
func ExternalValue(row *Config) *string {
	if row == nil {
		return nil
	}
	vt := NormalizeValueType(row.ValueType)
	if vt == "JSON" {
		if len(row.ConfigJSON) > 0 && string(row.ConfigJSON) != "null" {
			s := string(row.ConfigJSON)
			return &s
		}
		return row.ConfigValue
	}
	return row.ConfigValue
}

// ApplyExternalValue 按 value_type 将统一 config_value 拆分到列字段。
func ApplyExternalValue(valueType string, external *string) (scalar *string, jsonCol datatypes.JSON, normalized string) {
	// 1. 归一化类型
	normalized = NormalizeValueType(valueType)
	// 2. JSON 写入 config_json；标量写入 config_value
	if normalized == "JSON" {
		if external == nil || strings.TrimSpace(*external) == "" {
			return nil, nil, normalized
		}
		raw := []byte(strings.TrimSpace(*external))
		if !json.Valid(raw) {
			return nil, datatypes.JSON(raw), normalized
		}
		return nil, datatypes.JSON(raw), normalized
	}
	return external, nil, normalized
}
