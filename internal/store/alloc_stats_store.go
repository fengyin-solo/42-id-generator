package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateAllocStats(st *model.AllocStats) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := st.BizTypeID + ":" + st.NodeID + ":" + st.Date
	s.allocStats[key] = st
	return nil
}

func (s *MemoryStore) GetAllocStats(id string) (*model.AllocStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, st := range s.allocStats {
		if st.ID == id {
			return st, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListAllocStats() []*model.AllocStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AllocStats, 0, len(s.allocStats))
	for _, st := range s.allocStats {
		list = append(list, st)
	}
	return list
}

func (s *MemoryStore) GetAllocStatsByBizNodeDate(bizTypeID, nodeID, date string) (*model.AllocStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := bizTypeID + ":" + nodeID + ":" + date
	st, ok := s.allocStats[key]
	if !ok {
		return nil, ErrNotFound
	}
	return st, nil
}

func (s *MemoryStore) UpdateAllocStats(st *model.AllocStats) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := st.BizTypeID + ":" + st.NodeID + ":" + st.Date
	old, ok := s.allocStats[key]
	if !ok {
		for _, exist := range s.allocStats {
			if exist.ID == st.ID {
				old = exist
				break
			}
		}
	}
	if old == nil {
		return ErrNotFound
	}
	s.allocStats[key] = st
	return nil
}

func (s *MemoryStore) DeleteAllocStats(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, st := range s.allocStats {
		if st.ID == id {
			delete(s.allocStats, key)
			return nil
		}
	}
	return ErrNotFound
}
