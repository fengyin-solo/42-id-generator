package model

import (
	"strings"
	"time"
)

const (
	LeaseStatusActive  = "active"
	LeaseStatusExpired = "expired"
)

var leaseTransitions = map[string]map[string]bool{
	LeaseStatusActive:  {LeaseStatusExpired: true},
	LeaseStatusExpired: {},
}

// CanTransitionLease 判断租约状态是否允许流转。
func CanTransitionLease(from, to string) bool {
	if m, ok := leaseTransitions[from]; ok {
		return m[to]
	}
	return false
}

// Lease 租约，节点对某业务类型的发号权租约。
type Lease struct {
	ID        string    `json:"id"`
	NodeID    string    `json:"node_id"`
	BizTypeID string    `json:"biz_type_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RenewedAt time.Time `json:"renewed_at"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (l *Lease) Validate() error {
	l.NodeID = strings.TrimSpace(l.NodeID)
	l.BizTypeID = strings.TrimSpace(l.BizTypeID)
	if l.NodeID == "" {
		return NewValidationError("node_id", "节点不能为空")
	}
	if l.BizTypeID == "" {
		return NewValidationError("biz_type_id", "业务类型不能为空")
	}
	if l.ExpiresAt.IsZero() {
		return NewValidationError("expires_at", "过期时间不能为空")
	}
	if l.Status == "" {
		l.Status = LeaseStatusActive
	}
	if l.Status != LeaseStatusActive && l.Status != LeaseStatusExpired {
		return NewValidationError("status", "租约状态不合法")
	}
	return nil
}

// IsExpired 判断租约是否已过期。
func (l *Lease) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

type LeaseFilter struct {
	NodeID    string
	BizTypeID string
	Status    string
	Expired   *bool
}

func (f LeaseFilter) Match(l *Lease) bool {
	if f.NodeID != "" && l.NodeID != f.NodeID {
		return false
	}
	if f.BizTypeID != "" && l.BizTypeID != f.BizTypeID {
		return false
	}
	if f.Status != "" && l.Status != f.Status {
		return false
	}
	if f.Expired != nil {
		if *f.Expired != l.IsExpired() {
			return false
		}
	}
	return true
}
