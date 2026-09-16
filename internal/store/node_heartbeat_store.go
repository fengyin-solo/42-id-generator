package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateNodeHeartbeat(h *model.NodeHeartbeat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodeHeartbeats[h.ID] = h
	return nil
}

func (s *MemoryStore) GetNodeHeartbeat(id string) (*model.NodeHeartbeat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.nodeHeartbeats[id]
	if !ok {
		return nil, ErrNotFound
	}
	return h, nil
}

func (s *MemoryStore) ListNodeHeartbeats() []*model.NodeHeartbeat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.NodeHeartbeat, 0, len(s.nodeHeartbeats))
	for _, h := range s.nodeHeartbeats {
		list = append(list, h)
	}
	return list
}

func (s *MemoryStore) GetLatestHeartbeatByNodeID(nodeID string) (*model.NodeHeartbeat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var latest *model.NodeHeartbeat
	for _, h := range s.nodeHeartbeats {
		if h.NodeID == nodeID {
			if latest == nil || h.BeatAt.After(latest.BeatAt) {
				latest = h
			}
		}
	}
	if latest == nil {
		return nil, ErrNotFound
	}
	return latest, nil
}

func (s *MemoryStore) UpdateNodeHeartbeat(h *model.NodeHeartbeat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodeHeartbeats[h.ID]; !ok {
		return ErrNotFound
	}
	s.nodeHeartbeats[h.ID] = h
	return nil
}

func (s *MemoryStore) DeleteNodeHeartbeat(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodeHeartbeats[id]; !ok {
		return ErrNotFound
	}
	delete(s.nodeHeartbeats, id)
	return nil
}
