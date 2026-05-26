# 社区生活 uni-app 前端

街镇级社区连接平台的前端项目，基于 uni-app + Vue 3 开发。

## 项目结构

```
client/uni-app/
├── src/
│   ├── pages/          # 页面
│   │   ├── login/      # 登录页面
│   │   ├── register/   # 注册页面
│   │   ├── index/      # 首页（动态Feed）
│   │   ├── market/     # 集市页面
│   │   ├── chat/       # 聊天列表
│   │   ├── chat-room/  # 聊天室
│   │   ├── communities/# 社区管理
│   │   └── profile/    # 个人主页
│   ├── components/     # 组件
│   │   ├── NavBar.vue     # 导航栏
│   │   ├── Card.vue       # 卡片容器
│   │   └── UserInfo.vue   # 用户信息
│   ├── utils/          # 工具函数
│   │   ├── api.js      # API 接口
│   │   ├── ws.js       # WebSocket 连接
│   │   ├── format.js   # 格式化工具
│   │   ├── storage.js  # 本地存储
│   │   ├── validator.js # 表单验证
│   │   └── toast.js    # 提示工具
│   ├── static/         # 静态资源
│   ├── App.vue         # 应用入口
│   └── main.js         # 主文件
├── pages.json          # 页面配置
├── manifest.json       # 应用配置
├── package.json        # 依赖管理
└── vite.config.js      # Vite 配置
```

## 功能特性

### 已实现功能

1. **用户认证**
   - 手机号 + 密码登录
   - 注册新用户
   - Token 刷新机制

2. **社区管理**
   - 搜索街镇
   - 订阅/取消订阅社区
   - 设置主社区
   - 切换当前社区

3. **广场动态**
   - 发布动态（文字 + 图片）
   - 动态 Feed 流
   - 点赞/取消点赞
   - 评论功能
   - 标签分类（求助/分享/通知/问答）

4. **集市**
   - 发布商品
   - 商品列表浏览
   - 分类筛选
   - 新旧程度标识

5. **聊天功能**
   - 社区频道（公共群聊）
   - 私信
   - WebSocket 实时消息
   - 历史消息加载

6. **个人中心**
   - 查看个人信息
   - 修改昵称和头像
   - 查看我的社区
   - 退出登录

### 待实现功能

1. 微信小程序登录
2. 互助/招工模块
3. 积分/信用体系
4. 消息通知中心
5. 用户主页详情

## 开发指南

### 安装依赖

```bash
cd client/uni-app
npm install
```

### 本地开发

H5 开发：
```bash
npm run dev:h5
```

微信小程序开发：
```bash
npm run dev:mp-weixin
```

### 构建生产版本

H5 构建：
```bash
npm run build:h5
```

微信小程序构建：
```bash
npm run build:mp-weixin
```

## 配置说明

### API 地址

修改 `src/utils/api.js` 中的 `BASE` 常量：

```javascript
const BASE = 'http://localhost:8080/api/v1' // 本地开发
// const BASE = 'https://your-domain.com/api/v1' // 生产环境
```

### WebSocket 地址

修改 `src/utils/ws.js` 中的 `WS_URL` 常量：

```javascript
const WS_URL = 'ws://localhost:8080/ws' // 本地开发
// const WS_URL = 'wss://your-domain.com/ws' // 生产环境
```

### 微信小程序 AppID

修改 `manifest.json` 中的 `appid`：

```json
"mp-weixin": {
  "appid": "your-wechat-appid"
}
```

## API 接口说明

所有 API 接口定义在 `src/utils/api.js` 中，包含：

- 认证接口：`login`, `register`, `refresh`
- 用户接口：`getProfile`, `updateProfile`, `getUser`
- 社区接口：`searchCommunities`, `subscribeCommunity`, `getMyCommunities`
- 帖子接口：`createPost`, `listPosts`, `likePost`, `createComment`
- 集市接口：`createItem`, `listItems`, `markSold`
- 聊天接口：`getConversations`, `getMessages`, `createPrivateConversation`

## 组件使用

### NavBar 导航栏

```vue
<NavBar title="页面标题" :show-back="true">
  <template #right>
    <button>按钮</button>
  </template>
</NavBar>
```

### Card 卡片容器

```vue
<Card @click="handleClick">
  <view>卡片内容</view>
</Card>
```

### UserInfo 用户信息

```vue
<UserInfo
  :user-id="userId"
  :nickname="nickname"
  :avatar-url="avatarUrl"
  :time="createdAt"
  @click="handleUserClick"
/>
```

## 工具函数使用

```javascript
import { formatTime, formatPrice } from '@/utils/format'
import { validatePhone, validatePassword } from '@/utils/validator'
import { showSuccess, showError, showModal } from '@/utils/toast'
import { storage } from '@/utils/storage'

// 格式化时间
const timeStr = formatTime(post.created_at)

// 格式化价格
const priceStr = formatPrice(item.price)

// 验证表单
if (!validatePhone(phone)) {
  showError('请输入正确的手机号')
}

// 显示提示
showSuccess('操作成功')
showError('操作失败')

// 存储管理
storage.setToken(token)
const token = storage.getToken()
```

## 注意事项

1. **Token 管理**：使用 `storage` 工具管理 Token，确保刷新 Token 逻辑正确
2. **图片上传**：当前使用本地文件上传，需要后端支持文件存储服务
3. **WebSocket**：聊天功能依赖 WebSocket，确保后端服务正常运行
4. **权限控制**：所有需要登录的页面都需要验证 Token
5. **错误处理**：统一使用 `toast` 工具显示错误信息

## 下一步计划

1. 接入微信小程序登录
2. 添加图片裁剪和压缩功能
3. 实现消息推送
4. 添加分享功能
5. 优化用户体验和性能

## 技术栈

- **框架**：uni-app + Vue 3
- **构建工具**：Vite
- **HTTP 客户端**：uni.request
- **实时通信**：WebSocket
- **UI 组件**：自定义组件
- **状态管理**：Pinia（可选）