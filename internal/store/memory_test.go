package store

import (
	"testing"
	"time"

	"idgenerator/internal/model"
)

func TestMemoryStore_BizType(t *testing.T) {
	s := NewMemoryStore()
	b := &model.BizType{ID: "b1", Name: "订单", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100, Status: model.BizStatusEnabled, CreatedAt: time.Now()}
	if err := s.CreateBizType(b); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	got, err := s.GetBizType("b1")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if got.Code != "order" {
		t.Fatalf("Code 不匹配")
	}
	_, err = s.GetBizTypeByCode("order")
	if err != nil {
		t.Fatalf("按 Code 查询失败: %v", err)
	}
	if len(s.ListBizTypes()) != 1 {
		t.Fatalf("列表长度不匹配")
	}
	// 冲突
	b2 := &model.BizType{ID: "b2", Name: "订单2", Code: "order", Mode: model.BizModeSegment, SegmentStep: 100}
	if err := s.CreateBizType(b2); err != ErrConflict {
		t.Fatalf("应返回冲突: %v", err)
	}
	b.Name = "订单更新"
	if err := s.UpdateBizType(b); err != nil {
		t.Fatalf("更新失败: %v", err)
	}
	if err := s.DeleteBizType("b1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
	if _, err := s.GetBizType("b1"); err != ErrNotFound {
		t.Fatalf("应返回不存在: %v", err)
	}
}

func TestMemoryStore_MachineNode(t *testing.T) {
	s := NewMemoryStore()
	n := &model.MachineNode{ID: "n1", NodeID: "node-1", Hostname: "host1", IP: "127.0.0.1", WorkerID: 1, DatacenterID: 1, Status: model.NodeStatusActive, CreatedAt: time.Now()}
	if err := s.CreateMachineNode(n); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	got, err := s.GetMachineNode("n1")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	if got.NodeID != "node-1" {
		t.Fatalf("NodeID 不匹配")
	}
	_, err = s.GetMachineNodeByNodeID("node-1")
	if err != nil {
		t.Fatalf("按 NodeID 查询失败: %v", err)
	}
	n2 := &model.MachineNode{ID: "n2", NodeID: "node-1", Hostname: "host2", IP: "127.0.0.2", WorkerID: 2}
	if err := s.CreateMachineNode(n2); err != ErrConflict {
		t.Fatalf("应返回冲突: %v", err)
	}
	if err := s.DeleteMachineNode("n1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_IDRule(t *testing.T) {
	s := NewMemoryStore()
	r := &model.IDRule{ID: "r1", Name: "规则1", BizTypeID: "b1", Mode: model.BizModeSnowflake, SignBits: 1, TimestampBits: 41, DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12, Status: model.IDRuleStatusEnabled, CreatedAt: time.Now()}
	if err := s.CreateIDRule(r); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	_, err := s.GetIDRuleByBizTypeID("b1")
	if err != nil {
		t.Fatalf("按 BizTypeID 查询失败: %v", err)
	}
	r2 := &model.IDRule{ID: "r2", Name: "规则2", BizTypeID: "b1", Mode: model.BizModeSnowflake, SignBits: 1, TimestampBits: 41, DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12}
	if err := s.CreateIDRule(r2); err != ErrConflict {
		t.Fatalf("应返回冲突: %v", err)
	}
	if err := s.DeleteIDRule("r1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_Segment(t *testing.T) {
	s := NewMemoryStore()
	seg := &model.Segment{ID: "s1", BizTypeID: "b1", StartID: 1, EndID: 1001, Cursor: 1, Status: model.SegmentStatusUsing, CreatedAt: time.Now()}
	if err := s.CreateSegment(seg); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	got, err := s.GetActiveSegmentByBizTypeID("b1")
	if err != nil {
		t.Fatalf("查询活跃号段失败: %v", err)
	}
	if got.ID != "s1" {
		t.Fatalf("ID 不匹配")
	}
	if len(s.ListSegmentsByBizTypeID("b1")) != 1 {
		t.Fatalf("按业务查询长度不匹配")
	}
	if err := s.DeleteSegment("s1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_AllocRecord(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AllocRecord{ID: "a1", BizTypeID: "b1", NodeID: "n1", BatchSize: 10, Mode: model.BizModeSegment, AllocatedAt: time.Now(), CreatedAt: time.Now()}
	if err := s.CreateAllocRecord(a); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	if len(s.ListAllocRecords()) != 1 {
		t.Fatalf("列表长度不匹配")
	}
	if err := s.DeleteAllocRecord("a1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_Lease(t *testing.T) {
	s := NewMemoryStore()
	l := &model.Lease{ID: "l1", NodeID: "n1", BizTypeID: "b1", ExpiresAt: time.Now().Add(time.Hour), Status: model.LeaseStatusActive, CreatedAt: time.Now()}
	if err := s.CreateLease(l); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	got, err := s.GetActiveLeaseByNodeAndBiz("n1", "b1")
	if err != nil {
		t.Fatalf("查询活跃租约失败: %v", err)
	}
	if got.ID != "l1" {
		t.Fatalf("ID 不匹配")
	}
	if err := s.DeleteLease("l1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_NodeHeartbeat(t *testing.T) {
	s := NewMemoryStore()
	h := &model.NodeHeartbeat{ID: "h1", NodeID: "n1", Load: 10.5, CPUUsage: 20.0, MemoryUsage: 30.0, BeatAt: time.Now(), CreatedAt: time.Now()}
	if err := s.CreateNodeHeartbeat(h); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	got, err := s.GetLatestHeartbeatByNodeID("n1")
	if err != nil {
		t.Fatalf("查询最新心跳失败: %v", err)
	}
	if got.ID != "h1" {
		t.Fatalf("ID 不匹配")
	}
	if err := s.DeleteNodeHeartbeat("h1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_SnowflakeConfig(t *testing.T) {
	s := NewMemoryStore()
	c := &model.SnowflakeConfig{ID: "c1", BizTypeID: "b1", EpochMs: 0, DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12, Twepoch: 1704067200000, CreatedAt: time.Now()}
	if err := s.CreateSnowflakeConfig(c); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	_, err := s.GetSnowflakeConfigByBizTypeID("b1")
	if err != nil {
		t.Fatalf("按 BizTypeID 查询失败: %v", err)
	}
	c2 := &model.SnowflakeConfig{ID: "c2", BizTypeID: "b1", DatacenterBits: 5, WorkerBits: 5, SequenceBits: 12, Twepoch: 1704067200000}
	if err := s.CreateSnowflakeConfig(c2); err != ErrConflict {
		t.Fatalf("应返回冲突: %v", err)
	}
}

func TestMemoryStore_AllocStats(t *testing.T) {
	s := NewMemoryStore()
	st := &model.AllocStats{ID: "st1", BizTypeID: "b1", NodeID: "n1", Date: "2024-01-01", TotalAllocated: 100, PeakQPS: 10.5, AvgBatchSize: 5.0, CreatedAt: time.Now()}
	if err := s.CreateAllocStats(st); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	got, err := s.GetAllocStatsByBizNodeDate("b1", "n1", "2024-01-01")
	if err != nil {
		t.Fatalf("按组合键查询失败: %v", err)
	}
	if got.ID != "st1" {
		t.Fatalf("ID 不匹配")
	}
	st.TotalAllocated = 200
	if err := s.UpdateAllocStats(st); err != nil {
		t.Fatalf("更新失败: %v", err)
	}
	if err := s.DeleteAllocStats("st1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}

func TestMemoryStore_RecycleRecord(t *testing.T) {
	s := NewMemoryStore()
	rec := &model.RecycleRecord{ID: "rec1", SegmentID: "s1", BizTypeID: "b1", Reason: "节点下线", Status: model.RecycleStatusPending, CreatedAt: time.Now()}
	if err := s.CreateRecycleRecord(rec); err != nil {
		t.Fatalf("创建失败: %v", err)
	}
	if len(s.ListRecycleRecords()) != 1 {
		t.Fatalf("列表长度不匹配")
	}
	if err := s.DeleteRecycleRecord("rec1"); err != nil {
		t.Fatalf("删除失败: %v", err)
	}
}
