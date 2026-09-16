package store

import (
	"sync"

	"idgenerator/internal/model"
)

type MemoryStore struct {
	mu               sync.RWMutex
	bizTypes         map[string]*model.BizType
	machineNodes     map[string]*model.MachineNode
	idRules          map[string]*model.IDRule
	segments         map[string]*model.Segment
	allocRecords     map[string]*model.AllocRecord
	leases           map[string]*model.Lease
	nodeHeartbeats   map[string]*model.NodeHeartbeat
	snowflakeConfigs map[string]*model.SnowflakeConfig
	allocStats       map[string]*model.AllocStats
	recycleRecords   map[string]*model.RecycleRecord
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bizTypes:         make(map[string]*model.BizType),
		machineNodes:     make(map[string]*model.MachineNode),
		idRules:          make(map[string]*model.IDRule),
		segments:         make(map[string]*model.Segment),
		allocRecords:     make(map[string]*model.AllocRecord),
		leases:           make(map[string]*model.Lease),
		nodeHeartbeats:   make(map[string]*model.NodeHeartbeat),
		snowflakeConfigs: make(map[string]*model.SnowflakeConfig),
		allocStats:       make(map[string]*model.AllocStats),
		recycleRecords:   make(map[string]*model.RecycleRecord),
	}
}

var _ Store = (*MemoryStore)(nil)
