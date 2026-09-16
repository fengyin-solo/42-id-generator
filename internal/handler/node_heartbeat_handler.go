package handler

import (
	"net/http"
	"time"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerNodeHeartbeatRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/heartbeats", s.createNodeHeartbeat)
	mux.HandleFunc("GET /api/heartbeats", s.listNodeHeartbeats)
	mux.HandleFunc("GET /api/heartbeats/{id}", s.getNodeHeartbeat)
	mux.HandleFunc("DELETE /api/heartbeats/{id}", s.deleteNodeHeartbeat)
	mux.HandleFunc("GET /api/nodes/{id}/latest-heartbeat", s.getLatestHeartbeat)
}

type createNodeHeartbeatRequest struct {
	NodeID      string  `json:"node_id"`
	Load        float64 `json:"load"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

func (s *Server) createNodeHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req createNodeHeartbeatRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	hb, err := s.svc.CreateNodeHeartbeat(model.NodeHeartbeat{
		NodeID:      req.NodeID,
		Load:        req.Load,
		CPUUsage:    req.CPUUsage,
		MemoryUsage: req.MemoryUsage,
		BeatAt:      time.Now(),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, hb)
}

func (s *Server) listNodeHeartbeats(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NodeHeartbeatFilter{
		NodeID: r.URL.Query().Get("node_id"),
	}
	items, total, err := s.svc.ListNodeHeartbeats(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNodeHeartbeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	hb, err := s.svc.GetNodeHeartbeat(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, hb)
}

func (s *Server) deleteNodeHeartbeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNodeHeartbeat(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) getLatestHeartbeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	node, err := s.svc.GetMachineNode(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	hb, err := s.svc.GetLatestHeartbeatByNodeID(node.NodeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, hb)
}
