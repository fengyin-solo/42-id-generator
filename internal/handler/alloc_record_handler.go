package handler

import (
	"net/http"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerAllocRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alloc-records", s.createAllocRecord)
	mux.HandleFunc("GET /api/alloc-records", s.listAllocRecords)
	mux.HandleFunc("GET /api/alloc-records/{id}", s.getAllocRecord)
	mux.HandleFunc("DELETE /api/alloc-records/{id}", s.deleteAllocRecord)
}

type createAllocRecordRequest struct {
	BizTypeID string `json:"biz_type_id"`
	NodeID    string `json:"node_id"`
	SegmentID string `json:"segment_id"`
	BatchSize int64  `json:"batch_size"`
	StartID   int64  `json:"start_id"`
	EndID     int64  `json:"end_id"`
	Mode      string `json:"mode"`
}

func (s *Server) createAllocRecord(w http.ResponseWriter, r *http.Request) {
	var req createAllocRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ar, err := s.svc.CreateAllocRecord(model.AllocRecord{
		BizTypeID: req.BizTypeID,
		NodeID:    req.NodeID,
		SegmentID: req.SegmentID,
		BatchSize: req.BatchSize,
		StartID:   req.StartID,
		EndID:     req.EndID,
		Mode:      req.Mode,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, ar)
}

func (s *Server) listAllocRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AllocRecordFilter{
		BizTypeID: r.URL.Query().Get("biz_type_id"),
		NodeID:    r.URL.Query().Get("node_id"),
		SegmentID: r.URL.Query().Get("segment_id"),
		Mode:      r.URL.Query().Get("mode"),
	}
	items, total, err := s.svc.ListAllocRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAllocRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ar, err := s.svc.GetAllocRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ar)
}

func (s *Server) deleteAllocRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAllocRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
