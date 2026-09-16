package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateSnowflakeConfig(c *model.SnowflakeConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.snowflakeConfigs {
		if exist.BizTypeID == c.BizTypeID {
			return ErrConflict
		}
	}
	s.snowflakeConfigs[c.ID] = c
	return nil
}

func (s *MemoryStore) GetSnowflakeConfig(id string) (*model.SnowflakeConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.snowflakeConfigs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetSnowflakeConfigByBizTypeID(bizTypeID string) (*model.SnowflakeConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.snowflakeConfigs {
		if c.BizTypeID == bizTypeID {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListSnowflakeConfigs() []*model.SnowflakeConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SnowflakeConfig, 0, len(s.snowflakeConfigs))
	for _, c := range s.snowflakeConfigs {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateSnowflakeConfig(c *model.SnowflakeConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.snowflakeConfigs[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.snowflakeConfigs {
		if exist.ID != c.ID && exist.BizTypeID == c.BizTypeID {
			return ErrConflict
		}
	}
	s.snowflakeConfigs[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteSnowflakeConfig(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.snowflakeConfigs[id]; !ok {
		return ErrNotFound
	}
	delete(s.snowflakeConfigs, id)
	return nil
}
