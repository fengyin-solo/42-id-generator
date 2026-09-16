package handler

import (
	"net/http"
	"strconv"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerSegmentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/segments", s.createSegment)
	mux.HandleFunc("GET /api/segments", s.listSegments)
	mux.HandleFunc("GET /api/segments/{id}", s.getSegment)
	mux.HandleFunc("PUT /api/segments/{id}", s.updateSegment)
	mux.HandleFunc("DELETE /api/segments/{id}", s.deleteSegment)
	mux.HandleFunc("POST /api/segments/{id}/advance", s.advanceSegmentCursor)
}

type createSegmentRequest struct {
	BizTypeID string `json:"biz_type_id"`
	StartID   int64  `json:"start_id"`
	EndID     int64  `json:"end_id"`
	NodeID    string `json:"node_id"`
}

func (s *Server) createSegment(w http.ResponseWriter, r *http.Request) {
	var req createSegmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	seg, err := s.svc.CreateSegment(model.Segment{
		BizTypeID: req.BizTypeID,
		StartID:   req.StartID,
		EndID:     req.EndID,
		Cursor:    req.StartID,
		Status:    model.SegmentStatusUsing,
		NodeID:    req.NodeID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, seg)
}

func (s *Server) listSegments(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SegmentFilter{
		BizTypeID: r.URL.Query().Get("biz_type_id"),
		Status:    r.URL.Query().Get("status"),
		NodeID:    r.URL.Query().Get("node_id"),
	}
	if exhaustedStr := r.URL.Query().Get("exhausted"); exhaustedStr != "" {
		v, _ := strconv.ParseBool(exhaustedStr)
		filter.Exhausted = &v
	}
	items, total, err := s.svc.ListSegments(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSegment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	seg, err := s.svc.GetSegment(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, seg)
}

type updateSegmentRequest struct {
	Cursor int64  `json:"cursor"`
	Status string `json:"status"`
	NodeID string `json:"node_id"`
}

func (s *Server) updateSegment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSegmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	seg, err := s.svc.UpdateSegment(id, model.Segment{
		Cursor: req.Cursor,
		Status: req.Status,
		NodeID: req.NodeID,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, seg)
}

func (s *Server) deleteSegment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSegment(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type advanceSegmentRequest struct {
	BatchSize int `json:"batch_size"`
}

func (s *Server) advanceSegmentCursor(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req advanceSegmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 1
	}
	seg, ids, err := s.svc.AdvanceSegmentCursor(id, req.BatchSize)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{
		"segment": seg,
		"ids":     ids,
	})
}
