package handler

import (
	"net/http"
	"strconv"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerAllocStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/alloc-stats", s.createAllocStats)
	mux.HandleFunc("GET /api/alloc-stats", s.listAllocStats)
	mux.HandleFunc("GET /api/alloc-stats/{id}", s.getAllocStats)
	mux.HandleFunc("PUT /api/alloc-stats/{id}", s.updateAllocStats)
	mux.HandleFunc("DELETE /api/alloc-stats/{id}", s.deleteAllocStats)
	mux.HandleFunc("GET /api/stats/overview", s.getStatsOverview)
	mux.HandleFunc("GET /api/stats/top", s.getTopAllocStats)
	mux.HandleFunc("GET /api/stats/group-by-biz", s.getStatsGroupByBizType)
	mux.HandleFunc("GET /api/stats/group-by-node", s.getStatsGroupByNode)
}

type createAllocStatsRequest struct {
	BizTypeID      string  `json:"biz_type_id"`
	NodeID         string  `json:"node_id"`
	Date           string  `json:"date"`
	TotalAllocated int64   `json:"total_allocated"`
	PeakQPS        float64 `json:"peak_qps"`
	AvgBatchSize   float64 `json:"avg_batch_size"`
}

func (s *Server) createAllocStats(w http.ResponseWriter, r *http.Request) {
	var req createAllocStatsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.CreateAllocStats(model.AllocStats{
		BizTypeID:      req.BizTypeID,
		NodeID:         req.NodeID,
		Date:           req.Date,
		TotalAllocated: req.TotalAllocated,
		PeakQPS:        req.PeakQPS,
		AvgBatchSize:   req.AvgBatchSize,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, st)
}

func (s *Server) listAllocStats(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.AllocStatsFilter{
		BizTypeID: r.URL.Query().Get("biz_type_id"),
		NodeID:    r.URL.Query().Get("node_id"),
		Date:      r.URL.Query().Get("date"),
		DateFrom:  r.URL.Query().Get("date_from"),
		DateTo:    r.URL.Query().Get("date_to"),
	}
	items, total, err := s.svc.ListAllocStats(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getAllocStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, err := s.svc.GetAllocStats(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

type updateAllocStatsRequest struct {
	TotalAllocated int64   `json:"total_allocated"`
	PeakQPS        float64 `json:"peak_qps"`
	AvgBatchSize   float64 `json:"avg_batch_size"`
}

func (s *Server) updateAllocStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateAllocStatsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.UpdateAllocStats(id, model.AllocStats{
		TotalAllocated: req.TotalAllocated,
		PeakQPS:        req.PeakQPS,
		AvgBatchSize:   req.AvgBatchSize,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) deleteAllocStats(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteAllocStats(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getStatsOverview(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.GetStatsOverview())
}

func (s *Server) getTopAllocStats(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		httpx.BadRequest(w, "date 参数必填")
		return
	}
	topN := 10
	if nStr := r.URL.Query().Get("top_n"); nStr != "" {
		if n, err := strconv.Atoi(nStr); err == nil && n > 0 {
			topN = n
		}
	}
	httpx.OK(w, s.svc.GetTopAllocStatsByDate(date, topN))
}

func (s *Server) getStatsGroupByBizType(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.GetStatsGroupByBizType())
}

func (s *Server) getStatsGroupByNode(w http.ResponseWriter, r *http.Request) {
	httpx.OK(w, s.svc.GetStatsGroupByNode())
}
