package handler

import (
	"net/http"
	"strconv"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerMachineNodeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/nodes", s.createMachineNode)
	mux.HandleFunc("GET /api/nodes", s.listMachineNodes)
	mux.HandleFunc("GET /api/nodes/{id}", s.getMachineNode)
	mux.HandleFunc("PUT /api/nodes/{id}", s.updateMachineNode)
	mux.HandleFunc("DELETE /api/nodes/{id}", s.deleteMachineNode)
	mux.HandleFunc("GET /api/nodes/{id}/health", s.checkNodeHealth)
}

type createMachineNodeRequest struct {
	NodeID       string `json:"node_id"`
	Hostname     string `json:"hostname"`
	IP           string `json:"ip"`
	WorkerID     int64  `json:"worker_id"`
	DatacenterID int64  `json:"datacenter_id"`
	Status       string `json:"status"`
}

func (s *Server) createMachineNode(w http.ResponseWriter, r *http.Request) {
	var req createMachineNodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateMachineNode(model.MachineNode{
		NodeID:       req.NodeID,
		Hostname:     req.Hostname,
		IP:           req.IP,
		WorkerID:     req.WorkerID,
		DatacenterID: req.DatacenterID,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listMachineNodes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MachineNodeFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	if dcStr := r.URL.Query().Get("datacenter_id"); dcStr != "" {
		if dc, err := strconv.ParseInt(dcStr, 10, 64); err == nil {
			filter.DatacenterID = &dc
		}
	}
	items, total, err := s.svc.ListMachineNodes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMachineNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetMachineNode(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

type updateMachineNodeRequest struct {
	Hostname     string `json:"hostname"`
	IP           string `json:"ip"`
	WorkerID     int64  `json:"worker_id"`
	DatacenterID int64  `json:"datacenter_id"`
	Status       string `json:"status"`
}

func (s *Server) updateMachineNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateMachineNodeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.UpdateMachineNode(id, model.MachineNode{
		Hostname:     req.Hostname,
		IP:           req.IP,
		WorkerID:     req.WorkerID,
		DatacenterID: req.DatacenterID,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteMachineNode(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteMachineNode(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) checkNodeHealth(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	node, err := s.svc.GetMachineNode(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	healthy, lastBeat, err := s.svc.CheckNodeHealth(node.NodeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{
		"node_id":    node.NodeID,
		"healthy":    healthy,
		"last_beat":  lastBeat,
	})
}
