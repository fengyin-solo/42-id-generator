package handler

import (
	"net/http"
	"strconv"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerBizTypeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/biz-types", s.createBizType)
	mux.HandleFunc("GET /api/biz-types", s.listBizTypes)
	mux.HandleFunc("GET /api/biz-types/{id}", s.getBizType)
	mux.HandleFunc("PUT /api/biz-types/{id}", s.updateBizType)
	mux.HandleFunc("DELETE /api/biz-types/{id}", s.deleteBizType)
	mux.HandleFunc("POST /api/biz-types/{id}/toggle", s.toggleBizTypeStatus)
}

type createBizTypeRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Mode        string `json:"mode"`
	SegmentStep int    `json:"segment_step"`
	Description string `json:"description"`
}

func (s *Server) createBizType(w http.ResponseWriter, r *http.Request) {
	var req createBizTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBizType(model.BizType{
		Name:        req.Name,
		Code:        req.Code,
		Mode:        req.Mode,
		SegmentStep: req.SegmentStep,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBizTypes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	enabledStr := r.URL.Query().Get("enabled")
	var enabled *bool
	if enabledStr != "" {
		v, _ := strconv.ParseBool(enabledStr)
		enabled = &v
	}
	filter := model.BizTypeFilter{
		Mode:    r.URL.Query().Get("mode"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
		Enabled: enabled,
	}
	items, total, err := s.svc.ListBizTypes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBizType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.GetBizType(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type updateBizTypeRequest struct {
	Name        string `json:"name"`
	SegmentStep int    `json:"segment_step"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Status      string `json:"status"`
}

func (s *Server) updateBizType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateBizTypeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateBizType(id, model.BizType{
		Name:        req.Name,
		SegmentStep: req.SegmentStep,
		Description: req.Description,
		Enabled:     req.Enabled,
		Status:      req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBizType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBizType(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) toggleBizTypeStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.ToggleBizTypeStatus(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}
