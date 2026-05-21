# 图片上传模块设计

## 目标

支持用户上传图片，用于帖子、商品等场景。开发阶段存本地，生产环境迁移至 Cloudflare R2。

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|---|---|---|---|
| POST | /api/v1/upload | 上传单张或多张图片 | 是 |

## 请求格式

multipart/form-data，字段名 `files`，支持多文件。

## 响应

```json
{
  "urls": [
    "http://localhost:8080/uploads/1234567890.jpg"
  ]
}
```

## 限制

- 单文件最大 5MB
- 仅支持 jpg、jpeg、png、gif、webp
- 最多一次上传 9 张

## 开发阶段存储

文件保存到 `server/assets/uploads/`，通过 `/uploads/:filename` 提供静态访问。

## 涉及文件

- internal/handler/upload.go
- internal/service/upload.go
- internal/router/router.go
- cmd/community_server/main.go
