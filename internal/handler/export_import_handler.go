package handler

import (
	"net/http"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
)

func (s *Server) registerExportImportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export", s.exportSnapshot)
	mux.HandleFunc("POST /api/import", s.importSnapshot)
}

type snapshotData struct {
	BizTypes         []*model.BizType         `json:"biz_types"`
	MachineNodes     []*model.MachineNode     `json:"machine_nodes"`
	IDRules          []*model.IDRule          `json:"id_rules"`
	Segments         []*model.Segment         `json:"segments"`
	AllocRecords     []*model.AllocRecord     `json:"alloc_records"`
	Leases           []*model.Lease           `json:"leases"`
	NodeHeartbeats   []*model.NodeHeartbeat   `json:"node_heartbeats"`
	SnowflakeConfigs []*model.SnowflakeConfig `json:"snowflake_configs"`
	AllocStats       []*model.AllocStats      `json:"alloc_stats"`
	RecycleRecords   []*model.RecycleRecord   `json:"recycle_records"`
}

func (s *Server) exportSnapshot(w http.ResponseWriter, r *http.Request) {
	bizTypes, _, _ := s.svc.ListBizTypes(model.BizTypeFilter{}, 1, 99999)
	machineNodes, _, _ := s.svc.ListMachineNodes(model.MachineNodeFilter{}, 1, 99999)
	idRules, _, _ := s.svc.ListIDRules(model.IDRuleFilter{}, 1, 99999)
	segments, _, _ := s.svc.ListSegments(model.SegmentFilter{}, 1, 99999)
	allocRecords, _, _ := s.svc.ListAllocRecords(model.AllocRecordFilter{}, 1, 99999)
	leases, _, _ := s.svc.ListLeases(model.LeaseFilter{}, 1, 99999)
	nodeHeartbeats, _, _ := s.svc.ListNodeHeartbeats(model.NodeHeartbeatFilter{}, 1, 99999)
	snowflakeConfigs, _, _ := s.svc.ListSnowflakeConfigs(model.SnowflakeConfigFilter{}, 1, 99999)
	allocStats, _, _ := s.svc.ListAllocStats(model.AllocStatsFilter{}, 1, 99999)
	recycleRecords, _, _ := s.svc.ListRecycleRecords(model.RecycleRecordFilter{}, 1, 99999)
	data := snapshotData{
		BizTypes:         bizTypes,
		MachineNodes:     machineNodes,
		IDRules:          idRules,
		Segments:         segments,
		AllocRecords:     allocRecords,
		Leases:           leases,
		NodeHeartbeats:   nodeHeartbeats,
		SnowflakeConfigs: snowflakeConfigs,
		AllocStats:       allocStats,
		RecycleRecords:   recycleRecords,
	}
	httpx.OK(w, data)
}

func (s *Server) importSnapshot(w http.ResponseWriter, r *http.Request) {
	var data snapshotData
	if err := httpx.Decode(r, &data); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	imported := 0
	for _, b := range data.BizTypes {
		if _, err := s.svc.GetBizType(b.ID); err != nil {
			_ = s.store.CreateBizType(b)
			imported++
		}
	}
	for _, n := range data.MachineNodes {
		if _, err := s.svc.GetMachineNode(n.ID); err != nil {
			_ = s.store.CreateMachineNode(n)
			imported++
		}
	}
	for _, ir := range data.IDRules {
		if _, err := s.svc.GetIDRule(ir.ID); err != nil {
			_ = s.store.CreateIDRule(ir)
			imported++
		}
	}
	for _, seg := range data.Segments {
		if _, err := s.svc.GetSegment(seg.ID); err != nil {
			_ = s.store.CreateSegment(seg)
			imported++
		}
	}
	for _, ar := range data.AllocRecords {
		if _, err := s.svc.GetAllocRecord(ar.ID); err != nil {
			_ = s.store.CreateAllocRecord(ar)
			imported++
		}
	}
	for _, l := range data.Leases {
		if _, err := s.svc.GetLease(l.ID); err != nil {
			_ = s.store.CreateLease(l)
			imported++
		}
	}
	for _, hb := range data.NodeHeartbeats {
		if _, err := s.svc.GetNodeHeartbeat(hb.ID); err != nil {
			_ = s.store.CreateNodeHeartbeat(hb)
			imported++
		}
	}
	for _, sc := range data.SnowflakeConfigs {
		if _, err := s.svc.GetSnowflakeConfig(sc.ID); err != nil {
			_ = s.store.CreateSnowflakeConfig(sc)
			imported++
		}
	}
	for _, st := range data.AllocStats {
		if _, err := s.svc.GetAllocStats(st.ID); err != nil {
			_ = s.store.CreateAllocStats(st)
			imported++
		}
	}
	for _, rec := range data.RecycleRecords {
		if _, err := s.svc.GetRecycleRecord(rec.ID); err != nil {
			_ = s.store.CreateRecycleRecord(rec)
			imported++
		}
	}
	httpx.OK(w, map[string]interface{}{
		"imported": imported,
	})
}
