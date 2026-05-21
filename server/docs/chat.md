# 聊天模块设计（自实现 WebSocket）

## 目标

自实现即时通讯，覆盖社区频道（公共群聊）、私信（一对一）、消息通知。

## 技术选型

- WebSocket 库：`gorilla/websocket`
- 消息持久化：PostgreSQL
- 在线状态：内存 Hub 维护 userID → conn 映射

## 数据模型

### conversations 表

- id BIGINT PK — Snowflake
- type VARCHAR(20) — private / channel
- community_id BIGINT FK — channel 类型关联街镇
- name VARCHAR(100) — 频道/会话名称
- created_by BIGINT FK — 创建者
- created_at TIMESTAMPTZ

### conversation_participants 表

- conversation_id BIGINT FK
- user_id BIGINT FK
- joined_at TIMESTAMPTZ
- PK(conversation_id, user_id)

### messages 表

- id BIGINT PK — Snowflake
- conversation_id BIGINT FK
- sender_id BIGINT FK
- content TEXT
- type VARCHAR(20) — text / image / system
- created_at TIMESTAMPTZ

## WebSocket 协议

### 连接建立

`GET /ws?token={access_token}`

服务端校验 JWT，提取 userID，注册到 Hub。

### 客户端 → 服务端

```json
// 发送消息
{"action": "send_message", "conversation_id": 123, "content": "hello"}

// 心跳
{"action": "ping"}
```

### 服务端 → 客户端

```json
// 新消息推送
{"type": "message", "data": {"id": 1, "conversation_id": 123, "sender_id": 1, "sender_nickname": "张三", "content": "hello", "created_at": "..."}}

// 心跳响应
{"type": "pong"}

// 错误
{"type": "error", "msg": "..."}
```

## 推送逻辑

1. 客户端通过 WebSocket 发送 `send_message`
2. 服务端校验用户是否属于该 conversation
3. 消息落库（messages 表）
4. 查询该 conversation 所有在线参与者
5. 通过 WebSocket 逐个推送
6. 离线用户下次上线后通过 HTTP API 拉取历史消息

## 社区频道机制

- 每个街镇在创建时（或首次有用户订阅时）自动生成一个 type=channel 的 conversation
- 用户订阅该街镇时，自动加入对应 conversation_participants
- 用户取消订阅时，自动从 participants 中移除
- 新用户加入频道后可查看全部历史消息

## HTTP API

| 方法 | 路径 | 说明 | 认证 |
|---|---|---|---|
| GET | /ws | WebSocket 连接 | query token |
| GET | /api/v1/conversations | 我的会话列表 | 是 |
| GET | /api/v1/conversations/:id/messages?last_id=xxx | 历史消息 | 是 |
| POST | /api/v1/conversations/private | 发起/获取私聊会话 | 是 |

## 涉及文件

- go.mod — 新增 gorilla/websocket
- sql/005_chat_schema.sql
- internal/model/chat.go
- internal/ws/hub.go / client.go
- internal/service/chat.go
- internal/handler/chat.go / ws.go
- internal/router/router.go
- cmd/community_server/main.go
