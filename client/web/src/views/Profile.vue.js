import { api } from '../api.js'

export default {
  template: `
    <div class="page">
      <div class="header">我的</div>
      <div class="card" style="display:flex;align-items:center;gap:12px;">
        <div class="avatar" style="width:60px;height:60px;"><img :src="profile.avatar_url || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + profile.id" alt=""></div>
        <div>
          <div style="font-size:17px;font-weight:600;">{{ profile.nickname || '未设置昵称' }}</div>
          <div style="font-size:13px;color:#999;margin-top:4px;">手机号：{{ profile.phone }}</div>
          <div style="font-size:13px;color:#999;">信用分：{{ profile.credit_score }}</div>
        </div>
      </div>

      <div class="card">
        <div class="form-group">
          <span class="form-label">修改昵称</span>
          <input class="form-input" v-model="editNickname" placeholder="输入新昵称">
        </div>
        <div class="form-group">
          <span class="form-label">头像URL</span>
          <input class="form-input" v-model="editAvatar" placeholder="输入图片URL">
        </div>
        <button class="btn btn-primary" style="margin-top:8px;" @click="saveProfile">保存资料</button>
      </div>

      <div class="card" style="padding:0;">
        <div style="padding:14px 16px;border-bottom:1px solid #f5f5f5;display:flex;justify-content:space-between;align-items:center;" @click="$router.push('/communities')">
          <span>我的社区</span>
          <span style="color:#999;">></span>
        </div>
      </div>

      <button class="btn btn-danger" style="margin:20px 16px;" @click="logout">退出登录</button>

      <div class="tabbar">
        <div class="tabbar-item" @click="$router.push('/feed')"><span class="icon">📰</span>广场</div>
        <div class="tabbar-item" @click="$router.push('/market')"><span class="icon">🛒</span>集市</div>
        <div class="tabbar-item" @click="$router.push('/chat')"><span class="icon">💬</span>聊天</div>
        <div class="tabbar-item active"><span class="icon">👤</span>我的</div>
      </div>
    </div>
  `,
  data() {
    return {
      profile: {},
      editNickname: '',
      editAvatar: '',
    }
  },
  async mounted() {
    this.loadProfile()
  },
  methods: {
    async loadProfile() {
      try {
        this.profile = await api.getProfile()
        this.editNickname = this.profile.nickname || ''
        this.editAvatar = this.profile.avatar_url || ''
        localStorage.setItem('userID', this.profile.id)
      } catch (e) { console.error(e) }
    },
    async saveProfile() {
      try {
        await api.updateProfile(this.editNickname, this.editAvatar)
        alert('保存成功')
        this.loadProfile()
      } catch (e) { alert(e.message) }
    },
    logout() {
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('userID')
      this.$router.replace('/login')
    }
  }
}
