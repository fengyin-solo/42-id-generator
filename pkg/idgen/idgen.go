// Package idgen 提供 ID 与短码生成能力，包含雪花算法和号段模式。
package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Hex 生成 8 字节随机值的十六进制字符串（16 位）。
func Hex() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// HexN 生成指定字节数的十六进制 ID。
func HexN(n int) string {
	if n <= 0 {
		n = 8
	}
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var seq uint64

// Short 生成递增 base62 短码。
func Short() string {
	n := atomic.AddUint64(&seq, 1)
	if n == 0 {
		return string(base62Chars[0])
	}
	buf := make([]byte, 0, 11)
	for n > 0 {
		buf = append(buf, base62Chars[n%62])
		n /= 62
	}
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}

// Snowflake 实现雪花算法（时间戳 + 数据中心 + 机器节点 + 序列号）。
type Snowflake struct {
	mu            sync.Mutex
	lastTimestamp int64
	sequence      int64
	datacenterID  int64
	workerID      int64

	// 可配置位数
	sequenceBits  uint8
	workerBits    uint8
	datacenterBits uint8
	epoch         int64
}

// NewSnowflake 创建雪花发号器。
func NewSnowflake(datacenterID, workerID int64, sequenceBits, workerBits, datacenterBits uint8, epoch int64) (*Snowflake, error) {
	maxWorker := int64(-1) ^ (int64(-1) << workerBits)
	maxDatacenter := int64(-1) ^ (int64(-1) << datacenterBits)
	if workerID < 0 || workerID > maxWorker {
		return nil, errors.New("workerID 超出范围")
	}
	if datacenterID < 0 || datacenterID > maxDatacenter {
		return nil, errors.New("datacenterID 超出范围")
	}
	if epoch == 0 {
		epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	}
	return &Snowflake{
		datacenterID:   datacenterID,
		workerID:       workerID,
		sequenceBits:   sequenceBits,
		workerBits:     workerBits,
		datacenterBits: datacenterBits,
		epoch:          epoch,
	}, nil
}

// NextID 生成下一个 64 位唯一 ID。
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ts := time.Now().UnixMilli()
	if ts < s.lastTimestamp {
		return 0, errors.New("时钟回拨")
	}

	maxSequence := int64(-1) ^ (int64(-1) << s.sequenceBits)
	if ts == s.lastTimestamp {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			for ts <= s.lastTimestamp {
				ts = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTimestamp = ts

	workerShift := s.sequenceBits
	datacenterShift := workerShift + s.workerBits
	timestampShift := datacenterShift + s.datacenterBits

	id := ((ts - s.epoch) << timestampShift) |
		(s.datacenterID << datacenterShift) |
		(s.workerID << workerShift) |
		s.sequence

	return id, nil
}

// ParseID 解析雪花 ID 的组成部分。
func (s *Snowflake) ParseID(id int64) (timestamp int64, datacenterID, workerID, sequence int64) {
	workerShift := s.sequenceBits
	datacenterShift := workerShift + s.workerBits
	timestampShift := datacenterShift + s.datacenterBits

	maxSequence := int64(-1) ^ (int64(-1) << s.sequenceBits)
	maxWorker := int64(-1) ^ (int64(-1) << s.workerBits)
	maxDatacenter := int64(-1) ^ (int64(-1) << s.datacenterBits)

	sequence = id & maxSequence
	workerID = (id >> workerShift) & maxWorker
	datacenterID = (id >> datacenterShift) & maxDatacenter
	timestamp = (id >> timestampShift) + s.epoch
	return
}

// SegmentGenerator 号段发号器（游标推进）。
type SegmentGenerator struct {
	mu     sync.Mutex
	start  int64
	end    int64
	cursor int64
}

// NewSegmentGenerator 创建号段发号器。
func NewSegmentGenerator(start, end int64) *SegmentGenerator {
	return &SegmentGenerator{
		start:  start,
		end:    end,
		cursor: start,
	}
}

// Next 取下一个号，返回是否耗尽。
func (g *SegmentGenerator) Next() (int64, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.cursor >= g.end {
		return 0, true
	}
	id := g.cursor
	g.cursor++
	return id, false
}

// Remain 返回剩余可用数量。
func (g *SegmentGenerator) Remain() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.cursor >= g.end {
		return 0
	}
	return g.end - g.cursor
}

// Cursor 返回当前游标位置。
func (g *SegmentGenerator) Cursor() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.cursor
}

// BatchNext 批量取号，返回 ids 和是否全部满足。
func (g *SegmentGenerator) BatchNext(n int) ([]int64, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if n <= 0 {
		return []int64{}, true
	}
	remain := g.end - g.cursor
	if int64(n) > remain {
		n = int(remain)
	}
	ids := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		if g.cursor >= g.end {
			break
		}
		ids = append(ids, g.cursor)
		g.cursor++
	}
	all := len(ids) == n
	return ids, all
}

// SequenceGenerator 数据库序列模拟发号器（纯递增）。
type SequenceGenerator struct {
	mu  sync.Mutex
	seq int64
}

// NewSequenceGenerator 创建序列发号器。
func NewSequenceGenerator(start int64) *SequenceGenerator {
	return &SequenceGenerator{seq: start}
}

// Next 返回下一个序列号。
func (g *SequenceGenerator) Next() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.seq++
	return g.seq
}

// Current 返回当前序列值。
func (g *SequenceGenerator) Current() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.seq
}
