package model

import (
	"strings"
	"time"
)

const (
	BizModeSegment   = "segment"
	BizModeSnowflake = "snowflake"
	BizModeSequence  = "sequence"
)

const (
	BizStatusEnabled  = "enabled"
	BizStatusDisabled = "disabled"
)

// BizType 业务标识，每种业务对应一种发号模式。
type BizType struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Mode        string    `json:"mode"`
	SegmentStep int       `json:"segment_step"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *BizType) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	b.Code = strings.TrimSpace(b.Code)
	b.Mode = strings.TrimSpace(b.Mode)
	if b.Name == "" {
		return NewValidationError("name", "业务名称不能为空")
	}
	if b.Code == "" {
		return NewValidationError("code", "业务编码不能为空")
	}
	if b.Mode != BizModeSegment && b.Mode != BizModeSnowflake && b.Mode != BizModeSequence {
		return NewValidationError("mode", "发号模式必须是 segment/snowflake/sequence 之一")
	}
	if b.Mode == BizModeSegment && b.SegmentStep <= 0 {
		return NewValidationError("segment_step", "号段步长必须大于 0")
	}
	if b.Status == "" {
		b.Status = BizStatusEnabled
	}
	if b.Status != BizStatusEnabled && b.Status != BizStatusDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type BizTypeFilter struct {
	Mode    string
	Status  string
	Keyword string
	Enabled *bool
}

func (f BizTypeFilter) Match(b *BizType) bool {
	if f.Mode != "" && b.Mode != f.Mode {
		return false
	}
	if f.Status != "" && b.Status != f.Status {
		return false
	}
	if f.Enabled != nil && b.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(b.Name), k) &&
			!strings.Contains(strings.ToLower(b.Code), k) {
			return false
		}
	}
	return true
}
