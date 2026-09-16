package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateIDRule(input model.IDRule) (*model.IDRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// 外键校验：业务类型必须存在
	if _, err := s.store.GetBizType(input.BizTypeID); err != nil {
		return nil, model.NewValidationError("biz_type_id", "业务类型不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateIDRule(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建发号规则: %s", input.ID)
	return &input, nil
}

func (s *Service) GetIDRule(id string) (*model.IDRule, error) {
	return s.store.GetIDRule(id)
}

func (s *Service) GetIDRuleByBizTypeID(bizTypeID string) (*model.IDRule, error) {
	return s.store.GetIDRuleByBizTypeID(bizTypeID)
}

func (s *Service) ListIDRules(filter model.IDRuleFilter, page, size int) ([]*model.IDRule, int, error) {
	all := s.store.ListIDRules()
	matched := make([]*model.IDRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.IDRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateIDRule(id string, input model.IDRule) (*model.IDRule, error) {
	r, err := s.store.GetIDRule(id)
	if err != nil {
		return nil, err
	}
	// 外键校验
	if input.BizTypeID != "" && input.BizTypeID != r.BizTypeID {
		if _, err := s.store.GetBizType(input.BizTypeID); err != nil {
			return nil, model.NewValidationError("biz_type_id", "业务类型不存在")
		}
	}
	r.Name = input.Name
	r.Mode = input.Mode
	r.SignBits = input.SignBits
	r.TimestampBits = input.TimestampBits
	r.DatacenterBits = input.DatacenterBits
	r.WorkerBits = input.WorkerBits
	r.SequenceBits = input.SequenceBits
	r.Enabled = input.Enabled
	if input.Status != "" {
		r.Status = input.Status
	}
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateIDRule(r); err != nil {
		return nil, err
	}
	s.log.Infof("更新发号规则: %s", id)
	return r, nil
}

func (s *Service) DeleteIDRule(id string) error {
	if _, err := s.store.GetIDRule(id); err != nil {
		return err
	}
	s.log.Infof("删除发号规则: %s", id)
	return s.store.DeleteIDRule(id)
}

func (s *Service) ToggleIDRuleEnabled(id string) (*model.IDRule, error) {
	r, err := s.store.GetIDRule(id)
	if err != nil {
		return nil, err
	}
	r.Enabled = !r.Enabled
	if r.Enabled {
		r.Status = model.IDRuleStatusEnabled
	} else {
		r.Status = model.IDRuleStatusDisabled
	}
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateIDRule(r); err != nil {
		return nil, err
	}
	return r, nil
}
