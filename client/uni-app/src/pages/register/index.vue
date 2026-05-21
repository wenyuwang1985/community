<template>
  <view class="page">
    <view style="text-align:center;margin:80px 0 40px;">
      <text style="color:#07c160;font-size:24px;font-weight:600;">注册账号</text>
    </view>
    <view style="padding:0 24px;">
      <input v-model="phone" placeholder="手机号" maxlength="11" type="number" class="input" />
      <input v-model="password" placeholder="密码（至少6位）" password class="input" />
      <button class="btn-primary" @click="register" :disabled="loading">{{ loading ? '注册中...' : '注册' }}</button>
      <view style="text-align:center;margin-top:16px;">
        <text style="color:#576b95;font-size:14px;" @click="goLogin">已有账号？去登录</text>
      </view>
      <view v-if="error" style="color:#fa5151;text-align:center;margin-top:12px;font-size:13px;">{{ error }}</view>
    </view>
  </view>
</template>

<script>
import { api } from '../../utils/api.js'

export default {
  data() {
    return { phone: '', password: '', loading: false, error: '' }
  },
  methods: {
    async register() {
      if (!this.phone || this.password.length < 6) { this.error = '请填写手机号和密码（至少6位）'; return }
      this.loading = true; this.error = ''
      try {
        await api.register(this.phone, this.password)
        uni.navigateBack()
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    goLogin() {
      uni.navigateBack()
    }
  }
}
</script>

<style scoped>
.page { min-height: 100vh; background: #fff; }
.input { width: 100%; padding: 12px 0; border-bottom: 1px solid #eee; font-size: 15px; margin-bottom: 12px; }
.btn-primary { width: 100%; padding: 12px; background: #07c160; color: #fff; border-radius: 6px; font-size: 16px; margin-top: 24px; }
</style>
