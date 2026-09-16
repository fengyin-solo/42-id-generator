package model

import (
	"strings"
	"time"
)

const (
	SegmentStatusUsing    = "using"
	SegmentStatusUsed     = "used"
	SegmentStatusExhausted = "exhausted"
)

var segmentTransitions = map[string]map[string]bool{
	SegmentStatusUsing:    {SegmentStatusUsed: true, SegmentStatusExhausted: true},
	SegmentStatusUsed:     {SegmentStatusExhausted: true},
	SegmentStatusExhausted: {},
}

// CanTransitionSegment 判断号段状态是否允许流转。
func CanTransitionSegment(from, to string) bool {
	if m, ok := segmentTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Segment 号段，用于号段发号模式。
type Segment struct {
	ID        string    `json:"id"`
	BizTypeID string    `json:"biz_type_id"`
	StartID   int64     `json:"start_id"`
	EndID     int64     `json:"end_id"`
	Cursor    int64     `json:"cursor"`
	Status    string    `json:"status"`
	NodeID    string    `json:"node_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Segment) Validate() error {
	s.BizTypeID = strings.TrimSpace(s.BizTypeID)
	if s.BizTypeID == "" {
		return NewValidationError("biz_type_id", "业务类型不能为空")
	}
	if s.StartID < 0 {
		return NewValidationError("start_id", "起始 ID 不能为负数")
	}
	if s.EndID <= s.StartID {
		return NewValidationError("end_id", "结束 ID 必须大于起始 ID")
	}
	if s.Cursor < s.StartID {
		s.Cursor = s.StartID
	}
	if s.Status == "" {
		s.Status = SegmentStatusUsing
	}
	if s.Status != SegmentStatusUsing && s.Status != SegmentStatusUsed && s.Status != SegmentStatusExhausted {
		return NewValidationError("status", "号段状态不合法")
	}
	return nil
}

// IsExhausted 判断号段是否已耗尽。
func (s *Segment) IsExhausted() bool {
	return s.Cursor >= s.EndID
}

// Remain 返回剩余可用数量。
func (s *Segment) Remain() int64 {
	if s.Cursor >= s.EndID {
		return 0
	}
	return s.EndID - s.Cursor
}

type SegmentFilter struct {
	BizTypeID string
	Status    string
	NodeID    string
	Exhausted *bool
}

func (f SegmentFilter) Match(s *Segment) bool {
	if f.BizTypeID != "" && s.BizTypeID != f.BizTypeID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.NodeID != "" && s.NodeID != f.NodeID {
		return false
	}
	if f.Exhausted != nil {
		if *f.Exhausted != s.IsExhausted() {
			return false
		}
	}
	return true
}
