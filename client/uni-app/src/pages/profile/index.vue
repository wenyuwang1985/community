<template>
  <view class="page">
    <view class="card" style="display:flex;align-items:center;gap:12px;">
      <image class="avatar" :src="profile.avatar_url || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + profile.id" mode="aspectFill" />
      <view>
        <text style="font-size:17px;font-weight:600;display:block;">{{ profile.nickname || '未设置昵称' }}</text>
        <text style="font-size:13px;color:#999;display:block;margin-top:4px;">手机号：{{ profile.phone }}</text>
        <text style="font-size:13px;color:#999;display:block;">信用分：{{ profile.credit_score }}</text>
      </view>
    </view>

    <view class="card">
      <view style="margin-bottom:8px;">
        <text style="font-size:14px;color:#666;display:block;margin-bottom:6px;">修改昵称</text>
        <input v-model="editNickname" placeholder="输入新昵称" style="width:100%;padding:8px 0;border-bottom:1px solid #eee;" />
      </view>
      <view style="margin-bottom:8px;">
        <text style="font-size:14px;color:#666;display:block;margin-bottom:6px;">头像URL</text>
        <input v-model="editAvatar" placeholder="输入图片URL" style="width:100%;padding:8px 0;border-bottom:1px solid #eee;" />
      </view>
      <button class="btn-primary" style="margin-top:8px;" @click="saveProfile">保存资料</button>
    </view>

    <view class="card" style="padding:0;">
      <view style="padding:14px 16px;border-bottom:1px solid #f5f5f5;display:flex;justify-content:space-between;align-items:center;" @click="goCommunities">
        <text>我的社区</text>
        <text style="color:#999;">></text>
      </view>
    </view>

    <button style="margin:20px 16px;background:#fa5151;color:#fff;border-radius:6px;font-size:16px;" @click="logout">退出登录</button>
  </view>
</template>

<script>
import { api } from '../../utils/api.js'

export default {
  data() {
    return { profile: {}, editNickname: '', editAvatar: '' }
  },
  onShow() { this.loadProfile() },
  methods: {
    async loadProfile() {
      try {
        this.profile = await api.getProfile()
        this.editNickname = this.profile.nickname || ''
        this.editAvatar = this.profile.avatar_url || ''
        uni.setStorageSync('userID', this.profile.id)
      } catch (e) { console.error(e) }
    },
    async saveProfile() {
      try {
        await api.updateProfile(this.editNickname, this.editAvatar)
        uni.showToast({ title: '保存成功' })
        this.loadProfile()
      } catch (e) { uni.showToast({ title: e.message, icon: 'none' }) }
    },
    goCommunities() {
      uni.navigateTo({ url: '/src/pages/communities/index' })
    },
    logout() {
      uni.removeStorageSync('token')
      uni.removeStorageSync('refresh_token')
      uni.removeStorageSync('userID')
      uni.reLaunch({ url: '/src/pages/login/index' })
    }
  }
}
</script>

<style scoped>
.page { padding-bottom: 20px; }
.card { background: #fff; padding: 12px 16px; margin-bottom: 8px; }
.avatar { width: 60px; height: 60px; border-radius: 50%; background: #ddd; flex-shrink: 0; }
.btn-primary { background: #07c160; color: #fff; border-radius: 6px; font-size: 16px; }
</style>
