import { api } from '../api.js'

export default {
  template: `
    <div class="page">
      <div class="header" style="display:flex;align-items:center;justify-content:space-between;">
        <span @click="$router.back()" style="font-size:20px;">←</span>
        <span>我的社区</span>
        <span style="width:24px;"></span>
      </div>

      <!-- 搜索 -->
      <div style="padding:12px 16px;background:#fff;border-bottom:1px solid #f5f5f5;">
        <div style="display:flex;gap:8px;">
          <input v-model="keyword" placeholder="搜索街镇名称" style="flex:1;padding:8px 12px;border:1px solid #ddd;border-radius:20px;font-size:14px;">
          <button style="padding:8px 16px;background:#07c160;color:#fff;border-radius:20px;font-size:14px;" @click="search">搜索</button>
        </div>
      </div>

      <!-- 搜索结果 -->
      <div v-if="searchResults.length">
        <div style="padding:8px 16px;font-size:14px;color:#999;">搜索结果</div>
        <div v-for="c in searchResults" :key="c.id" class="card" style="display:flex;justify-content:space-between;align-items:center;">
          <div>
            <div style="font-size:15px;color:#333;">{{ c.name }}</div>
            <div style="font-size:12px;color:#999;">{{ c.province }} {{ c.city }} {{ c.district }}</div>
          </div>
          <button style="padding:6px 14px;background:#07c160;color:#fff;border-radius:16px;font-size:13px;" @click="subscribe(c)">订阅</button>
        </div>
      </div>

      <!-- 已订阅 -->
      <div style="padding:8px 16px;font-size:14px;color:#999;">已订阅社区</div>
      <div v-for="c in communities" :key="c.id" class="card" style="display:flex;justify-content:space-between;align-items:center;">
        <div>
          <div style="font-size:15px;color:#333;">{{ c.name }} {{ c.is_primary ? '(主社区)' : '' }}</div>
          <div style="font-size:12px;color:#999;">{{ c.province }} {{ c.city }} {{ c.district }}</div>
        </div>
        <div style="display:flex;gap:8px;">
          <button v-if="!c.is_primary" style="padding:6px 12px;background:#f2f2f2;color:#333;border-radius:16px;font-size:13px;" @click="setPrimary(c)">设为主</button>
          <button style="padding:6px 12px;background:#fff2f0;color:#fa5151;border-radius:16px;font-size:13px;" @click="unsubscribe(c)">取消</button>
        </div>
      </div>

      <div v-if="communities.length === 0" class="empty">暂无订阅社区，搜索并订阅一个吧</div>
    </div>
  `,
  data() {
    return { keyword: '', searchResults: [], communities: [] }
  },
  async mounted() {
    this.loadCommunities()
  },
  methods: {
    async loadCommunities() {
      try {
        this.communities = await api.getMyCommunities()
      } catch (e) { console.error(e) }
    },
    async search() {
      if (!this.keyword.trim()) return
      try {
        this.searchResults = await api.searchCommunities(this.keyword.trim())
      } catch (e) { alert(e.message) }
    },
    async subscribe(c) {
      try {
        await api.subscribeCommunity(c.id)
        alert('订阅成功')
        this.searchResults = []
        this.loadCommunities()
      } catch (e) { alert(e.message) }
    },
    async unsubscribe(c) {
      if (!confirm('确定取消订阅？')) return
      try {
        await api.unsubscribeCommunity(c.id)
        this.loadCommunities()
      } catch (e) { alert(e.message) }
    },
    async setPrimary(c) {
      try {
        await api.setPrimaryCommunity(c.id)
        this.loadCommunities()
      } catch (e) { alert(e.message) }
    }
  }
}
