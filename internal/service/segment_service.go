package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateSegment(input model.Segment) (*model.Segment, error) {
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
	if err := s.store.CreateSegment(&input); err != nil {
		return nil, err
	}
	s.log.Infof("创建号段: %s [%d, %d)", input.ID, input.StartID, input.EndID)
	return &input, nil
}

func (s *Service) GetSegment(id string) (*model.Segment, error) {
	return s.store.GetSegment(id)
}

func (s *Service) ListSegments(filter model.SegmentFilter, page, size int) ([]*model.Segment, int, error) {
	all := s.store.ListSegments()
	matched := make([]*model.Segment, 0, len(all))
	for _, seg := range all {
		if filter.Match(seg) {
			matched = append(matched, seg)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Segment{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSegment(id string, input model.Segment) (*model.Segment, error) {
	seg, err := s.store.GetSegment(id)
	if err != nil {
		return nil, err
	}
	// 状态机校验
	if input.Status != "" && input.Status != seg.Status {
		if !model.CanTransitionSegment(seg.Status, input.Status) {
			return nil, model.NewValidationError("status", "号段状态不允许从 "+seg.Status+" 流转到 "+input.Status)
		}
		seg.Status = input.Status
	}
	// 游标推进校验
	if input.Cursor > 0 {
		if input.Cursor < seg.StartID || input.Cursor > seg.EndID {
			return nil, model.NewValidationError("cursor", "游标必须在号段范围内")
		}
		seg.Cursor = input.Cursor
	}
	if input.NodeID != "" {
		seg.NodeID = input.NodeID
	}
	seg.UpdatedAt = time.Now()
	if err := seg.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSegment(seg); err != nil {
		return nil, err
	}
	s.log.Infof("更新号段: %s cursor=%d", id, seg.Cursor)
	return seg, nil
}

func (s *Service) DeleteSegment(id string) error {
	seg, err := s.store.GetSegment(id)
	if err != nil {
		return err
	}
	if seg.Status == model.SegmentStatusUsing {
		return model.NewValidationError("segment", "正在使用的号段不能删除")
	}
	s.log.Infof("删除号段: %s", id)
	return s.store.DeleteSegment(id)
}

func (s *Service) AdvanceSegmentCursor(id string, n int) (*model.Segment, []int64, error) {
	seg, err := s.store.GetSegment(id)
	if err != nil {
		return nil, nil, err
	}
	if seg.Status != model.SegmentStatusUsing {
		return nil, nil, model.NewValidationError("status", "只有 using 状态的号段可以取号")
	}
	ids := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		if seg.Cursor >= seg.EndID {
			seg.Status = model.SegmentStatusExhausted
			break
		}
		ids = append(ids, seg.Cursor)
		seg.Cursor++
	}
	if seg.Cursor >= seg.EndID && seg.Status == model.SegmentStatusUsing {
		seg.Status = model.SegmentStatusExhausted
	}
	seg.UpdatedAt = time.Now()
	if err := s.store.UpdateSegment(seg); err != nil {
		return nil, nil, err
	}
	s.log.Infof("号段 %s 推进游标: %d 个", id, len(ids))
	return seg, ids, nil
}

func (s *Service) SwitchSegment(bizTypeID string) (*model.Segment, error) {
	// 查找该业务类型下一个可用的 using 号段
	for _, seg := range s.store.ListSegmentsByBizTypeID(bizTypeID) {
		if seg.Status == model.SegmentStatusUsing && !seg.IsExhausted() {
			return seg, nil
		}
	}
	return nil, model.NewValidationError("segment", "没有可用的号段")
}

func (s *Service) CountSegmentsByStatus(status string) int {
	count := 0
	for _, seg := range s.store.ListSegments() {
		if seg.Status == status {
			count++
		}
	}
	return count
}

func (s *Service) CountExhaustedSegments() int {
	count := 0
	for _, seg := range s.store.ListSegments() {
		if seg.IsExhausted() {
			count++
		}
	}
	return count
}
