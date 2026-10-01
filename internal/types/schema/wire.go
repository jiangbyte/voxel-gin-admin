// Package schema wire helpers — JSON 数值一律字符串（对齐 voxel-boot StringlyTypedJacksonModule）。
//
// Author: Charlie
package schema

import (
	"encoding/json"
	"strconv"
	"strings"
)

// WireInt JSON 线型整数：出站为十进制字符串（"0"/"1"/"300" 同一规则）。
type WireInt int

// MarshalJSON implements json.Marshaler.
func (v WireInt) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.Itoa(int(v)))), nil
}

// WireIntPtr converts *int to *WireInt.
func WireIntPtr(v *int) *WireInt {
	if v == nil {
		return nil
	}
	w := WireInt(*v)
	return &w
}

// WireIntValue converts int to WireInt.
func WireIntValue(v int) WireInt { return WireInt(v) }

// WireFlag 业务 0/1 标志的类型别名；线型与 WireInt 完全相同，仅语义标注。
type WireFlag = WireInt

// WireFlagValue converts int flag to WireFlag.
func WireFlagValue(v int) WireFlag { return WireIntValue(v) }

// IntStringPtr formats *int as decimal string pointer (nil stays nil).
func IntStringPtr(v *int) *string {
	if v == nil {
		return nil
	}
	s := strconv.Itoa(*v)
	return &s
}

// IntString formats int as decimal string.
func IntString(v int) string { return strconv.Itoa(v) }

// StringPtr returns nil for blank strings (boot nullable String fields).
func StringPtr(v *string) *string {
	if v == nil {
		return nil
	}
	if strings.TrimSpace(*v) == "" {
		return nil
	}
	return v
}

// JSONOrNull marshals empty JSON as null (align voxel-boot Map/JSON null).
type JSONOrNull json.RawMessage

// MarshalJSON implements json.Marshaler.
func (j JSONOrNull) MarshalJSON() ([]byte, error) {
	raw := []byte(j)
	if len(raw) == 0 || string(raw) == "null" || string(raw) == "{}" || string(raw) == "[]" {
		return []byte("null"), nil
	}
	return json.RawMessage(raw).MarshalJSON()
}

// JSONOrNullFromBytes builds JSONOrNull from DB bytes.
func JSONOrNullFromBytes(raw []byte) JSONOrNull {
	if len(raw) == 0 || string(raw) == "null" {
		return JSONOrNull(nil)
	}
	return JSONOrNull(raw)
}
