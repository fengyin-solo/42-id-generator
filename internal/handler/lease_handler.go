package handler

import (
	"net/http"
	"strconv"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerLeaseRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/leases", s.createLease)
	mux.HandleFunc("GET /api/leases", s.listLeases)
	mux.HandleFunc("GET /api/leases/{id}", s.getLease)
	mux.HandleFunc("POST /api/leases/{id}/renew", s.renewLease)
	mux.HandleFunc("POST /api/leases/{id}/expire", s.expireLease)
	mux.HandleFunc("DELETE /api/leases/{id}", s.deleteLease)
}

type createLeaseRequest struct {
	NodeID    string `json:"node_id"`
	BizTypeID string `json:"biz_type_id"`
	ExpiresAt string `json:"expires_at"`
}

func (s *Server) createLease(w http.ResponseWriter, r *http.Request) {
	var req createLeaseRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	expiresAt, err := time.Parse(time.RFC3339, req.ExpiresAt)
	if err != nil {
		httpx.BadRequest(w, "过期时间格式错误: "+err.Error())
		return
	}
	l, err := s.svc.CreateLease(model.Lease{
		NodeID:    req.NodeID,
		BizTypeID: req.BizTypeID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, l)
}

func (s *Server) listLeases(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.LeaseFilter{
		NodeID:    r.URL.Query().Get("node_id"),
		BizTypeID: r.URL.Query().Get("biz_type_id"),
		Status:    r.URL.Query().Get("status"),
	}
	if expiredStr := r.URL.Query().Get("expired"); expiredStr != "" {
		v, _ := strconv.ParseBool(expiredStr)
		filter.Expired = &v
	}
	items, total, err := s.svc.ListLeases(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getLease(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	l, err := s.svc.GetLease(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

type renewLeaseRequest struct {
	DurationSeconds int `json:"duration_seconds"`
}

func (s *Server) renewLease(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req renewLeaseRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.DurationSeconds <= 0 {
		req.DurationSeconds = 60
	}
	l, err := s.svc.RenewLease(id, time.Duration(req.DurationSeconds)*time.Second)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

func (s *Server) expireLease(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	l, err := s.svc.ExpireLease(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

func (s *Server) deleteLease(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteLease(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
