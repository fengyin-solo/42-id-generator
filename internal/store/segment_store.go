package store

import "idgenerator/internal/model"

func (s *MemoryStore) CreateSegment(seg *model.Segment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.segments[seg.ID] = seg
	return nil
}

func (s *MemoryStore) GetSegment(id string) (*model.Segment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	seg, ok := s.segments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return seg, nil
}

func (s *MemoryStore) ListSegments() []*model.Segment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Segment, 0, len(s.segments))
	for _, seg := range s.segments {
		list = append(list, seg)
	}
	return list
}

func (s *MemoryStore) ListSegmentsByBizTypeID(bizTypeID string) []*model.Segment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Segment, 0)
	for _, seg := range s.segments {
		if seg.BizTypeID == bizTypeID {
			list = append(list, seg)
		}
	}
	return list
}

func (s *MemoryStore) GetActiveSegmentByBizTypeID(bizTypeID string) (*model.Segment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, seg := range s.segments {
		if seg.BizTypeID == bizTypeID && seg.Status == model.SegmentStatusUsing {
			return seg, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) UpdateSegment(seg *model.Segment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.segments[seg.ID]; !ok {
		return ErrNotFound
	}
	s.segments[seg.ID] = seg
	return nil
}

func (s *MemoryStore) DeleteSegment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.segments[id]; !ok {
		return ErrNotFound
	}
	delete(s.segments, id)
	return nil
}
