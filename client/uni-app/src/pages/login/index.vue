<template>
  <view class="page">
    <view style="text-align:center;margin:80px 0 40px;">
      <text style="color:#07c160;font-size:24px;font-weight:600;">社区生活</text>
    </view>
    <view style="padding:0 24px;">
      <input v-model="phone" placeholder="手机号" maxlength="11" type="number" class="input" />
      <input v-model="password" placeholder="密码" password class="input" />
      <button class="btn-primary" @click="login" :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
      <view style="text-align:center;margin-top:16px;">
        <text style="color:#576b95;font-size:14px;" @click="goRegister">还没有账号？去注册</text>
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
    async login() {
      if (!this.phone || !this.password) { this.error = '请填写手机号和密码'; return }
      this.loading = true; this.error = ''
      try {
        const data = await api.login(this.phone, this.password)
        uni.setStorageSync('token', data.access_token)
        uni.setStorageSync('refresh_token', data.refresh_token)
        uni.switchTab({ url: '/src/pages/index/index' })
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    },
    goRegister() {
      uni.navigateTo({ url: '/src/pages/register/index' })
    }
  }
}
</script>

<style scoped>
.page { min-height: 100vh; background: #fff; }
.input { width: 100%; padding: 12px 0; border-bottom: 1px solid #eee; font-size: 15px; margin-bottom: 12px; }
.btn-primary { width: 100%; padding: 12px; background: #07c160; color: #fff; border-radius: 6px; font-size: 16px; margin-top: 24px; }
</style>
