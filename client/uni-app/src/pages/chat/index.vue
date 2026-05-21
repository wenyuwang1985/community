<template>
  <view class="page">
    <view v-if="conversations.length === 0" class="empty">暂无会话</view>
    <view v-for="conv in conversations" :key="conv.id" class="conv-item" @click="enterRoom(conv)">
      <image class="avatar" :src="'https://api.dicebear.com/7.x/avataaars/svg?seed=' + conv.id" mode="aspectFill" />
      <view class="conv-info">
        <text class="conv-name">{{ conv.name || (conv.type === 'channel' ? '社区频道' : '私聊') }}</text>
        <text class="conv-preview">{{ conv.type === 'channel' ? '社区公共频道' : '点击开始聊天' }}</text>
      </view>
    </view>
  </view>
</template>

<script>
import { api } from '../../utils/api.js'
import { chatSocket } from '../../utils/ws.js'

export default {
  data() { return { conversations: [] } },
  onShow() {
    this.loadConversations()
    const token = uni.getStorageSync('token')
    if (token) chatSocket.connect()
  },
  methods: {
    async loadConversations() {
      try { this.conversations = await api.getConversations() }
      catch (e) { console.error(e) }
    },
    enterRoom(conv) {
      uni.navigateTo({ url: `/src/pages/chat-room/index?id=${conv.id}&name=${encodeURIComponent(conv.name || (conv.type === 'channel' ? '社区频道' : '私聊'))}` })
    }
  }
}
</script>

<style scoped>
.page { padding-bottom: 20px; }
.conv-item { display: flex; align-items: center; padding: 12px 16px; background: #fff; border-bottom: 1px solid #f5f5f5; }
.avatar { width: 48px; height: 48px; border-radius: 50%; background: #ddd; flex-shrink: 0; }
.conv-info { flex: 1; margin-left: 10px; }
.conv-name { font-size: 15px; color: #333; display: block; }
.conv-preview { font-size: 13px; color: #999; margin-top: 2px; }
.empty { text-align: center; padding: 60px 20px; color: #999; font-size: 14px; }
</style>
