import { api } from '../api.js'
import { chatSocket } from '../ws.js'

export default {
  template: `
    <div class="page">
      <div class="header">聊天</div>
      <div v-if="conversations.length === 0" class="empty">暂无会话</div>
      <div v-for="conv in conversations" :key="conv.id" class="conv-item" @click="enterRoom(conv)">
        <div class="avatar"><img :src="'https://api.dicebear.com/7.x/avataaars/svg?seed=' + conv.id" alt=""></div>
        <div class="conv-info">
          <div class="conv-name">{{ conv.name || (conv.type === 'channel' ? '社区频道' : '私聊') }}</div>
          <div class="conv-preview">{{ conv.type === 'channel' ? '社区公共频道' : '点击开始聊天' }}</div>
        </div>
      </div>

      <div class="tabbar">
        <div class="tabbar-item" @click="$router.push('/feed')"><span class="icon">📰</span>广场</div>
        <div class="tabbar-item" @click="$router.push('/market')"><span class="icon">🛒</span>集市</div>
        <div class="tabbar-item active"><span class="icon">💬</span>聊天</div>
        <div class="tabbar-item" @click="$router.push('/profile')"><span class="icon">👤</span>我的</div>
      </div>
    </div>
  `,
  data() {
    return { conversations: [] }
  },
  async mounted() {
    this.loadConversations()
    // 建立 WebSocket
    const token = localStorage.getItem('token')
    if (token) chatSocket.connect(token)
  },
  methods: {
    async loadConversations() {
      try {
        this.conversations = await api.getConversations()
      } catch (e) { console.error(e) }
    },
    enterRoom(conv) {
      this.$router.push({ path: '/chat-room', query: { id: conv.id, name: conv.name || (conv.type === 'channel' ? '社区频道' : '私聊') } })
    }
  }
}
