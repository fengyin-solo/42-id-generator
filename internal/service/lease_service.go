package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateLease(input model.Lease) (*model.Lease, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// 外键校验
	if _, err := s.store.GetMachineNodeByNodeID(input.NodeID); err != nil {
		return nil, model.NewValidationError("node_id", "节点不存在")
	}
	if _, err := s.store.GetBizType(input.BizTypeID); err != nil {
		return nil, model.NewValidationError("biz_type_id", "业务类型不存在")
	}
	// 检查是否已存在活跃租约
	if existing, err := s.store.GetActiveLeaseByNodeAndBiz(input.NodeID, input.BizTypeID); err == nil {
		return nil, model.NewValidationError("lease", "该节点已存在此业务类型的活跃租约: "+existing.ID)
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	input.RenewedAt = input.CreatedAt
	if err := s.store.CreateLease(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建租约: %s node=%s biz=%s", input.ID, input.NodeID, input.BizTypeID)
	return &input, nil
}

func (s *Service) GetLease(id string) (*model.Lease, error) {
	return s.store.GetLease(id)
}

func (s *Service) ListLeases(filter model.LeaseFilter, page, size int) ([]*model.Lease, int, error) {
	all := s.store.ListLeases()
	matched := make([]*model.Lease, 0, len(all))
	for _, l := range all {
		if filter.Match(l) {
			matched = append(matched, l)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Lease{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) RenewLease(id string, duration time.Duration) (*model.Lease, error) {
	l, err := s.store.GetLease(id)
	if err != nil {
		return nil, err
	}
	if l.Status != model.LeaseStatusActive {
		return nil, model.NewValidationError("status", "只有活跃租约可以续期")
	}
	if l.IsExpired() {
		l.Status = model.LeaseStatusExpired
		l.UpdatedAt = time.Now()
		_ = s.store.UpdateLease(l)
		return nil, model.NewValidationError("lease", "租约已过期，无法续期")
	}
	l.ExpiresAt = time.Now().Add(duration)
	l.RenewedAt = time.Now()
	l.UpdatedAt = time.Now()
	if err := s.store.UpdateLease(l); err != nil {
		return nil, err
	}
	s.log.Infof("续期租约: %s expires=%s", id, l.ExpiresAt.Format(time.RFC3339))
	return l, nil
}

func (s *Service) ExpireLease(id string) (*model.Lease, error) {
	l, err := s.store.GetLease(id)
	if err != nil {
		return nil, err
	}
	if l.Status == model.LeaseStatusExpired {
		return l, nil
	}
	if !model.CanTransitionLease(l.Status, model.LeaseStatusExpired) {
		return nil, model.NewValidationError("status", "租约状态不允许流转到 expired")
	}
	l.Status = model.LeaseStatusExpired
	l.ExpiresAt = time.Now()
	l.UpdatedAt = time.Now()
	if err := s.store.UpdateLease(l); err != nil {
		return nil, err
	}
	// 过期后，释放关联的号段
	for _, seg := range s.store.ListSegments() {
		if seg.NodeID == l.NodeID && seg.Status == model.SegmentStatusUsing {
			seg.Status = model.SegmentStatusUsed
			seg.NodeID = ""
			seg.UpdatedAt = time.Now()
			_ = s.store.UpdateSegment(seg)
		}
	}
	s.log.Infof("租约过期: %s node=%s", id, l.NodeID)
	return l, nil
}

func (s *Service) DeleteLease(id string) error {
	if _, err := s.store.GetLease(id); err != nil {
		return err
	}
	return s.store.DeleteLease(id)
}

func (s *Service) CountActiveLeases() int {
	count := 0
	now := time.Now()
	for _, l := range s.store.ListLeases() {
		if l.Status == model.LeaseStatusActive && l.ExpiresAt.After(now) {
			count++
		}
	}
	return count
}

func (s *Service) ExpireOverdueLeases() int {
	expired := 0
	now := time.Now()
	for _, l := range s.store.ListLeases() {
		if l.Status == model.LeaseStatusActive && l.ExpiresAt.Before(now) {
			l.Status = model.LeaseStatusExpired
			l.UpdatedAt = now
			_ = s.store.UpdateLease(l)
			// 释放号段
			for _, seg := range s.store.ListSegments() {
				if seg.NodeID == l.NodeID && seg.Status == model.SegmentStatusUsing {
					seg.Status = model.SegmentStatusUsed
					seg.NodeID = ""
					seg.UpdatedAt = now
					_ = s.store.UpdateSegment(seg)
				}
			}
			expired++
		}
	}
	return expired
}
