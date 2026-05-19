# 用户模块设计

## 目标

完善用户资料管理功能，支持更新昵称、头像，查看他人主页。

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|---|---|---|---|
| PUT | /api/v1/user/profile | 更新当前用户资料 | 是 |
| GET | /api/v1/users/:id | 获取指定用户公开资料 | 否 |

## 更新资料说明

- 昵称：长度 1-20 字符
- 头像：传入图片 URL（由上层上传服务提供）
- 只允许修改自己的资料

## 涉及文件

- internal/handler/auth.go — Profile 补充 avatar_url、created_at
- internal/handler/user.go — 用户资料接口处理器
- internal/service/user.go — 用户资料业务逻辑
- internal/router/router.go — 注册路由
- cmd/community_server/main.go — 注入服务
