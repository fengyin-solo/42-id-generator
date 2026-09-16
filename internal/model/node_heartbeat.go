package model

import (
	"strings"
	"time"
)

// NodeHeartbeat 心跳记录，用于节点健康监控。
type NodeHeartbeat struct {
	ID           string    `json:"id"`
	NodeID       string    `json:"node_id"`
	Load         float64   `json:"load"`
	CPUUsage     float64   `json:"cpu_usage"`
	MemoryUsage  float64   `json:"memory_usage"`
	BeatAt       time.Time `json:"beat_at"`
	CreatedAt    time.Time `json:"created_at"`
}

func (h *NodeHeartbeat) Validate() error {
	h.NodeID = strings.TrimSpace(h.NodeID)
	if h.NodeID == "" {
		return NewValidationError("node_id", "节点不能为空")
	}
	if h.Load < 0 || h.Load > 100 {
		return NewValidationError("load", "负载必须在 0-100 之间")
	}
	if h.CPUUsage < 0 || h.CPUUsage > 100 {
		return NewValidationError("cpu_usage", "CPU 使用率必须在 0-100 之间")
	}
	if h.MemoryUsage < 0 || h.MemoryUsage > 100 {
		return NewValidationError("memory_usage", "内存使用率必须在 0-100 之间")
	}
	if h.BeatAt.IsZero() {
		h.BeatAt = time.Now()
	}
	return nil
}

type NodeHeartbeatFilter struct {
	NodeID       string
	MinLoad      *float64
	MaxLoad      *float64
	MinCPUUsage  *float64
	MaxCPUUsage  *float64
	StartTime    *time.Time
	EndTime      *time.Time
}

func (f NodeHeartbeatFilter) Match(h *NodeHeartbeat) bool {
	if f.NodeID != "" && h.NodeID != f.NodeID {
		return false
	}
	if f.MinLoad != nil && h.Load < *f.MinLoad {
		return false
	}
	if f.MaxLoad != nil && h.Load > *f.MaxLoad {
		return false
	}
	if f.MinCPUUsage != nil && h.CPUUsage < *f.MinCPUUsage {
		return false
	}
	if f.MaxCPUUsage != nil && h.CPUUsage > *f.MaxCPUUsage {
		return false
	}
	if f.StartTime != nil && h.BeatAt.Before(*f.StartTime) {
		return false
	}
	if f.EndTime != nil && h.BeatAt.After(*f.EndTime) {
		return false
	}
	return true
}
