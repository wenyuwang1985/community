<template>
  <view class="user-info" @click="$emit('click')">
    <image
      class="avatar"
      :src="avatarUrl || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + userId"
      mode="aspectFill"
    />
    <view class="info">
      <text class="nickname">{{ nickname || '用户' + userId }}</text>
      <text class="time" v-if="time">{{ formatTime(time) }}</text>
    </view>
  </view>
</template>

<script>
export default {
  props: {
    userId: { type: Number, default: 0 },
    nickname: { type: String, default: '' },
    avatarUrl: { type: String, default: '' },
    time: { type: String, default: '' }
  },
  methods: {
    formatTime(t) {
      const d = new Date(t)
      const now = new Date()
      const diff = now - d

      if (diff < 60000) return '刚刚'
      if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
      if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
      if (diff < 604800000) return Math.floor(diff / 86400000) + '天前'

      return `${d.getMonth() + 1}月${d.getDate()}日`
    }
  }
}
</script>

<style scoped>
.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #ddd;
  flex-shrink: 0;
}

.info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.nickname {
  font-size: 15px;
  font-weight: 500;
  color: #333;
}

.time {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
}
</style>