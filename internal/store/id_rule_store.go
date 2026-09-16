package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateIDRule(r *model.IDRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.idRules {
		if exist.BizTypeID == r.BizTypeID {
			return ErrConflict
		}
	}
	s.idRules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetIDRule(id string) (*model.IDRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.idRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) GetIDRuleByBizTypeID(bizTypeID string) (*model.IDRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.idRules {
		if r.BizTypeID == bizTypeID {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListIDRules() []*model.IDRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.IDRule, 0, len(s.idRules))
	for _, r := range s.idRules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateIDRule(r *model.IDRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.idRules[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.idRules {
		if exist.ID != r.ID && exist.BizTypeID == r.BizTypeID {
			return ErrConflict
		}
	}
	s.idRules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteIDRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.idRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.idRules, id)
	return nil
}
