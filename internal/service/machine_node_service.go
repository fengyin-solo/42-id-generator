package service

import (
	"sort"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/idgen"
)

func (s *Service) CreateMachineNode(input model.MachineNode) (*model.MachineNode, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	input.UpdatedAt = input.CreatedAt
	if input.LastBeatAt.IsZero() {
		input.LastBeatAt = input.CreatedAt
	}
	if err := s.store.CreateMachineNode(&input); err != nil {
		return nil, err
	}
	s.log.Infof("注册机器节点: %s %s", input.ID, input.NodeID)
	return &input, nil
}

func (s *Service) GetMachineNode(id string) (*model.MachineNode, error) {
	return s.store.GetMachineNode(id)
}

func (s *Service) GetMachineNodeByNodeID(nodeID string) (*model.MachineNode, error) {
	return s.store.GetMachineNodeByNodeID(nodeID)
}

func (s *Service) ListMachineNodes(filter model.MachineNodeFilter, page, size int) ([]*model.MachineNode, int, error) {
	all := s.store.ListMachineNodes()
	matched := make([]*model.MachineNode, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.MachineNode{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateMachineNode(id string, input model.MachineNode) (*model.MachineNode, error) {
	n, err := s.store.GetMachineNode(id)
	if err != nil {
		return nil, err
	}
	// 状态机校验
	if input.Status != "" && input.Status != n.Status {
		if !model.CanTransitionNode(n.Status, input.Status) {
			return nil, model.NewValidationError("status", "节点状态不允许从 "+n.Status+" 流转到 "+input.Status)
		}
		n.Status = input.Status
	}
	if input.Hostname != "" {
		n.Hostname = input.Hostname
	}
	if input.IP != "" {
		n.IP = input.IP
	}
	if input.WorkerID >= 0 {
		n.WorkerID = input.WorkerID
	}
	if input.DatacenterID >= 0 {
		n.DatacenterID = input.DatacenterID
	}
	n.UpdatedAt = time.Now()
	if err := n.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateMachineNode(n); err != nil {
		return nil, err
	}
	s.log.Infof("更新机器节点: %s", id)
	return n, nil
}

func (s *Service) DeleteMachineNode(id string) error {
	if _, err := s.store.GetMachineNode(id); err != nil {
		return err
	}
	// 检查是否有关联的租约
	for _, l := range s.store.ListLeases() {
		if l.NodeID == id {
			return model.NewValidationError("node", "该节点存在租约，不能删除")
		}
	}
	// 检查是否有关联的心跳
	for _, h := range s.store.ListNodeHeartbeats() {
		if h.NodeID == id {
			return model.NewValidationError("node", "该节点存在心跳记录，不能删除")
		}
	}
	s.log.Infof("删除机器节点: %s", id)
	return s.store.DeleteMachineNode(id)
}

func (s *Service) CountActiveMachineNodes() int {
	count := 0
	for _, n := range s.store.ListMachineNodes() {
		if n.Status == model.NodeStatusActive {
			count++
		}
	}
	return count
}

func (s *Service) ListMachineNodesByDatacenter(dcID int64) []*model.MachineNode {
	list := make([]*model.MachineNode, 0)
	for _, n := range s.store.ListMachineNodes() {
		if n.DatacenterID == dcID {
			list = append(list, n)
		}
	}
	return list
}
