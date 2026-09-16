// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"idgenerator/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// BizType
	CreateBizType(b *model.BizType) error
	GetBizType(id string) (*model.BizType, error)
	GetBizTypeByCode(code string) (*model.BizType, error)
	ListBizTypes() []*model.BizType
	UpdateBizType(b *model.BizType) error
	DeleteBizType(id string) error

	// MachineNode
	CreateMachineNode(n *model.MachineNode) error
	GetMachineNode(id string) (*model.MachineNode, error)
	GetMachineNodeByNodeID(nodeID string) (*model.MachineNode, error)
	ListMachineNodes() []*model.MachineNode
	UpdateMachineNode(n *model.MachineNode) error
	DeleteMachineNode(id string) error

	// IDRule
	CreateIDRule(r *model.IDRule) error
	GetIDRule(id string) (*model.IDRule, error)
	GetIDRuleByBizTypeID(bizTypeID string) (*model.IDRule, error)
	ListIDRules() []*model.IDRule
	UpdateIDRule(r *model.IDRule) error
	DeleteIDRule(id string) error

	// Segment
	CreateSegment(s *model.Segment) error
	GetSegment(id string) (*model.Segment, error)
	ListSegments() []*model.Segment
	ListSegmentsByBizTypeID(bizTypeID string) []*model.Segment
	GetActiveSegmentByBizTypeID(bizTypeID string) (*model.Segment, error)
	UpdateSegment(s *model.Segment) error
	DeleteSegment(id string) error

	// AllocRecord
	CreateAllocRecord(a *model.AllocRecord) error
	GetAllocRecord(id string) (*model.AllocRecord, error)
	ListAllocRecords() []*model.AllocRecord
	UpdateAllocRecord(a *model.AllocRecord) error
	DeleteAllocRecord(id string) error

	// Lease
	CreateLease(l *model.Lease) error
	GetLease(id string) (*model.Lease, error)
	ListLeases() []*model.Lease
	GetActiveLeaseByNodeAndBiz(nodeID, bizTypeID string) (*model.Lease, error)
	UpdateLease(l *model.Lease) error
	DeleteLease(id string) error

	// NodeHeartbeat
	CreateNodeHeartbeat(h *model.NodeHeartbeat) error
	GetNodeHeartbeat(id string) (*model.NodeHeartbeat, error)
	ListNodeHeartbeats() []*model.NodeHeartbeat
	GetLatestHeartbeatByNodeID(nodeID string) (*model.NodeHeartbeat, error)
	UpdateNodeHeartbeat(h *model.NodeHeartbeat) error
	DeleteNodeHeartbeat(id string) error

	// SnowflakeConfig
	CreateSnowflakeConfig(c *model.SnowflakeConfig) error
	GetSnowflakeConfig(id string) (*model.SnowflakeConfig, error)
	GetSnowflakeConfigByBizTypeID(bizTypeID string) (*model.SnowflakeConfig, error)
	ListSnowflakeConfigs() []*model.SnowflakeConfig
	UpdateSnowflakeConfig(c *model.SnowflakeConfig) error
	DeleteSnowflakeConfig(id string) error

	// AllocStats
	CreateAllocStats(s *model.AllocStats) error
	GetAllocStats(id string) (*model.AllocStats, error)
	ListAllocStats() []*model.AllocStats
	GetAllocStatsByBizNodeDate(bizTypeID, nodeID, date string) (*model.AllocStats, error)
	UpdateAllocStats(s *model.AllocStats) error
	DeleteAllocStats(id string) error

	// RecycleRecord
	CreateRecycleRecord(r *model.RecycleRecord) error
	GetRecycleRecord(id string) (*model.RecycleRecord, error)
	ListRecycleRecords() []*model.RecycleRecord
	UpdateRecycleRecord(r *model.RecycleRecord) error
	DeleteRecycleRecord(id string) error
}
