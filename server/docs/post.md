# 广场动态模块设计

## 目标

实现社区动态（Feed）功能：用户可在已订阅的社区发帖，浏览社区 Feed 流，对帖子点赞、评论。

## 发帖规则

- 只能在自己已订阅的社区发帖
- 内容长度 1-2000 字符
- 图片 URL 数组，最多 9 张
- Tag 可选：help（求助）、share（分享）、notice（通知）、qa（问答）

## Feed 流规则

- 按社区隔离，不同社区内容不混合
- 按发布时间倒序排列
- Cursor 分页（传入 last_id，返回下一页）
- 每页默认 10 条
- 返回作者昵称、头像、当前用户是否点赞

## 评论规则

- 内容长度 1-500 字符
- 按时间正序排列
- 返回评论者昵称、头像

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|---|---|---|---|
| POST | /api/v1/posts | 发帖 | 是 |
| GET | /api/v1/posts?community_id=xxx&last_id=xxx | Feed 流 | 是 |
| GET | /api/v1/posts/:id | 帖子详情 | 是 |
| DELETE | /api/v1/posts/:id | 删除自己的帖子 | 是 |
| POST | /api/v1/posts/:id/like | 点赞 | 是 |
| DELETE | /api/v1/posts/:id/like | 取消点赞 | 是 |
| POST | /api/v1/posts/:id/comments | 发表评论 | 是 |
| GET | /api/v1/posts/:id/comments | 评论列表 | 是 |

## 数据库变更

- posts 表 — 帖子主表
- comments 表 — 评论表
- post_likes 表 — 点赞关系表

## 涉及文件

- sql/003_posts_schema.sql — 数据库迁移
- internal/model/post.go — 帖子、评论模型
- internal/service/post.go — 帖子业务逻辑
- internal/handler/post.go — 帖子接口处理器
- internal/router/router.go — 注册路由
- cmd/community_server/main.go — 注入服务
