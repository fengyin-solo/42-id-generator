package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateSnowflakeConfig(input model.SnowflakeConfig) (*model.SnowflakeConfig, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// 外键校验
	if _, err := s.store.GetBizType(input.BizTypeID); err != nil {
		return nil, model.NewValidationError("biz_type_id", "业务类型不存在")
	}
	// 检查是否已存在
	if _, err := s.store.GetSnowflakeConfigByBizTypeID(input.BizTypeID); err == nil {
		return nil, model.NewValidationError("biz_type_id", "该业务类型已存在雪花配置")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateSnowflakeConfig(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建雪花配置: %s biz=%s", input.ID, input.BizTypeID)
	return &input, nil
}

func (s *Service) GetSnowflakeConfig(id string) (*model.SnowflakeConfig, error) {
	return s.store.GetSnowflakeConfig(id)
}

func (s *Service) GetSnowflakeConfigByBizTypeID(bizTypeID string) (*model.SnowflakeConfig, error) {
	return s.store.GetSnowflakeConfigByBizTypeID(bizTypeID)
}

func (s *Service) ListSnowflakeConfigs(filter model.SnowflakeConfigFilter, page, size int) ([]*model.SnowflakeConfig, int, error) {
	all := s.store.ListSnowflakeConfigs()
	matched := make([]*model.SnowflakeConfig, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.SnowflakeConfig{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSnowflakeConfig(id string, input model.SnowflakeConfig) (*model.SnowflakeConfig, error) {
	c, err := s.store.GetSnowflakeConfig(id)
	if err != nil {
		return nil, err
	}
	if input.EpochMs > 0 {
		c.EpochMs = input.EpochMs
	}
	if input.DatacenterBits > 0 {
		c.DatacenterBits = input.DatacenterBits
	}
	if input.WorkerBits > 0 {
		c.WorkerBits = input.WorkerBits
	}
	if input.SequenceBits > 0 {
		c.SequenceBits = input.SequenceBits
	}
	if input.Twepoch > 0 {
		c.Twepoch = input.Twepoch
	}
	c.UpdatedAt = time.Now()
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSnowflakeConfig(c); err != nil {
		return nil, err
	}
	s.log.Infof("更新雪花配置: %s", id)
	return c, nil
}

func (s *Service) DeleteSnowflakeConfig(id string) error {
	if _, err := s.store.GetSnowflakeConfig(id); err != nil {
		return err
	}
	return s.store.DeleteSnowflakeConfig(id)
}
