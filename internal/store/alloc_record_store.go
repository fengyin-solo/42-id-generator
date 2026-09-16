package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateAllocRecord(a *model.AllocRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.allocRecords[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAllocRecord(id string) (*model.AllocRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.allocRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAllocRecords() []*model.AllocRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AllocRecord, 0, len(s.allocRecords))
	for _, a := range s.allocRecords {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateAllocRecord(a *model.AllocRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.allocRecords[a.ID]; !ok {
		return ErrNotFound
	}
	s.allocRecords[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteAllocRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.allocRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.allocRecords, id)
	return nil
}
