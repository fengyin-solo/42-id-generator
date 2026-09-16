package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateAllocStats(input model.AllocStats) (*model.AllocStats, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if err := s.store.CreateAllocStats(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetAllocStats(id string) (*model.AllocStats, error) {
	return s.store.GetAllocStats(id)
}

func (s *Service) ListAllocStats(filter model.AllocStatsFilter, page, size int) ([]*model.AllocStats, int, error) {
	all := s.store.ListAllocStats()
	matched := make([]*model.AllocStats, 0, len(all))
	for _, st := range all {
		if filter.Match(st) {
			matched = append(matched, st)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Date != matched[j].Date {
			return matched[i].Date > matched[j].Date
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AllocStats{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAllocStats(id string, input model.AllocStats) (*model.AllocStats, error) {
	st, err := s.store.GetAllocStats(id)
	if err != nil {
		return nil, err
	}
	st.TotalAllocated = input.TotalAllocated
	st.PeakQPS = input.PeakQPS
	st.AvgBatchSize = input.AvgBatchSize
	st.UpdatedAt = time.Now()
	if err := st.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAllocStats(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) DeleteAllocStats(id string) error {
	if _, err := s.store.GetAllocStats(id); err != nil {
		return err
	}
	return s.store.DeleteAllocStats(id)
}

func (s *Service) GetStatsOverview() map[string]interface{} {
	bizCount := len(s.store.ListBizTypes())
	activeNodes := 0
	for _, n := range s.store.ListMachineNodes() {
		if n.Status == model.NodeStatusActive {
			activeNodes++
		}
	}
	today := time.Now().Format("2006-01-02")
	var todayAllocated int64
	var peakQPS float64
	for _, st := range s.store.ListAllocStats() {
		if st.Date == today {
			todayAllocated += st.TotalAllocated
			if st.PeakQPS > peakQPS {
				peakQPS = st.PeakQPS
			}
		}
	}
	return map[string]interface{}{
		"biz_count":       bizCount,
		"active_nodes":    activeNodes,
		"today_allocated": todayAllocated,
		"peak_qps":        peakQPS,
	}
}

func (s *Service) GetTopAllocStatsByDate(date string, topN int) []*model.AllocStats {
	list := make([]*model.AllocStats, 0)
	for _, st := range s.store.ListAllocStats() {
		if st.Date == date {
			list = append(list, st)
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].TotalAllocated > list[j].TotalAllocated
	})
	if topN > len(list) {
		topN = len(list)
	}
	return list[:topN]
}

func (s *Service) GetStatsGroupByBizType() map[string]int64 {
	result := make(map[string]int64)
	for _, st := range s.store.ListAllocStats() {
		result[st.BizTypeID] += st.TotalAllocated
	}
	return result
}

func (s *Service) GetStatsGroupByNode() map[string]int64 {
	result := make(map[string]int64)
	for _, st := range s.store.ListAllocStats() {
		result[st.NodeID] += st.TotalAllocated
	}
	return result
}
