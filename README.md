# 分布式 ID 生成服务 (Distributed ID Generator)

纯 Go 标准库实现的分布式 ID 生成后端服务，零第三方依赖，支持雪花算法（Snowflake）、号段（Segment）与递增序列三种发号模式，并管理节点注册、租约续期、心跳、号段回收的完整生命周期。

## 运行说明

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`，API Key 为 `default-api-key`（通过 Header `X-Api-Key` 或 Query `api_key` 传递）。

访问 `http://localhost:8080/` 可查看前端看板页面。

## API 表格

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/biz-types | 创建业务标识 |
| GET | /api/biz-types | 列出业务标识（支持 name/code/mode/enabled/keyword 筛选） |
| GET | /api/biz-types/{id} | 获取业务标识 |
| PUT | /api/biz-types/{id} | 更新业务标识 |
| DELETE | /api/biz-types/{id} | 删除业务标识 |
| POST | /api/biz-types/{id}/toggle | 启停业务标识 |
| POST | /api/nodes | 注册机器节点 |
| GET | /api/nodes | 列出机器节点（支持 node_id/status/hostname/keyword 筛选） |
| GET | /api/nodes/{id} | 获取机器节点 |
| PUT | /api/nodes/{id} | 更新机器节点 |
| DELETE | /api/nodes/{id} | 删除机器节点 |
| GET | /api/nodes/{id}/health | 检查节点健康状态 |
| POST | /api/id-rules | 创建发号规则 |
| GET | /api/id-rules | 列出发号规则（支持 name/biz_type_id/mode/enabled/keyword 筛选） |
| GET | /api/id-rules/{id} | 获取发号规则 |
| PUT | /api/id-rules/{id} | 更新发号规则 |
| DELETE | /api/id-rules/{id} | 删除发号规则 |
| POST | /api/id-rules/{id}/toggle | 启停发号规则 |
| POST | /api/segments | 创建号段 |
| GET | /api/segments | 列出号段（支持 biz_type_id/status/keyword 筛选） |
| GET | /api/segments/{id} | 获取号段 |
| PUT | /api/segments/{id} | 更新号段 |
| DELETE | /api/segments/{id} | 删除号段 |
| POST | /api/segments/{id}/advance | 推进号段游标（状态机 using->used->exhausted） |
| POST | /api/alloc-records | 创建分配记录 |
| GET | /api/alloc-records | 列出分配记录（支持 biz_type_id/node_id/keyword 筛选） |
| GET | /api/alloc-records/{id} | 获取分配记录 |
| DELETE | /api/alloc-records/{id} | 删除分配记录 |
| POST | /api/leases | 创建租约 |
| GET | /api/leases | 列出租约（支持 node_id/biz_type_id/status/keyword 筛选） |
| GET | /api/leases/{id} | 获取租约 |
| POST | /api/leases/{id}/renew | 续约 |
| POST | /api/leases/{id}/expire | 过期（状态机 active->expired） |
| DELETE | /api/leases/{id} | 删除租约 |
| POST | /api/heartbeats | 上报节点心跳 |
| GET | /api/heartbeats | 列出心跳（支持 node_id/keyword 筛选） |
| GET | /api/heartbeats/{id} | 获取心跳 |
| DELETE | /api/heartbeats/{id} | 删除心跳 |
| GET | /api/nodes/{id}/latest-heartbeat | 获取节点最近心跳 |
| POST | /api/snowflake-configs | 创建雪花配置 |
| GET | /api/snowflake-configs | 列出雪花配置（支持 biz_type_id/keyword 筛选） |
| GET | /api/snowflake-configs/{id} | 获取雪花配置 |
| PUT | /api/snowflake-configs/{id} | 更新雪花配置 |
| DELETE | /api/snowflake-configs/{id} | 删除雪花配置 |
| POST | /api/alloc-stats | 创建分配统计 |
| GET | /api/alloc-stats | 列出分配统计（支持 biz_type_id/node_id/date 筛选） |
| GET | /api/alloc-stats/{id} | 获取分配统计 |
| PUT | /api/alloc-stats/{id} | 更新分配统计 |
| DELETE | /api/alloc-stats/{id} | 删除分配统计 |
| POST | /api/recycle-records | 创建号段回收记录 |
| GET | /api/recycle-records | 列出回收记录（支持 biz_type_id/status/keyword 筛选） |
| GET | /api/recycle-records/{id} | 获取回收记录 |
| PUT | /api/recycle-records/{id} | 更新回收记录 |
| DELETE | /api/recycle-records/{id} | 删除回收记录 |
| POST | /api/id-gen/snowflake | 生成雪花 ID |
| POST | /api/id-gen/segment | 生成号段 ID |
| POST | /api/id-gen/sequence | 生成递增序列 ID |
| POST | /api/id-gen/batch | 批量生成 ID |
| GET | /api/stats/overview | 全局统计概览 |
| GET | /api/stats/top | TOP N 分配统计 |
| GET | /api/stats/group-by-biz | 按业务分组统计 |
| GET | /api/stats/group-by-node | 按节点分组统计 |
| GET | /api/export | 导出全量 JSON 快照 |
| POST | /api/import | 导入 JSON 快照 |
| GET | / | 前端看板页面 |

## 实体列表

1. BizType（业务标识）
2. MachineNode（机器节点）
3. IDRule（发号规则）
4. Segment（号段）
5. AllocRecord（分配记录）
6. Lease（租约）
7. NodeHeartbeat（节点心跳）
8. SnowflakeConfig（雪花配置）
9. AllocStats（分配统计）
10. RecycleRecord（号段回收）
