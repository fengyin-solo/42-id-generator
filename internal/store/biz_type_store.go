package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateBizType(b *model.BizType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.bizTypes {
		if exist.Code == b.Code {
			return ErrConflict
		}
	}
	s.bizTypes[b.ID] = b
	return nil
}

func (s *MemoryStore) GetBizType(id string) (*model.BizType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.bizTypes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (s *MemoryStore) GetBizTypeByCode(code string) (*model.BizType, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, b := range s.bizTypes {
		if b.Code == code {
			return b, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListBizTypes() []*model.BizType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.BizType, 0, len(s.bizTypes))
	for _, b := range s.bizTypes {
		list = append(list, b)
	}
	return list
}

func (s *MemoryStore) UpdateBizType(b *model.BizType) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bizTypes[b.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.bizTypes {
		if exist.ID != b.ID && exist.Code == b.Code {
			return ErrConflict
		}
	}
	s.bizTypes[b.ID] = b
	return nil
}

func (s *MemoryStore) DeleteBizType(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.bizTypes[id]; !ok {
		return ErrNotFound
	}
	delete(s.bizTypes, id)
	return nil
}
