package handler

import (
	"net/http"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerSnowflakeConfigRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/snowflake-configs", s.createSnowflakeConfig)
	mux.HandleFunc("GET /api/snowflake-configs", s.listSnowflakeConfigs)
	mux.HandleFunc("GET /api/snowflake-configs/{id}", s.getSnowflakeConfig)
	mux.HandleFunc("PUT /api/snowflake-configs/{id}", s.updateSnowflakeConfig)
	mux.HandleFunc("DELETE /api/snowflake-configs/{id}", s.deleteSnowflakeConfig)
}

type createSnowflakeConfigRequest struct {
	BizTypeID      string `json:"biz_type_id"`
	EpochMs        int64  `json:"epoch_ms"`
	DatacenterBits int    `json:"datacenter_bits"`
	WorkerBits     int    `json:"worker_bits"`
	SequenceBits   int    `json:"sequence_bits"`
	Twepoch        int64  `json:"twepoch"`
}

func (s *Server) createSnowflakeConfig(w http.ResponseWriter, r *http.Request) {
	var req createSnowflakeConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cfg, err := s.svc.CreateSnowflakeConfig(model.SnowflakeConfig{
		BizTypeID:      req.BizTypeID,
		EpochMs:        req.EpochMs,
		DatacenterBits: req.DatacenterBits,
		WorkerBits:     req.WorkerBits,
		SequenceBits:   req.SequenceBits,
		Twepoch:        req.Twepoch,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, cfg)
}

func (s *Server) listSnowflakeConfigs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SnowflakeConfigFilter{
		BizTypeID: r.URL.Query().Get("biz_type_id"),
		Keyword:   r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSnowflakeConfigs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSnowflakeConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, err := s.svc.GetSnowflakeConfig(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cfg)
}

type updateSnowflakeConfigRequest struct {
	EpochMs        int64 `json:"epoch_ms"`
	DatacenterBits int   `json:"datacenter_bits"`
	WorkerBits     int   `json:"worker_bits"`
	SequenceBits   int   `json:"sequence_bits"`
	Twepoch        int64 `json:"twepoch"`
}

func (s *Server) updateSnowflakeConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSnowflakeConfigRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	cfg, err := s.svc.UpdateSnowflakeConfig(id, model.SnowflakeConfig{
		EpochMs:        req.EpochMs,
		DatacenterBits: req.DatacenterBits,
		WorkerBits:     req.WorkerBits,
		SequenceBits:   req.SequenceBits,
		Twepoch:        req.Twepoch,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, cfg)
}

func (s *Server) deleteSnowflakeConfig(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSnowflakeConfig(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
