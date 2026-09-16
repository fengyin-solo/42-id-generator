package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateNodeHeartbeat(input model.NodeHeartbeat) (*model.NodeHeartbeat, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	// 外键校验：节点必须存在
	if _, err := s.store.GetMachineNodeByNodeID(input.NodeID); err != nil {
		return nil, model.NewValidationError("node_id", "节点不存在")
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	if err := s.store.CreateNodeHeartbeat(&input); err != nil {
		return nil, err
	}
	// 更新节点最后心跳时间
	if node, err := s.store.GetMachineNodeByNodeID(input.NodeID); err == nil {
		node.LastBeatAt = input.BeatAt
		node.UpdatedAt = time.Now()
		_ = s.store.UpdateMachineNode(node)
	}
	s.log.Infof("节点心跳: %s load=%.1f cpu=%.1f mem=%.1f", input.NodeID, input.Load, input.CPUUsage, input.MemoryUsage)
	return &input, nil
}

func (s *Service) GetNodeHeartbeat(id string) (*model.NodeHeartbeat, error) {
	return s.store.GetNodeHeartbeat(id)
}

func (s *Service) ListNodeHeartbeats(filter model.NodeHeartbeatFilter, page, size int) ([]*model.NodeHeartbeat, int, error) {
	all := s.store.ListNodeHeartbeats()
	matched := make([]*model.NodeHeartbeat, 0, len(all))
	for _, h := range all {
		if filter.Match(h) {
			matched = append(matched, h)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].BeatAt.After(matched[j].BeatAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.NodeHeartbeat{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetLatestHeartbeatByNodeID(nodeID string) (*model.NodeHeartbeat, error) {
	return s.store.GetLatestHeartbeatByNodeID(nodeID)
}

func (s *Service) DeleteNodeHeartbeat(id string) error {
	if _, err := s.store.GetNodeHeartbeat(id); err != nil {
		return err
	}
	return s.store.DeleteNodeHeartbeat(id)
}

func (s *Service) CheckNodeHealth(nodeID string) (healthy bool, lastBeat time.Time, err error) {
	h, err := s.store.GetLatestHeartbeatByNodeID(nodeID)
	if err != nil {
		return false, time.Time{}, err
	}
	// 超过 30 秒无心跳视为不健康
	threshold := 30 * time.Second
	healthy = time.Since(h.BeatAt) < threshold
	return healthy, h.BeatAt, nil
}

func (s *Service) CleanupOldHeartbeats(before time.Time) int {
	count := 0
	for _, h := range s.store.ListNodeHeartbeats() {
		if h.BeatAt.Before(before) {
			_ = s.store.DeleteNodeHeartbeat(h.ID)
			count++
		}
	}
	return count
}

func (s *Service) GetNodeLoadAvg(nodeID string) (float64, error) {
	var sum float64
	var cnt int
	for _, h := range s.store.ListNodeHeartbeats() {
		if h.NodeID == nodeID {
			sum += h.Load
			cnt++
		}
	}
	if cnt == 0 {
		return 0, model.NewValidationError("node_id", "该节点无心跳记录")
	}
	return sum / float64(cnt), nil
}
