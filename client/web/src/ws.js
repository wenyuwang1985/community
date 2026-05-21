const WS_URL = 'ws://localhost:8080/ws'

class ChatSocket {
  constructor() {
    this.ws = null
    this.listeners = new Map()
    this.reconnectTimer = null
    this.userID = null
  }

  connect(token) {
    if (this.ws) return
    this.ws = new WebSocket(`${WS_URL}?token=${token}`)
    this.ws.onopen = () => {
      console.log('WebSocket 连接成功')
    }
    this.ws.onmessage = (e) => {
      const msg = JSON.parse(e.data)
      this.emit(msg.type, msg)
    }
    this.ws.onclose = () => {
      this.ws = null
      this.reconnect(token)
    }
    this.ws.onerror = (err) => {
      console.error('WebSocket 错误', err)
    }
  }

  reconnect(token) {
    if (this.reconnectTimer) return
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null
      this.connect(token)
    }, 3000)
  }

  sendMessage(conversation_id, content) {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error('WebSocket 未连接')
    }
    this.ws.send(JSON.stringify({ action: 'send_message', conversation_id, content }))
  }

  ping() {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ action: 'ping' }))
    }
  }

  on(type, callback) {
    if (!this.listeners.has(type)) this.listeners.set(type, [])
    this.listeners.get(type).push(callback)
    return () => {
      const arr = this.listeners.get(type)
      const idx = arr.indexOf(callback)
      if (idx > -1) arr.splice(idx, 1)
    }
  }

  emit(type, data) {
    const arr = this.listeners.get(type) || []
    arr.forEach(cb => cb(data))
  }

  close() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }
}

export const chatSocket = new ChatSocket()
