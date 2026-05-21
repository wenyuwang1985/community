# 前端目录说明

本项目前端分为两个独立目录，分别对应不同场景：

```
client/
├── web/        # 浏览器端 SPA（Vue 3 CDN，免构建）
└── uni-app/    # uni-app 项目（Vue 3 + Vite，编译到小程序/H5/App）
```

---

## web/ — 浏览器端

**技术栈**：Vue 3 全局构建版 + Vue Router 4 CDN + 原生 fetch/WebSocket

**适用场景**：
- 快速在浏览器中预览效果
- 不方便安装 Node.js 的环境
- PC 浏览器访问

**运行方式**：

```bash
cd client/web
# 需要本地 HTTP 服务器（ES Modules 不支持 file:// 协议）
python3 -m http.server 3000
# 然后访问 http://localhost:3000
```

**特性**：
- 单页应用（SPA），Hash 路由
- 底部 TabBar 模拟移动端 App 体验
- 对接后端 localhost:8080

---

## uni-app/ — 多端项目

**技术栈**：uni-app 3 (Vue 3 + Vite)

**适用场景**：
- 微信小程序（第一期主战场）
- H5 浏览器
- 未来可编译到 iOS/Android App

**运行方式**：

```bash
cd client/uni-app
npm install
npm run dev:h5         # H5 开发
npm run dev:mp-weixin  # 微信小程序开发
```

**项目结构**：

| 文件/目录 | 说明 |
|---|---|
| `manifest.json` | 各端应用配置（AppID、H5 代理、小程序设置） |
| `pages.json` | 页面路由 + 全局样式 + TabBar |
| `src/utils/api.js` | 封装 `uni.request` 调用后端 API |
| `src/utils/ws.js` | 封装 `uni.connectSocket` 聊天长连接 |
| `src/pages/*/index.vue` | 页面组件 |

**页面清单**：

| 页面 | 路径 | 类型 |
|---|---|---|
| 登录 | `src/pages/login/index` | 普通页面 |
| 注册 | `src/pages/register/index` | 普通页面 |
| 广场 | `src/pages/index/index` | **TabBar** |
| 集市 | `src/pages/market/index` | **TabBar** |
| 聊天 | `src/pages/chat/index` | **TabBar** |
| 聊天房间 | `src/pages/chat-room/index` | 普通页面 |
| 社区管理 | `src/pages/communities/index` | 普通页面 |
| 我的 | `src/pages/profile/index` | **TabBar** |

---

## 开发规范

- 两个目录共享同一套业务逻辑和 API 接口定义
- `uni-app/` 中的组件使用 uni-app 规范（`<view>`、`<text>`、`<image>`）
- `web/` 中的组件使用标准 HTML（`<div>`、`<span>`、`<img>`）
- 样式变量主色：`#07c160`（微信绿）

## 后端对接

默认后端地址：`http://localhost:8080`

如需修改，请编辑：
- `web/src/api.js` 中的 `BASE`
- `uni-app/src/utils/api.js` 中的 `BASE`
- `uni-app/manifest.json` 中 `h5.devServer.proxy` 的 `target`
