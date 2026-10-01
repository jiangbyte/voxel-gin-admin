// Package trigger 是 xfg-ddd 触发器层（HTTP/MQ/Job 入口）。
// HTTP 适配位于 trigger/http/...，只做 bind/鉴权转发并调用 Case。
package trigger
