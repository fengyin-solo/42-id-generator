package handler

import (
	"net/http"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerRecycleRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/recycle-records", s.createRecycleRecord)
	mux.HandleFunc("GET /api/recycle-records", s.listRecycleRecords)
	mux.HandleFunc("GET /api/recycle-records/{id}", s.getRecycleRecord)
	mux.HandleFunc("PUT /api/recycle-records/{id}", s.updateRecycleRecord)
	mux.HandleFunc("DELETE /api/recycle-records/{id}", s.deleteRecycleRecord)
}

type createRecycleRecordRequest struct {
	SegmentID string `json:"segment_id"`
	BizTypeID string `json:"biz_type_id"`
	Reason    string `json:"reason"`
}

func (s *Server) createRecycleRecord(w http.ResponseWriter, r *http.Request) {
	var req createRecycleRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.CreateRecycleRecord(model.RecycleRecord{
		SegmentID: req.SegmentID,
		BizTypeID: req.BizTypeID,
		Reason:    req.Reason,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) listRecycleRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RecycleRecordFilter{
		SegmentID: r.URL.Query().Get("segment_id"),
		BizTypeID: r.URL.Query().Get("biz_type_id"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRecycleRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRecycleRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rec, err := s.svc.GetRecycleRecord(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

type updateRecycleRecordRequest struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

func (s *Server) updateRecycleRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRecycleRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.UpdateRecycleRecord(id, model.RecycleRecord{
		Status: req.Status,
		Reason: req.Reason,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) deleteRecycleRecord(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRecycleRecord(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
