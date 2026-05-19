# 社区订阅模块设计

## 目标

实现街镇搜索与订阅功能，用户可搜索街镇、订阅/取消订阅、管理已订阅列表、切换主社区。

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|---|---|---|---|
| GET | /api/v1/communities/search?q=xxx | 按名称搜索街镇 | 否 |
| POST | /api/v1/communities/:id/subscribe | 订阅街镇 | 是 |
| DELETE | /api/v1/communities/:id/subscribe | 取消订阅 | 是 |
| GET | /api/v1/user/communities | 获取已订阅社区列表 | 是 |
| PUT | /api/v1/user/communities/:id/primary | 设为主社区 | 是 |

## 业务规则

- 用户必须至少订阅一个社区
- 取消订阅时，若只剩最后一个订阅，禁止取消
- 设为主社区时，自动将该社区设为 is_primary=true，其他设为 false
- 搜索按名称模糊匹配，返回前 20 条

## 涉及文件

- internal/model/community.go — 社区结构体
- internal/handler/community.go — 社区接口处理器
- internal/service/community.go — 社区业务逻辑
- internal/router/router.go — 注册路由
- cmd/community_server/main.go — 注入服务
