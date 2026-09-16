package model

import (
	"strings"
	"time"
)

// AllocRecord 分配记录，记录每次号段/ID 分配的详情。
type AllocRecord struct {
	ID         string    `json:"id"`
	BizTypeID  string    `json:"biz_type_id"`
	NodeID     string    `json:"node_id"`
	SegmentID  string    `json:"segment_id"`
	BatchSize  int64     `json:"batch_size"`
	StartID    int64     `json:"start_id"`
	EndID      int64     `json:"end_id"`
	Mode       string    `json:"mode"`
	AllocatedAt time.Time `json:"allocated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a *AllocRecord) Validate() error {
	a.BizTypeID = strings.TrimSpace(a.BizTypeID)
	a.NodeID = strings.TrimSpace(a.NodeID)
	if a.BizTypeID == "" {
		return NewValidationError("biz_type_id", "业务类型不能为空")
	}
	if a.NodeID == "" {
		return NewValidationError("node_id", "节点不能为空")
	}
	if a.BatchSize <= 0 {
		return NewValidationError("batch_size", "批次大小必须大于 0")
	}
	if a.Mode == "" {
		return NewValidationError("mode", "模式不能为空")
	}
	if a.Mode != BizModeSegment && a.Mode != BizModeSnowflake && a.Mode != BizModeSequence {
		return NewValidationError("mode", "发号模式不合法")
	}
	return nil
}

type AllocRecordFilter struct {
	BizTypeID string
	NodeID    string
	SegmentID string
	Mode      string
	StartDate *time.Time
	EndDate   *time.Time
}

func (f AllocRecordFilter) Match(a *AllocRecord) bool {
	if f.BizTypeID != "" && a.BizTypeID != f.BizTypeID {
		return false
	}
	if f.NodeID != "" && a.NodeID != f.NodeID {
		return false
	}
	if f.SegmentID != "" && a.SegmentID != f.SegmentID {
		return false
	}
	if f.Mode != "" && a.Mode != f.Mode {
		return false
	}
	if f.StartDate != nil && a.AllocatedAt.Before(*f.StartDate) {
		return false
	}
	if f.EndDate != nil && a.AllocatedAt.After(*f.EndDate) {
		return false
	}
	return true
}
