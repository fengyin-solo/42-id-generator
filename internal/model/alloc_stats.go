package model

import (
	"strings"
	"time"
)

// AllocStats 分配统计，按业务/节点/日期聚合。
type AllocStats struct {
	ID           string    `json:"id"`
	BizTypeID    string    `json:"biz_type_id"`
	NodeID       string    `json:"node_id"`
	Date         string    `json:"date"`
	TotalAllocated int64   `json:"total_allocated"`
	PeakQPS      float64   `json:"peak_qps"`
	AvgBatchSize float64   `json:"avg_batch_size"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (s *AllocStats) Validate() error {
	s.BizTypeID = strings.TrimSpace(s.BizTypeID)
	s.NodeID = strings.TrimSpace(s.NodeID)
	s.Date = strings.TrimSpace(s.Date)
	if s.BizTypeID == "" {
		return NewValidationError("biz_type_id", "业务类型不能为空")
	}
	if s.NodeID == "" {
		return NewValidationError("node_id", "节点不能为空")
	}
	if s.Date == "" {
		return NewValidationError("date", "日期不能为空")
	}
	if s.TotalAllocated < 0 {
		return NewValidationError("total_allocated", "分配总量不能为负数")
	}
	if s.PeakQPS < 0 {
		return NewValidationError("peak_qps", "峰值 QPS 不能为负数")
	}
	if s.AvgBatchSize < 0 {
		return NewValidationError("avg_batch_size", "平均批次大小不能为负数")
	}
	return nil
}

type AllocStatsFilter struct {
	BizTypeID string
	NodeID    string
	Date      string
	DateFrom  string
	DateTo    string
}

func (f AllocStatsFilter) Match(s *AllocStats) bool {
	if f.BizTypeID != "" && s.BizTypeID != f.BizTypeID {
		return false
	}
	if f.NodeID != "" && s.NodeID != f.NodeID {
		return false
	}
	if f.Date != "" && s.Date != f.Date {
		return false
	}
	if f.DateFrom != "" && s.Date < f.DateFrom {
		return false
	}
	if f.DateTo != "" && s.Date > f.DateTo {
		return false
	}
	return true
}
