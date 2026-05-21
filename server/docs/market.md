# 广场集市模块设计

## 目标

实现二手交易集市功能：用户可在已订阅的社区发布商品，浏览商品列表，管理自己的商品（下架、标记已售）。

## 商品属性

- 标题：长度 1-100 字符
- 价格：正整数（单位：分）
- 新旧程度：new（全新）, like_new（几乎全新）, lightly_used（轻微使用）, heavily_used（明显使用）
- 分类：appliance（家电）, furniture（家具）, book（书籍）, baby（母婴）, sports（运动）, other（其他）
- 图片：最多 9 张
- 状态：selling（出售中）, sold（已售出）, off（已下架）

## 业务规则

- 只能在自己已订阅的社区发布商品
- 列表按社区隔离，可按分类筛选
- 只有卖家本人可下架或标记已售
- 已下架/已售出的商品不在列表中展示
- 私信卖家由前端直接调用腾讯云 IM SDK，后端集市模块只提供卖家基本信息

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|---|---|---|---|
| POST | /api/v1/items | 发布商品 | 是 |
| GET | /api/v1/items?community_id=xxx&category=xxx&last_id=xxx | 商品列表 | 是 |
| GET | /api/v1/items/:id | 商品详情 | 是 |
| PUT | /api/v1/items/:id | 修改商品信息 | 是 |
| PUT | /api/v1/items/:id/sold | 标记已售 | 是 |
| PUT | /api/v1/items/:id/off | 下架商品 | 是 |

## 数据库变更

- items 表 — 商品主表

## 涉及文件

- sql/004_market_schema.sql — 数据库迁移
- internal/model/item.go — 商品模型
- internal/service/market.go — 集市业务逻辑
- internal/handler/market.go — 集市接口处理器
- internal/router/router.go — 注册路由
- cmd/community_server/main.go — 注入服务
