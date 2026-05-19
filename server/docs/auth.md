# 用户认证模块设计

## 目标

实现用户注册、登录、Token 刷新功能，为后续所有需要身份认证的接口提供基础。

## 登录方式

开发阶段使用手机号+密码登录，后续切换为微信小程序手机号授权。

## Token 机制

JWT 双 Token：
- access_token：有效期 2 小时，用于接口认证
- refresh_token：有效期 7 天，用于换取新的 access_token

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|---|---|---|---|
| POST | /api/v1/auth/register | 注册 | 否 |
| POST | /api/v1/auth/login | 登录 | 否 |
| POST | /api/v1/auth/refresh | 刷新 Token | 否 |
| GET | /api/v1/user/profile | 获取当前用户信息 | 是 |

## 数据库变更

users 表新增 password 字段（002_add_user_password.sql）。

## 涉及文件

- pkg/jwt/jwt.go — JWT 生成与解析
- internal/model/user.go — 用户结构体
- internal/service/auth.go — 认证业务逻辑
- internal/handler/auth.go — 接口处理器
- internal/middleware/auth.go — JWT 认证中间件
- internal/config/config.go — 新增 JWT 配置
- internal/router/router.go — 注册路由
- cmd/community_server/main.go — 串联各模块
