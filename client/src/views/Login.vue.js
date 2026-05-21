import { api } from '../api.js'

export default {
  template: `
    <div class="page" style="display:flex;flex-direction:column;justify-content:center;padding:40px 24px;">
      <h2 style="text-align:center;margin-bottom:40px;color:#07c160;font-size:24px;">社区生活</h2>
      <div class="form-group">
        <input class="form-input" v-model="phone" placeholder="手机号" maxlength="11" type="tel">
      </div>
      <div class="form-group">
        <input class="form-input" v-model="password" placeholder="密码" type="password">
      </div>
      <button class="btn btn-primary" style="margin-top:24px;" @click="login" :disabled="loading">
        {{ loading ? '登录中...' : '登录' }}
      </button>
      <div style="text-align:center;margin-top:16px;font-size:14px;color:#576b95;" @click="$router.push('/register')">
        还没有账号？去注册
      </div>
      <div v-if="error" style="color:#fa5151;text-align:center;margin-top:12px;font-size:13px;">{{ error }}</div>
    </div>
  `,
  data() {
    return { phone: '', password: '', loading: false, error: '' }
  },
  methods: {
    async login() {
      if (!this.phone || !this.password) { this.error = '请填写手机号和密码'; return }
      this.loading = true; this.error = ''
      try {
        const data = await api.login(this.phone, this.password)
        localStorage.setItem('token', data.access_token)
        localStorage.setItem('refresh_token', data.refresh_token)
        this.$router.replace('/feed')
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
    }
  }
}
