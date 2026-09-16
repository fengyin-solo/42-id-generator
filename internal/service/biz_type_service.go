package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateBizType(input model.BizType) (*model.BizType, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateBizType(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建业务类型: %s %s", input.ID, input.Code)
	return &input, nil
}

func (s *Service) GetBizType(id string) (*model.BizType, error) {
	return s.store.GetBizType(id)
}

func (s *Service) GetBizTypeByCode(code string) (*model.BizType, error) {
	return s.store.GetBizTypeByCode(code)
}

func (s *Service) ListBizTypes(filter model.BizTypeFilter, page, size int) ([]*model.BizType, int, error) {
	all := s.store.ListBizTypes()
	matched := make([]*model.BizType, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.BizType{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateBizType(id string, input model.BizType) (*model.BizType, error) {
	b, err := s.store.GetBizType(id)
	if err != nil {
		return nil, err
	}
	b.Name = input.Name
	b.Description = input.Description
	if input.SegmentStep > 0 {
		b.SegmentStep = input.SegmentStep
	}
	b.Enabled = input.Enabled
	if input.Status != "" {
		b.Status = input.Status
	}
	b.UpdatedAt = time.Now()
	if err := b.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateBizType(b); err != nil {
		return nil, err
	}
	s.log.Infof("更新业务类型: %s", id)
	return b, nil
}

func (s *Service) DeleteBizType(id string) error {
	if _, err := s.store.GetBizType(id); err != nil {
		return err
	}
	// 检查是否有关联的号段
	segments := s.store.ListSegmentsByBizTypeID(id)
	if len(segments) > 0 {
		return model.NewValidationError("biz_type", "该业务类型下存在号段，不能删除")
	}
	// 检查是否有关联规则
	if _, err := s.store.GetIDRuleByBizTypeID(id); err == nil {
		return model.NewValidationError("biz_type", "该业务类型下存在发号规则，不能删除")
	}
	// 检查是否有关联配置
	if _, err := s.store.GetSnowflakeConfigByBizTypeID(id); err == nil {
		return model.NewValidationError("biz_type", "该业务类型下存在雪花配置，不能删除")
	}
	s.log.Infof("删除业务类型: %s", id)
	return s.store.DeleteBizType(id)
}

func (s *Service) ToggleBizTypeStatus(id string) (*model.BizType, error) {
	b, err := s.store.GetBizType(id)
	if err != nil {
		return nil, err
	}
	if b.Status == model.BizStatusEnabled {
		b.Status = model.BizStatusDisabled
		b.Enabled = false
	} else {
		b.Status = model.BizStatusEnabled
		b.Enabled = true
	}
	b.UpdatedAt = time.Now()
	if err := s.store.UpdateBizType(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) CountBizTypes() int {
	return len(s.store.ListBizTypes())
}

func (s *Service) CountBizTypesByMode(mode string) int {
	count := 0
	for _, b := range s.store.ListBizTypes() {
		if b.Mode == mode {
			count++
		}
	}
	return count
}
