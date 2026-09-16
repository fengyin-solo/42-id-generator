package model

import (
	"strings"
	"time"
)

const (
	NodeStatusActive   = "active"
	NodeStatusInactive = "inactive"
	NodeStatusFailed   = "failed"
)

var nodeTransitions = map[string]map[string]bool{
	NodeStatusActive:   {NodeStatusInactive: true, NodeStatusFailed: true},
	NodeStatusInactive: {NodeStatusActive: true, NodeStatusFailed: true},
	NodeStatusFailed:   {NodeStatusInactive: true},
}

// CanTransitionNode 判断节点状态是否允许流转。
func CanTransitionNode(from, to string) bool {
	if m, ok := nodeTransitions[from]; ok {
		return m[to]
	}
	return false
}

// MachineNode 机器节点（worker），用于雪花算法和号段分配。
type MachineNode struct {
	ID           string    `json:"id"`
	NodeID       string    `json:"node_id"`
	Hostname     string    `json:"hostname"`
	IP           string    `json:"ip"`
	WorkerID     int64     `json:"worker_id"`
	DatacenterID int64     `json:"datacenter_id"`
	Status       string    `json:"status"`
	LastBeatAt   time.Time `json:"last_beat_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (n *MachineNode) Validate() error {
	n.Hostname = strings.TrimSpace(n.Hostname)
	n.IP = strings.TrimSpace(n.IP)
	n.NodeID = strings.TrimSpace(n.NodeID)
	if n.NodeID == "" {
		return NewValidationError("node_id", "节点标识不能为空")
	}
	if n.Hostname == "" {
		return NewValidationError("hostname", "主机名不能为空")
	}
	if n.IP == "" {
		return NewValidationError("ip", "IP 不能为空")
	}
	if n.WorkerID < 0 {
		return NewValidationError("worker_id", "worker_id 不能为负数")
	}
	if n.DatacenterID < 0 {
		return NewValidationError("datacenter_id", "datacenter_id 不能为负数")
	}
	if n.Status == "" {
		n.Status = NodeStatusActive
	}
	if n.Status != NodeStatusActive && n.Status != NodeStatusInactive && n.Status != NodeStatusFailed {
		return NewValidationError("status", "节点状态不合法")
	}
	return nil
}

type MachineNodeFilter struct {
	Status     string
	DatacenterID *int64
	Keyword    string
}

func (f MachineNodeFilter) Match(n *MachineNode) bool {
	if f.Status != "" && n.Status != f.Status {
		return false
	}
	if f.DatacenterID != nil && n.DatacenterID != *f.DatacenterID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(n.Hostname), k) &&
			!strings.Contains(strings.ToLower(n.NodeID), k) {
			return false
		}
	}
	return true
}
