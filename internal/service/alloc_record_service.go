package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateAllocRecord(input model.AllocRecord) (*model.AllocRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// 外键校验
	if _, err := s.store.GetBizType(input.BizTypeID); err != nil {
		return nil, model.NewValidationError("biz_type_id", "业务类型不存在")
	}
	if _, err := s.store.GetMachineNodeByNodeID(input.NodeID); err != nil {
		return nil, model.NewValidationError("node_id", "节点不存在")
	}
	if input.SegmentID != "" {
		if _, err := s.store.GetSegment(input.SegmentID); err != nil {
			return nil, model.NewValidationError("segment_id", "号段不存在")
		}
	}
	input.ID = idgen.Hex()
	input.AllocatedAt = time.Now()
	input.CreatedAt = input.AllocatedAt
	if err := s.store.CreateAllocRecord(&input); err != nil {
		return nil, err
	}
	// 更新统计
	dateStr := input.AllocatedAt.Format("2006-01-02")
	stats, err := s.store.GetAllocStatsByBizNodeDate(input.BizTypeID, input.NodeID, dateStr)
	if err != nil {
		stats = &model.AllocStats{
			ID:             idgen.Hex(),
			BizTypeID:      input.BizTypeID,
			NodeID:         input.NodeID,
			Date:           dateStr,
			TotalAllocated: input.BatchSize,
			PeakQPS:        1,
			AvgBatchSize:   float64(input.BatchSize),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		_ = s.store.CreateAllocStats(stats)
	} else {
		stats.TotalAllocated += input.BatchSize
		if input.BatchSize > int64(stats.PeakQPS) {
			stats.PeakQPS = float64(input.BatchSize)
		}
		// 简单更新平均批次大小
		stats.AvgBatchSize = float64(stats.TotalAllocated) / 100.0
		stats.UpdatedAt = time.Now()
		_ = s.store.UpdateAllocStats(stats)
	}
	s.log.Infof("创建分配记录: %s", input.ID)
	return &input, nil
}

func (s *Service) GetAllocRecord(id string) (*model.AllocRecord, error) {
	return s.store.GetAllocRecord(id)
}

func (s *Service) ListAllocRecords(filter model.AllocRecordFilter, page, size int) ([]*model.AllocRecord, int, error) {
	all := s.store.ListAllocRecords()
	matched := make([]*model.AllocRecord, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].AllocatedAt.After(matched[j].AllocatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AllocRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteAllocRecord(id string) error {
	if _, err := s.store.GetAllocRecord(id); err != nil {
		return err
	}
	return s.store.DeleteAllocRecord(id)
}

func (s *Service) CountAllocRecordsByBizType(bizTypeID string) int {
	count := 0
	for _, a := range s.store.ListAllocRecords() {
		if a.BizTypeID == bizTypeID {
			count++
		}
	}
	return count
}

func (s *Service) SumAllocatedByDate(date string) int64 {
	var sum int64
	for _, st := range s.store.ListAllocStats() {
		if st.Date == date {
			sum += st.TotalAllocated
		}
	}
	return sum
}
