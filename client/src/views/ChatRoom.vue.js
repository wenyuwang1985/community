import { api } from '../api.js'
import { chatSocket } from '../ws.js'

export default {
  template: `
    <div class="chat-page">
      <div class="header" style="display:flex;align-items:center;justify-content:space-between;">
        <span @click="$router.back()" style="font-size:20px;">←</span>
        <span>{{ roomName }}</span>
        <span style="width:24px;"></span>
      </div>
      <div class="chat-messages" ref="msgBox">
        <div v-for="msg in messages" :key="msg.id" class="msg-row" :class="{me: msg.sender_id == userID}">
          <div>
            <div class="msg-sender" v-if="msg.sender_id != userID">{{ msg.sender?.nickname || '用户' + msg.sender_id }}</div>
            <div class="msg-bubble">{{ msg.content }}</div>
          </div>
        </div>
      </div>
      <div class="chat-input-bar">
        <input class="chat-input" v-model="inputText" @keyup.enter="send" placeholder="输入消息...">
        <button class="chat-send" @click="send">发送</button>
      </div>
    </div>
  `,
  data() {
    return {
      roomName: '',
      convId: 0,
      messages: [],
      inputText: '',
      userID: 0,
      unsubscribe: null,
    }
  },
  async mounted() {
    this.convId = Number(this.$route.query.id)
    this.roomName = this.$route.query.name || '聊天'
    this.userID = Number(localStorage.getItem('userID') || 0)
    if (!this.userID) {
      try {
        const profile = await api.getProfile()
        this.userID = profile.id
        localStorage.setItem('userID', profile.id)
      } catch (e) { console.error(e) }
    }
    await this.loadHistory()
    this.unsubscribe = chatSocket.on('message', (payload) => {
      const msg = payload.data
      if (msg.conversation_id == this.convId) {
        this.messages.push(msg)
        this.scrollToBottom()
      }
    })
    this.scrollToBottom()
  },
  beforeUnmount() {
    if (this.unsubscribe) this.unsubscribe()
  },
  methods: {
    async loadHistory() {
      try {
        const list = await api.getMessages(this.convId)
        this.messages = list
      } catch (e) { console.error(e) }
    },
    send() {
      const text = this.inputText.trim()
      if (!text) return
      try {
        chatSocket.sendMessage(this.convId, text)
        this.inputText = ''
      } catch (e) {
        alert(e.message)
      }
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const box = this.$refs.msgBox
        if (box) box.scrollTop = box.scrollHeight
      })
    }
  }
}
