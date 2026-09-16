package model

import (
	"strings"
	"time"
)

// SnowflakeConfig 雪花算法配置，按业务类型存储。
type SnowflakeConfig struct {
	ID             string    `json:"id"`
	BizTypeID      string    `json:"biz_type_id"`
	EpochMs        int64     `json:"epoch_ms"`
	DatacenterBits int       `json:"datacenter_bits"`
	WorkerBits     int       `json:"worker_bits"`
	SequenceBits   int       `json:"sequence_bits"`
	Twepoch        int64     `json:"twepoch"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (c *SnowflakeConfig) Validate() error {
	c.BizTypeID = strings.TrimSpace(c.BizTypeID)
	if c.BizTypeID == "" {
		return NewValidationError("biz_type_id", "业务类型不能为空")
	}
	if c.DatacenterBits < 0 || c.WorkerBits < 0 || c.SequenceBits < 0 {
		return NewValidationError("bits", "位数不能为负数")
	}
	total := c.DatacenterBits + c.WorkerBits + c.SequenceBits
	if total > 63 {
		return NewValidationError("bits", "总位数不能超过 63（不含符号位）")
	}
	if c.Twepoch <= 0 {
		return NewValidationError("twepoch", "纪元时间戳必须大于 0")
	}
	return nil
}

type SnowflakeConfigFilter struct {
	BizTypeID string
	Keyword   string
}

func (f SnowflakeConfigFilter) Match(c *SnowflakeConfig) bool {
	if f.BizTypeID != "" && c.BizTypeID != f.BizTypeID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.BizTypeID), k) {
			return false
		}
	}
	return true
}
