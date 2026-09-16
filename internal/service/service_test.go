package service

import (
	"testing"
	"time"

	"idgenerator/internal/config"
	"idgenerator/internal/model"
	"idgenerator/internal/store"
	"idgenerator/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestService_BizType(t *testing.T) {
	s := newTestService()
	b, err := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	if b.ID == "" {
		t.Fatalf("ID 不应为空")
	}
	got, err := s.GetBizType(b.ID)
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if got.Name != "订单" {
		t.Fatalf("名称不匹配")
	}
	_, err = s.GetBizTypeByCode("order")
	if err != nil {
		t.Fatalf("按 Code 查询失败: %v", err)
	}
	items, total, err := s.ListBizTypes(model.BizTypeFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("列表失败: %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("列表长度不匹配")
	}
	_, err = s.UpdateBizType(b.ID, model.BizType{Name: "订单系统"})
	if err != nil {
		t.Fatalf("更新失败: %v", err)
	}
	_, err = s.ToggleBizTypeStatus(b.ID)
	if err != nil {
		t.Fatalf("切换状态失败: %v", err)
	}
	if s.CountBizTypes() != 1 {
		t.Fatalf("计数不匹配")
	}
	if s.CountBizTypesByMode(model.BizModeSegment) != 1 {
		t.Fatalf("按模式计数不匹配")
	}
}

func TestService_MachineNode(t *testing.T) {
	s := newTestService()
	n, err := s.CreateMachineNode(model.MachineNode{NodeID: "node-1", Hostname: "host1", IP: "127.0.0.1", WorkerID: 1, DatacenterID: 1})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	got, err := s.GetMachineNode(n.ID)
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if got.Status != model.NodeStatusActive {
		t.Fatalf("默认状态应为 active")
	}
	// 状态机：active -> failed
	_, err = s.UpdateMachineNode(n.ID, model.MachineNode{Status: model.NodeStatusFailed})
	if err != nil {
		t.Fatalf("状态流转失败: %v", err)
	}
	// 状态机：failed -> active (不允许)
	_, err = s.UpdateMachineNode(n.ID, model.MachineNode{Status: model.NodeStatusActive})
	if err == nil {
		t.Fatalf("failed -> active 应被禁止")
	}
	if s.CountActiveMachineNodes() != 0 {
		t.Fatalf("活跃节点数不匹配")
	}
}

func TestService_IDRule(t *testing.T) {
	s := newTestService()
	b, _ := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSnowflake, SegmentStep: 0})
	_, err := s.CreateIDRule(model.IDRule{Name: "规则1", BizTypeID: b.ID, Mode: model.BizModeSnowflake, SignBits: 1, TimestampBits: 41, DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	// 外键校验失败
	_, err = s.CreateIDRule(model.IDRule{Name: "规则2", BizTypeID: "not-exist", Mode: model.BizModeSnowflake, SignBits: 1, TimestampBits: 41, DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12})
	if err == nil {
		t.Fatalf("外键校验应失败")
	}
}

func TestService_Segment(t *testing.T) {
	s := newTestService()
	b, _ := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100})
	seg, err := s.CreateSegment(model.Segment{BizTypeID: b.ID, StartID: 1, EndID: 1001})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	if seg.Status != model.SegmentStatusUsing {
		t.Fatalf("默认状态应为 using")
	}
	// 推进游标
	seg2, ids, err := s.AdvanceSegmentCursor(seg.ID, 5)
	if err != nil {
		t.Fatalf("推进失败: %v", err)
	}
	if len(ids) != 5 {
		t.Fatalf("应返回 5 个 ID")
	}
	if seg2.Cursor != 6 {
		t.Fatalf("游标应为 6")
	}
	// 状态机校验
	_, err = s.UpdateSegment(seg.ID, model.Segment{Status: model.SegmentStatusExhausted})
	if err != nil {
		t.Fatalf("状态流转失败: %v", err)
	}
	// using -> exhausted 后不能再改
	_, err = s.UpdateSegment(seg.ID, model.Segment{Status: model.SegmentStatusUsing})
	if err == nil {
		t.Fatalf("exhausted -> using 应被禁止")
	}
}

func TestService_Lease(t *testing.T) {
	s := newTestService()
	b, _ := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100})
	n, _ := s.CreateMachineNode(model.MachineNode{NodeID: "node-1", Hostname: "host1", IP: "127.0.0.1"})
	l, err := s.CreateLease(model.Lease{NodeID: n.NodeID, BizTypeID: b.ID, ExpiresAt: time.Now().Add(time.Hour)})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	// 重复创建应失败
	_, err = s.CreateLease(model.Lease{NodeID: n.NodeID, BizTypeID: b.ID, ExpiresAt: time.Now().Add(time.Hour)})
	if err == nil {
		t.Fatalf("重复租约应失败")
	}
	// 续期
	_, err = s.RenewLease(l.ID, 2*time.Hour)
	if err != nil {
		t.Fatalf("续期失败: %v", err)
	}
	// 过期
	_, err = s.ExpireLease(l.ID)
	if err != nil {
		t.Fatalf("过期失败: %v", err)
	}
	// 过期后不能续期
	_, err = s.RenewLease(l.ID, time.Hour)
	if err == nil {
		t.Fatalf("过期后续期应失败")
	}
}

