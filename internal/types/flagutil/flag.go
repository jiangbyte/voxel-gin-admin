// Package flagutil 业务标志 0/1 与布尔语义的转换。
//
// Author: Charlie
package flagutil

// FromBool 将布尔语义转为 0/1。
func FromBool(v bool) int {
	if v {
		return 1
	}
	return 0
}

// IsOn 判断标志是否为 1。
func IsOn(v int) bool { return v == 1 }

// OrInt 解引用可选标志，nil 视为 0。
func OrInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// OrIntDefault 解引用可选标志，nil 时使用默认值。
func OrIntDefault(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}
