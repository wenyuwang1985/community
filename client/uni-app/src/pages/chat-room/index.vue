<template>
  <view class="chat-page">
    <view class="chat-messages" id="msgBox">
      <view v-for="msg in messages" :key="msg.id" class="msg-row" :class="{me: msg.sender_id == userID}">
        <view>
          <text v-if="msg.sender_id != userID" class="msg-sender">{{ msg.sender?.nickname || '用户' + msg.sender_id }}</text>
          <view class="msg-bubble">{{ msg.content }}</view>
        </view>
      </view>
    </view>
    <view class="input-bar">
      <input v-model="inputText" placeholder="输入消息..." class="chat-input" confirm-type="send" @confirm="send" />
      <button class="chat-send" @click="send">发送</button>
    </view>
  </view>
</template>

<script>
import { api } from '../../utils/api.js'
import { chatSocket } from '../../utils/ws.js'

export default {
  data() {
    return {
      roomName: '', convId: 0, messages: [], inputText: '', userID: 0, unsubscribe: null
    }
  },
  onLoad(options) {
    this.convId = Number(options.id)
    this.roomName = decodeURIComponent(options.name || '聊天')
    this.userID = Number(uni.getStorageSync('userID') || 0)
    if (!this.userID) {
      api.getProfile().then(p => {
        this.userID = p.id
        uni.setStorageSync('userID', p.id)
      }).catch(() => {})
    }
    this.loadHistory()
    this.unsubscribe = chatSocket.on('message', (payload) => {
      const msg = payload.data
      if (msg.conversation_id == this.convId) {
        this.messages.push(msg)
        this.scrollToBottom()
      }
    })
  },
  onUnload() {
    if (this.unsubscribe) this.unsubscribe()
  },
  methods: {
    async loadHistory() {
      try { this.messages = await api.getMessages(this.convId) }
      catch (e) { console.error(e) }
      this.scrollToBottom()
    },
    send() {
      const text = this.inputText.trim()
      if (!text) return
      try {
        chatSocket.sendMessage(this.convId, text)
        this.inputText = ''
      } catch (e) {
        uni.showToast({ title: e.message, icon: 'none' })
      }
    },
    scrollToBottom() {
      this.$nextTick(() => {
        uni.createSelectorQuery().in(this).select('#msgBox').boundingClientRect(rect => {
          if (rect) uni.pageScrollTo({ scrollTop: rect.height, duration: 0 })
        }).exec()
      })
    }
  }
}
</script>

<style scoped>
.chat-page { display: flex; flex-direction: column; height: 100vh; }
.chat-messages { flex: 1; overflow-y: auto; padding: 12px 16px; background: #f5f5f5; }
.msg-row { margin-bottom: 12px; display: flex; }
.msg-row.me { justify-content: flex-end; }
.msg-bubble { max-width: 70%; padding: 10px 14px; border-radius: 12px; font-size: 14px; line-height: 1.5; word-break: break-all; background: #fff; color: #333; }
.msg-row.me .msg-bubble { background: #95ec69; }
.msg-sender { font-size: 11px; color: #999; margin-bottom: 4px; display: block; }
.input-bar { display: flex; padding: 8px 12px; background: #fff; border-top: 1px solid #eee; gap: 8px; }
.chat-input { flex: 1; padding: 8px 12px; border: 1px solid #ddd; border-radius: 20px; font-size: 14px; }
.chat-send { padding: 8px 16px; background: #07c160; color: #fff; border-radius: 20px; font-size: 14px; }
</style>
