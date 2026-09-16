package handler

import (
	"net/http"

	"idgenerator/internal/model"
	"idgenerator/pkg/httpx"
	"idgenerator/pkg/idgen"
)

func (s *Server) registerIDGenRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/id-gen/snowflake", s.generateSnowflakeID)
	mux.HandleFunc("POST /api/id-gen/segment", s.generateSegmentID)
	mux.HandleFunc("POST /api/id-gen/sequence", s.generateSequenceID)
	mux.HandleFunc("POST /api/id-gen/batch", s.generateBatchIDs)
}

type generateSnowflakeRequest struct {
	BizTypeID    string `json:"biz_type_id"`
	NodeID       string `json:"node_id"`
}

func (s *Server) generateSnowflakeID(w http.ResponseWriter, r *http.Request) {
	var req generateSnowflakeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.BizTypeID == "" {
		httpx.BadRequest(w, "biz_type_id 不能为空")
		return
	}
	// 校验业务类型
	bt, err := s.svc.GetBizType(req.BizTypeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if bt.Mode != model.BizModeSnowflake {
		httpx.BadRequest(w, "该业务类型不是雪花模式")
		return
	}
	// 获取雪花配置
	sfConfig, err := s.svc.GetSnowflakeConfigByBizTypeID(req.BizTypeID)
	if err != nil {
		httpx.BadRequest(w, "该业务类型未配置雪花算法: "+err.Error())
		return
	}
	// 获取节点 worker_id
	var workerID, datacenterID int64
	if req.NodeID != "" {
		node, err := s.svc.GetMachineNodeByNodeID(req.NodeID)
		if err == nil {
			workerID = node.WorkerID
			datacenterID = node.DatacenterID
		}
	}
	sf, err := idgen.NewSnowflake(datacenterID, workerID,
		uint8(sfConfig.DatacenterBits),
		uint8(sfConfig.WorkerBits),
		uint8(sfConfig.SequenceBits),
		sfConfig.Twepoch)
	if err != nil {
		httpx.BadRequest(w, "雪花算法初始化失败: "+err.Error())
		return
	}
	idVal, err := sf.NextID()
	if err != nil {
		httpx.InternalError(w, "生成 ID 失败: "+err.Error())
		return
	}
	// 记录分配
	_, _ = s.svc.CreateAllocRecord(model.AllocRecord{
		BizTypeID: req.BizTypeID,
		NodeID:    req.NodeID,
		BatchSize: 1,
		Mode:      model.BizModeSnowflake,
	})
	httpx.OK(w, map[string]interface{}{
		"id":         idVal,
		"biz_type":   bt.Code,
		"mode":       model.BizModeSnowflake,
		"worker_id":  workerID,
	})
}

type generateSegmentRequest struct {
	BizTypeID string `json:"biz_type_id"`
	NodeID    string `json:"node_id"`
}

func (s *Server) generateSegmentID(w http.ResponseWriter, r *http.Request) {
	var req generateSegmentRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.BizTypeID == "" {
		httpx.BadRequest(w, "biz_type_id 不能为空")
		return
	}
	bt, err := s.svc.GetBizType(req.BizTypeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if bt.Mode != model.BizModeSegment {
		httpx.BadRequest(w, "该业务类型不是号段模式")
		return
	}
	// 获取活跃号段
	seg, err := s.svc.SwitchSegment(req.BizTypeID)
	if err != nil {
		httpx.BadRequest(w, "无可用号段: "+err.Error())
		return
	}
	// 推进游标
	seg, ids, err := s.svc.AdvanceSegmentCursor(seg.ID, 1)
	if err != nil {
		httpx.InternalError(w, "取号失败: "+err.Error())
		return
	}
	if len(ids) == 0 {
		httpx.BadRequest(w, "号段已耗尽")
		return
	}
	// 记录分配
	_, _ = s.svc.CreateAllocRecord(model.AllocRecord{
		BizTypeID: req.BizTypeID,
		NodeID:    req.NodeID,
		SegmentID: seg.ID,
		BatchSize: 1,
		StartID:   ids[0],
		EndID:     ids[0],
		Mode:      model.BizModeSegment,
	})
	httpx.OK(w, map[string]interface{}{
		"id":        ids[0],
		"biz_type":  bt.Code,
		"mode":      model.BizModeSegment,
		"segment_id": seg.ID,
	})
}

type generateSequenceRequest struct {
	BizTypeID string `json:"biz_type_id"`
	NodeID    string `json:"node_id"`
}

func (s *Server) generateSequenceID(w http.ResponseWriter, r *http.Request) {
	var req generateSequenceRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.BizTypeID == "" {
		httpx.BadRequest(w, "biz_type_id 不能为空")
		return
	}
	bt, err := s.svc.GetBizType(req.BizTypeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if bt.Mode != model.BizModeSequence {
		httpx.BadRequest(w, "该业务类型不是序列模式")
		return
	}
	// 使用内存中的序列生成器（这里简化，直接用自增逻辑）
	// 实际生产中需要持久化序列号
	seq := idgen.NewSequenceGenerator(1000000)
	idVal := seq.Next()
	_, _ = s.svc.CreateAllocRecord(model.AllocRecord{
		BizTypeID: req.BizTypeID,
		NodeID:    req.NodeID,
		BatchSize: 1,
		StartID:   idVal,
		EndID:     idVal,
		Mode:      model.BizModeSequence,
	})
	httpx.OK(w, map[string]interface{}{
		"id":       idVal,
		"biz_type": bt.Code,
		"mode":     model.BizModeSequence,
	})
}

type generateBatchRequest struct {
	BizTypeID string `json:"biz_type_id"`
	NodeID    string `json:"node_id"`
	BatchSize int    `json:"batch_size"`
}

func (s *Server) generateBatchIDs(w http.ResponseWriter, r *http.Request) {
	var req generateBatchRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if req.BizTypeID == "" {
		httpx.BadRequest(w, "biz_type_id 不能为空")
		return
	}
	if req.BatchSize <= 0 || req.BatchSize > 1000 {
		httpx.BadRequest(w, "batch_size 必须在 1-1000 之间")
		return
	}
	bt, err := s.svc.GetBizType(req.BizTypeID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	ids := make([]int64, 0, req.BatchSize)
	switch bt.Mode {
	case model.BizModeSnowflake:
		sfConfig, err := s.svc.GetSnowflakeConfigByBizTypeID(req.BizTypeID)
		if err != nil {
			httpx.BadRequest(w, "该业务类型未配置雪花算法: "+err.Error())
			return
		}
		var workerID, datacenterID int64
		if req.NodeID != "" {
			node, err := s.svc.GetMachineNodeByNodeID(req.NodeID)
			if err == nil {
				workerID = node.WorkerID
				datacenterID = node.DatacenterID
			}
		}
		sf, err := idgen.NewSnowflake(datacenterID, workerID,
			uint8(sfConfig.DatacenterBits),
			uint8(sfConfig.WorkerBits),
			uint8(sfConfig.SequenceBits),
			sfConfig.Twepoch)
		if err != nil {
			httpx.BadRequest(w, "雪花算法初始化失败: "+err.Error())
			return
		}
		for i := 0; i < req.BatchSize; i++ {
			idVal, err := sf.NextID()
			if err != nil {
				break
			}
			ids = append(ids, idVal)
		}
	case model.BizModeSegment:
		seg, err := s.svc.SwitchSegment(req.BizTypeID)
		if err != nil {
			httpx.BadRequest(w, "无可用号段: "+err.Error())
			return
		}
		seg, batchIDs, err := s.svc.AdvanceSegmentCursor(seg.ID, req.BatchSize)
		if err != nil {
			httpx.InternalError(w, "取号失败: "+err.Error())
			return
		}
		ids = append(ids, batchIDs...)
		if seg != nil && len(ids) > 0 {
			_, _ = s.svc.CreateAllocRecord(model.AllocRecord{
				BizTypeID: req.BizTypeID,
				NodeID:    req.NodeID,
				SegmentID: seg.ID,
				BatchSize: int64(len(ids)),
				StartID:   ids[0],
				EndID:     ids[len(ids)-1],
				Mode:      model.BizModeSegment,
			})
		}
	case model.BizModeSequence:
		seq := idgen.NewSequenceGenerator(1000000)
		for i := 0; i < req.BatchSize; i++ {
			ids = append(ids, seq.Next())
		}
	default:
		httpx.BadRequest(w, "不支持的发号模式")
		return
	}
	if len(ids) == 0 {
		httpx.BadRequest(w, "未生成任何 ID")
		return
	}
	httpx.OK(w, map[string]interface{}{
		"ids":       ids,
		"count":     len(ids),
		"biz_type":  bt.Code,
		"mode":      bt.Mode,
	})
}
