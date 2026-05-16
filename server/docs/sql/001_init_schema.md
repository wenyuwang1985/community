# 001_init_schema 表结构设计说明

## 概述

本文件定义社区生活系统的三张核心基础表，支撑用户注册、街镇管理、社区订阅功能。

## 主键策略：Snowflake ID

所有表主键使用 BIGINT 类型，由 Go 应用层通过 `bwmarrin/snowflake` 库生成。

- 趋势递增，对 B-Tree 索引友好
- 全局唯一，未来可无缝扩展到分库分表
- 64 位整数，比 UUID 更紧凑，URL 更短

每个服务实例需分配不同的 node ID（0~1023），单机部署时固定为 1 即可。

## 表说明

### communities（街镇表）

存储街镇级行政单位信息。数据来源为民政部公开数据或高德行政区划 API。

- `adcode`：12 位行政区划代码，用于与外部地图 API 对接
- `boundary`：PostGIS 几何字段，存储街镇边界多边形，用于 GPS 定位判断用户所在街镇。初期可为空，后续按需填充

### users（用户表）

- `phone`：唯一约束，微信小程序手机号授权获取，作为账号唯一标识
- `credit_score`：信用分，初始 100，后续通过交易评价、举报等机制调整
- `updated_at`：需要应用层在更新时手动设置，或后续添加触发器自动维护

### user_community_subscriptions（订阅关系表）

用户与街镇的多对多关系。

- `is_primary`：标记用户的主社区（首页默认展示），每个用户应只有一条 is_primary=true 的记录，此约束由应用层保证
- `UNIQUE(user_id, community_id)`：防止同一用户重复订阅同一街镇

## 索引说明

| 索引 | 用途 |
|---|---|
| idx_communities_adcode | 按行政区划代码快速查找街镇 |
| idx_subscriptions_user_id | 查询用户已订阅的所有社区 |
| idx_subscriptions_community_id | 查询某社区的所有订阅用户 |

## 后续扩展

下一批表（posts、items 等）将在 `002_*.sql` 中定义，届时会引用 users 和 communities 的外键。
