package handler

import (
	"net/http"
	"strconv"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerIDRuleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/id-rules", s.createIDRule)
	mux.HandleFunc("GET /api/id-rules", s.listIDRules)
	mux.HandleFunc("GET /api/id-rules/{id}", s.getIDRule)
	mux.HandleFunc("PUT /api/id-rules/{id}", s.updateIDRule)
	mux.HandleFunc("DELETE /api/id-rules/{id}", s.deleteIDRule)
	mux.HandleFunc("POST /api/id-rules/{id}/toggle", s.toggleIDRuleEnabled)
}

type createIDRuleRequest struct {
	Name           string `json:"name"`
	BizTypeID      string `json:"biz_type_id"`
	Mode           string `json:"mode"`
	SignBits       int    `json:"sign_bits"`
	TimestampBits  int    `json:"timestamp_bits"`
	DatacenterBits int    `json:"datacenter_bits"`
	WorkerBits     int    `json:"worker_bits"`
	SequenceBits   int    `json:"sequence_bits"`
}

func (s *Server) createIDRule(w http.ResponseWriter, r *http.Request) {
	var req createIDRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ir, err := s.svc.CreateIDRule(model.IDRule{
		Name:           req.Name,
		BizTypeID:      req.BizTypeID,
		Mode:           req.Mode,
		SignBits:       req.SignBits,
		TimestampBits:  req.TimestampBits,
		DatacenterBits: req.DatacenterBits,
		WorkerBits:     req.WorkerBits,
		SequenceBits:   req.SequenceBits,
		Enabled:        true,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, ir)
}

func (s *Server) listIDRules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.IDRuleFilter{
		BizTypeID: r.URL.Query().Get("biz_type_id"),
		Mode:      r.URL.Query().Get("mode"),
		Status:    r.URL.Query().Get("status"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	if enabledStr := r.URL.Query().Get("enabled"); enabledStr != "" {
		v, _ := strconv.ParseBool(enabledStr)
		filter.Enabled = &v
	}
	items, total, err := s.svc.ListIDRules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getIDRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ir, err := s.svc.GetIDRule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ir)
}

type updateIDRuleRequest struct {
	Name           string `json:"name"`
	Mode           string `json:"mode"`
	SignBits       int    `json:"sign_bits"`
	TimestampBits  int    `json:"timestamp_bits"`
	DatacenterBits int    `json:"datacenter_bits"`
	WorkerBits     int    `json:"worker_bits"`
	SequenceBits   int    `json:"sequence_bits"`
	Enabled        bool   `json:"enabled"`
	Status         string `json:"status"`
}

func (s *Server) updateIDRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateIDRuleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	ir, err := s.svc.UpdateIDRule(id, model.IDRule{
		Name:           req.Name,
		Mode:           req.Mode,
		SignBits:       req.SignBits,
		TimestampBits:  req.TimestampBits,
		DatacenterBits: req.DatacenterBits,
		WorkerBits:     req.WorkerBits,
		SequenceBits:   req.SequenceBits,
		Enabled:        req.Enabled,
		Status:         req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ir)
}

func (s *Server) deleteIDRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteIDRule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) toggleIDRuleEnabled(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ir, err := s.svc.ToggleIDRuleEnabled(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, ir)
}
