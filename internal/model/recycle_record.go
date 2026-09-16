package model

import (
	"strings"
	"time"
)

const (
	RecycleStatusPending   = "pending"
	RecycleStatusRecycled  = "recycled"
	RecycleStatusRejected  = "rejected"
)

var recycleTransitions = map[string]map[string]bool{
	RecycleStatusPending:  {RecycleStatusRecycled: true, RecycleStatusRejected: true},
	RecycleStatusRecycled: {},
	RecycleStatusRejected: {},
}

// CanTransitionRecycle 判断回收记录状态是否允许流转。
func CanTransitionRecycle(from, to string) bool {
	if m, ok := recycleTransitions[from]; ok {
		return m[to]
	}
	return false
}

// RecycleRecord 号段回收记录。
type RecycleRecord struct {
	ID          string    `json:"id"`
	SegmentID   string    `json:"segment_id"`
	BizTypeID   string    `json:"biz_type_id"`
	Reason      string    `json:"reason"`
	RecycledAt  time.Time `json:"recycled_at"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *RecycleRecord) Validate() error {
	r.SegmentID = strings.TrimSpace(r.SegmentID)
	r.BizTypeID = strings.TrimSpace(r.BizTypeID)
	if r.SegmentID == "" {
		return NewValidationError("segment_id", "号段不能为空")
	}
	if r.BizTypeID == "" {
		return NewValidationError("biz_type_id", "业务类型不能为空")
	}
	if r.Reason == "" {
		return NewValidationError("reason", "回收原因不能为空")
	}
	if r.Status == "" {
		r.Status = RecycleStatusPending
	}
	if r.Status != RecycleStatusPending && r.Status != RecycleStatusRecycled && r.Status != RecycleStatusRejected {
		return NewValidationError("status", "回收状态不合法")
	}
	return nil
}

type RecycleRecordFilter struct {
	SegmentID string
	BizTypeID string
	Status    string
	Keyword   string
}

func (f RecycleRecordFilter) Match(r *RecycleRecord) bool {
	if f.SegmentID != "" && r.SegmentID != f.SegmentID {
		return false
	}
	if f.BizTypeID != "" && r.BizTypeID != f.BizTypeID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Reason), k) {
			return false
		}
	}
	return true
}