func TestService_RecycleRecord(t *testing.T) {
	s := newTestService()
	b, _ := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100})
	seg, _ := s.CreateSegment(model.Segment{BizTypeID: b.ID, StartID: 1, EndID: 1001})
	rec, err := s.CreateRecycleRecord(model.RecycleRecord{SegmentID: seg.ID, BizTypeID: b.ID, Reason: "测试回收"})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	_, err = s.UpdateRecycleRecord(rec.ID, model.RecycleRecord{Status: model.RecycleStatusRecycled})
	if err != nil {
		t.Fatalf("状态流转失败: %v", err)
	}
	// 状态机：recycled 不能再流转
	_, err = s.UpdateRecycleRecord(rec.ID, model.RecycleRecord{Status: model.RecycleStatusPending})
	if err == nil {
		t.Fatalf("recycled -> pending 应被禁止")
	}
}

func TestService_AllocRecord(t *testing.T) {
	s := newTestService()
	b, _ := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100})
	n, _ := s.CreateMachineNode(model.MachineNode{NodeID: "node-1", Hostname: "host1", IP: "127.0.0.1"})
	ar, err := s.CreateAllocRecord(model.AllocRecord{BizTypeID: b.ID, NodeID: n.NodeID, BatchSize: 10, Mode: model.BizModeSegment})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	if ar.ID == "" {
		t.Fatalf("ID 不应为空")
	}
	items, total, err := s.ListAllocRecords(model.AllocRecordFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("列表失败: %v", err)
	}
	if total != 1 {
		t.Fatalf("总数不匹配")
	}
	if len(items) != 1 {
		t.Fatalf("列表长度不匹配")
	}
}

func TestService_NodeHeartbeat(t *testing.T) {
	s := newTestService()
	n, _ := s.CreateMachineNode(model.MachineNode{NodeID: "node-1", Hostname: "host1", IP: "127.0.0.1"})
	_, err := s.CreateNodeHeartbeat(model.NodeHeartbeat{NodeID: n.NodeID, Load: 10.5, CPUUsage: 20.0, MemoryUsage: 30.0})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	// 外键校验失败
	_, err = s.CreateNodeHeartbeat(model.NodeHeartbeat{NodeID: "not-exist", Load: 0})
	if err == nil {
		t.Fatalf("外键校验应失败")
	}
}

func TestService_SnowflakeConfig(t *testing.T) {
	s := newTestService()
	b, _ := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSnowflake, SegmentStep: 0})
	_, err := s.CreateSnowflakeConfig(model.SnowflakeConfig{BizTypeID: b.ID, DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12, Twepoch: 1704067200000})
	if err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	// 重复创建应失败
	_, err = s.CreateSnowflakeConfig(model.SnowflakeConfig{BizTypeID: b.ID, DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12, Twepoch: 1704067200000})
	if err == nil {
		t.Fatalf("重复创建应失败")
	}
}

func TestService_Stats(t *testing.T) {
	s := newTestService()
	b, _ := s.CreateBizType(model.BizType{Name: "订单", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100})
	n, _ := s.CreateMachineNode(model.MachineNode{NodeID: "node-1", Hostname: "host1", IP: "127.0.0.1"})
	_, _ = s.CreateAllocStats(model.AllocStats{BizTypeID: b.ID, NodeID: n.NodeID, Date: "2024-01-01", TotalAllocated: 1000, PeakQPS: 100.5, AvgBatchSize: 10.0})
	overview := s.GetStatsOverview()
	if overview["biz_count"] != 1 {
		t.Fatalf("业务计数不匹配")
	}
	if s.CountBizTypes() != 1 {
		t.Fatalf("CountBizTypes 不匹配")
	}
}

func TestService_Pagination(t *testing.T) {
	s := newTestService()
	for i := 0; i < 25; i++ {
		_, _ = s.CreateBizType(model.BizType{Name: "业务" + string(rune('A'+i)), Code: "code" + string(rune('A'+i)), Mode: model.BizModeSegment, SegmentStep: 100})
	}
	items, total, err := s.ListBizTypes(model.BizTypeFilter{}, 1, 10)
	if err != nil {
		t.Fatalf("列表失败: %v", err)
	}
	if total != 25 {
		t.Fatalf("总数应为 25")
	}
	if len(items) != 10 {
		t.Fatalf("第 1 页应为 10 条")
	}
	items2, _, _ := s.ListBizTypes(model.BizTypeFilter{}, 3, 10)
	if len(items2) != 5 {
		t.Fatalf("第 3 页应为 5 条")
	}
}
