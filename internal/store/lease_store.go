package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateLease(l *model.Lease) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.leases[l.ID] = l
	return nil
}

func (s *MemoryStore) GetLease(id string) (*model.Lease, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.leases[id]
	if !ok {
		return nil, ErrNotFound
	}
	return l, nil
}

func (s *MemoryStore) ListLeases() []*model.Lease {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Lease, 0, len(s.leases))
	for _, l := range s.leases {
		list = append(list, l)
	}
	return list
}

func (s *MemoryStore) GetActiveLeaseByNodeAndBiz(nodeID, bizTypeID string) (*model.Lease, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, l := range s.leases {
		if l.NodeID == nodeID && l.BizTypeID == bizTypeID && l.Status == model.LeaseStatusActive {
			return l, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) UpdateLease(l *model.Lease) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.leases[l.ID]; !ok {
		return ErrNotFound
	}
	s.leases[l.ID] = l
	return nil
}

func (s *MemoryStore) DeleteLease(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.leases[id]; !ok {
		return ErrNotFound
	}
	delete(s.leases, id)
	return nil
}
