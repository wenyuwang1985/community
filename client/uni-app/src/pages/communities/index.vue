<template>
  <view class="page">
    <view style="padding:12px 16px;background:#fff;border-bottom:1px solid #f5f5f5;">
      <view style="display:flex;gap:8px;">
        <input v-model="keyword" placeholder="搜索街镇名称" style="flex:1;padding:8px 12px;border:1px solid #ddd;border-radius:20px;font-size:14px;" />
        <button style="padding:8px 16px;background:#07c160;color:#fff;border-radius:20px;font-size:14px;" @click="search">搜索</button>
      </view>
    </view>

    <view v-if="searchResults.length">
      <view style="padding:8px 16px;font-size:14px;color:#999;">搜索结果</view>
      <view v-for="c in searchResults" :key="c.id" class="card" style="display:flex;justify-content:space-between;align-items:center;">
        <view>
          <text style="font-size:15px;color:#333;display:block;">{{ c.name }}</text>
          <text style="font-size:12px;color:#999;">{{ c.province }} {{ c.city }} {{ c.district }}</text>
        </view>
        <button style="padding:6px 14px;background:#07c160;color:#fff;border-radius:16px;font-size:13px;" @click="subscribe(c)">订阅</button>
      </view>
    </view>

    <view style="padding:8px 16px;font-size:14px;color:#999;">已订阅社区</view>
    <view v-for="c in communities" :key="c.id" class="card" style="display:flex;justify-content:space-between;align-items:center;">
      <view>
        <text style="font-size:15px;color:#333;display:block;">{{ c.name }} {{ c.is_primary ? '(主社区)' : '' }}</text>
        <text style="font-size:12px;color:#999;">{{ c.province }} {{ c.city }} {{ c.district }}</text>
      </view>
      <view style="display:flex;gap:8px;">
        <button v-if="!c.is_primary" style="padding:6px 12px;background:#f2f2f2;color:#333;border-radius:16px;font-size:13px;" @click="setPrimary(c)">设为主</button>
        <button style="padding:6px 12px;background:#fff2f0;color:#fa5151;border-radius:16px;font-size:13px;" @click="unsubscribe(c)">取消</button>
      </view>
    </view>

    <view v-if="communities.length === 0" class="empty">暂无订阅社区，搜索并订阅一个吧</view>
  </view>
</template>

<script>
import { api } from '../../utils/api.js'

export default {
  data() { return { keyword: '', searchResults: [], communities: [] } },
  onShow() { this.loadCommunities() },
  methods: {
    async loadCommunities() {
      try { this.communities = await api.getMyCommunities() }
      catch (e) { console.error(e) }
    },
    async search() {
      if (!this.keyword.trim()) return
      try { this.searchResults = await api.searchCommunities(this.keyword.trim()) }
      catch (e) { uni.showToast({ title: e.message, icon: 'none' }) }
    },
    async subscribe(c) {
      try {
        await api.subscribeCommunity(c.id)
        uni.showToast({ title: '订阅成功' })
        this.searchResults = []
        this.loadCommunities()
      } catch (e) { uni.showToast({ title: e.message, icon: 'none' }) }
    },
    async unsubscribe(c) {
      uni.showModal({
        title: '确认',
        content: '确定取消订阅？',
        success: async (res) => {
          if (res.confirm) {
            try { await api.unsubscribeCommunity(c.id); this.loadCommunities() }
            catch (e) { uni.showToast({ title: e.message, icon: 'none' }) }
          }
        }
      })
    },
    async setPrimary(c) {
      try { await api.setPrimaryCommunity(c.id); this.loadCommunities() }
      catch (e) { uni.showToast({ title: e.message, icon: 'none' }) }
    }
  }
}
</script>

<style scoped>
.page { padding-bottom: 20px; }
.card { background: #fff; padding: 12px 16px; margin-bottom: 1px; }
.empty { text-align: center; padding: 60px 20px; color: #999; font-size: 14px; }
</style>
