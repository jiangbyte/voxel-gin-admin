// Package codegen 标志位 0/1 辅助（对齐 abolish-boolean）。
//
// Author: Charlie
package codegen

// flagOn 判断持久化标志是否为 1。
func flagOn(v int) bool { return v == 1 }

// boolFlag 将布尔语义转为 0/1 存库。
func boolFlag(v bool) int {
	if v {
		return 1
	}
	return 0
}
