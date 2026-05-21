const WS_URL = 'ws://localhost:8080/ws'

class ChatSocket {
  constructor() {
    this.socket = null
    this.listeners = new Map()
    this.reconnectTimer = null
  }

  connect() {
    const token = uni.getStorageSync('token')
    if (!token) return
    if (this.socket) return

    this.socket = uni.connectSocket({
      url: `${WS_URL}?token=${token}`,
      complete: () => {}
    })

    this.socket.onOpen(() => {
      console.log('WebSocket 连接成功')
    })

    this.socket.onMessage((res) => {
      try {
        const msg = JSON.parse(res.data)
        this.emit(msg.type, msg)
      } catch (e) {
        console.error('WebSocket 消息解析失败', e)
      }
    })

    this.socket.onClose(() => {
      this.socket = null
      this.reconnect()
    })

    this.socket.onError((err) => {
      console.error('WebSocket 错误', err)
    })
  }

  reconnect() {
    if (this.reconnectTimer) return
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null
      this.connect()
    }, 3000)
  }

  sendMessage(conversation_id, content) {
    if (!this.socket) {
      throw new Error('WebSocket 未连接')
    }
    this.socket.send({
      data: JSON.stringify({ action: 'send_message', conversation_id, content })
    })
  }

  ping() {
    if (this.socket) {
      this.socket.send({
        data: JSON.stringify({ action: 'ping' })
      })
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
    if (this.socket) {
      this.socket.close()
      this.socket = null
    }
  }
}

export const chatSocket = new ChatSocket()
