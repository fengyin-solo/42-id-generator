package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateMachineNode(n *model.MachineNode) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.machineNodes {
		if exist.NodeID == n.NodeID {
			return ErrConflict
		}
	}
	s.machineNodes[n.ID] = n
	return nil
}

func (s *MemoryStore) GetMachineNode(id string) (*model.MachineNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.machineNodes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) GetMachineNodeByNodeID(nodeID string) (*model.MachineNode, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, n := range s.machineNodes {
		if n.NodeID == nodeID {
			return n, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListMachineNodes() []*model.MachineNode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.MachineNode, 0, len(s.machineNodes))
	for _, n := range s.machineNodes {
		list = append(list, n)
	}
	return list
}

func (s *MemoryStore) UpdateMachineNode(n *model.MachineNode) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.machineNodes[n.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.machineNodes {
		if exist.ID != n.ID && exist.NodeID == n.NodeID {
			return ErrConflict
		}
	}
	s.machineNodes[n.ID] = n
	return nil
}

func (s *MemoryStore) DeleteMachineNode(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.machineNodes[id]; !ok {
		return ErrNotFound
	}
	delete(s.machineNodes, id)
	return nil
}
