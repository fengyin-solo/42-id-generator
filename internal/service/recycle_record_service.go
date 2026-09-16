package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateRecycleRecord(input model.RecycleRecord) (*model.RecycleRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// 外键校验
	if _, err := s.store.GetSegment(input.SegmentID); err != nil {
		return nil, model.NewValidationError("segment_id", "号段不存在")
	}
	if _, err := s.store.GetBizType(input.BizTypeID); err != nil {
		return nil, model.NewValidationError("biz_type_id", "业务类型不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateRecycleRecord(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建回收记录: %s segment=%s", input.ID, input.SegmentID)
	return &input, nil
}

func (s *Service) GetRecycleRecord(id string) (*model.RecycleRecord, error) {
	return s.store.GetRecycleRecord(id)
}

func (s *Service) ListRecycleRecords(filter model.RecycleRecordFilter, page, size int) ([]*model.RecycleRecord, int, error) {
	all := s.store.ListRecycleRecords()
	matched := make([]*model.RecycleRecord, 0, len(all))
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
		return []*model.RecycleRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRecycleRecord(id string, input model.RecycleRecord) (*model.RecycleRecord, error) {
	r, err := s.store.GetRecycleRecord(id)
	if err != nil {
		return nil, err
	}
	// 状态机校验
	if input.Status != "" && input.Status != r.Status {
		if !model.CanTransitionRecycle(r.Status, input.Status) {
			return nil, model.NewValidationError("status", "回收状态不允许从 "+r.Status+" 流转到 "+input.Status)
		}
		r.Status = input.Status
		if input.Status == model.RecycleStatusRecycled {
			r.RecycledAt = time.Now()
			// 实际回收号段：标记号段为 exhausted
			if seg, err := s.store.GetSegment(r.SegmentID); err == nil {
				seg.Status = model.SegmentStatusExhausted
				seg.UpdatedAt = time.Now()
				_ = s.store.UpdateSegment(seg)
			}
		}
	}
	if input.Reason != "" {
		r.Reason = input.Reason
	}
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRecycleRecord(r); err != nil {
		return nil, err
	}
	s.log.Infof("更新回收记录: %s status=%s", id, r.Status)
	return r, nil
}

func (s *Service) DeleteRecycleRecord(id string) error {
	if _, err := s.store.GetRecycleRecord(id); err != nil {
		return err
	}
	return s.store.DeleteRecycleRecord(id)
}

func (s *Service) CountRecycleRecordsByStatus(status string) int {
	count := 0
	for _, r := range s.store.ListRecycleRecords() {
		if r.Status == status {
			count++
		}
	}
	return count
}
