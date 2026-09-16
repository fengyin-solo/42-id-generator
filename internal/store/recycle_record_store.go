package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateRecycleRecord(r *model.RecycleRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recycleRecords[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRecycleRecord(id string) (*model.RecycleRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.recycleRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRecycleRecords() []*model.RecycleRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RecycleRecord, 0, len(s.recycleRecords))
	for _, r := range s.recycleRecords {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRecycleRecord(r *model.RecycleRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recycleRecords[r.ID]; !ok {
		return ErrNotFound
	}
	s.recycleRecords[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRecycleRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recycleRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.recycleRecords, id)
	return nil
}
