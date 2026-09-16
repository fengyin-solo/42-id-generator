package model

import (
	"strings"
	"time"
)

const (
	IDRuleStatusEnabled  = "enabled"
	IDRuleStatusDisabled = "disabled"
)

// IDRule 发号规则，定义比特位布局。
type IDRule struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	BizTypeID      string    `json:"biz_type_id"`
	Mode           string    `json:"mode"`
	SignBits       int       `json:"sign_bits"`
	TimestampBits  int       `json:"timestamp_bits"`
	DatacenterBits int       `json:"datacenter_bits"`
	WorkerBits     int       `json:"worker_bits"`
	SequenceBits   int       `json:"sequence_bits"`
	Enabled        bool      `json:"enabled"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (r *IDRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.BizTypeID = strings.TrimSpace(r.BizTypeID)
	r.Mode = strings.TrimSpace(r.Mode)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.BizTypeID == "" {
		return NewValidationError("biz_type_id", "业务类型不能为空")
	}
	if r.Mode != BizModeSegment && r.Mode != BizModeSnowflake && r.Mode != BizModeSequence {
		return NewValidationError("mode", "发号模式不合法")
	}
	total := r.SignBits + r.TimestampBits + r.DatacenterBits + r.WorkerBits + r.SequenceBits
	if total != 64 {
		return NewValidationError("bit_layout", "总位数必须等于 64")
	}
	if r.SignBits < 0 || r.TimestampBits < 0 || r.DatacenterBits < 0 || r.WorkerBits < 0 || r.SequenceBits < 0 {
		return NewValidationError("bit_layout", "各段位数不能为负数")
	}
	if r.Status == "" {
		r.Status = IDRuleStatusEnabled
	}
	if r.Status != IDRuleStatusEnabled && r.Status != IDRuleStatusDisabled {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

type IDRuleFilter struct {
	BizTypeID string
	Mode      string
	Status    string
	Keyword   string
	Enabled   *bool
}

func (f IDRuleFilter) Match(r *IDRule) bool {
	if f.BizTypeID != "" && r.BizTypeID != f.BizTypeID {
		return false
	}
	if f.Mode != "" && r.Mode != f.Mode {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Enabled != nil && r.Enabled != *f.Enabled {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) {
			return false
		}
	}
	return true
}
