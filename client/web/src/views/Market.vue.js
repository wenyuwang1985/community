import { api, uploadFiles } from '../api.js'

const categories = [
  { key: '', label: '全部' },
  { key: 'appliance', label: '家电' },
  { key: 'furniture', label: '家具' },
  { key: 'book', label: '书籍' },
  { key: 'baby', label: '母婴' },
  { key: 'sports', label: '运动' },
  { key: 'other', label: '其他' }
]

const conditions = [
  { key: 'new', label: '全新' },
  { key: 'like_new', label: '几乎全新' },
  { key: 'lightly_used', label: '轻微使用' },
  { key: 'heavily_used', label: '明显使用' }
]

export default {
  template: `
    <div class="page">
      <div class="header">
        <span>集市 - {{ currentCommunity?.name || '选择社区' }}</span>
        <span style="position:absolute;right:16px;top:12px;font-size:14px;" @click="showCommunitySwitch=true">切换</span>
      </div>

      <!-- 社区选择 -->
      <div v-if="showCommunitySwitch" class="modal-overlay" @click.self="showCommunitySwitch=false">
        <div class="modal-content">
          <div class="modal-header">选择社区</div>
          <div v-for="c in communities" :key="c.id" class="form-group" style="display:flex;justify-content:space-between;align-items:center;cursor:pointer;" @click="switchCommunity(c)">
            <span>{{ c.name }} {{ c.is_primary ? '(主)' : '' }}</span>
            <span v-if="c.id === currentCommunity?.id" style="color:#07c160;">✓</span>
          </div>
        </div>
      </div>

      <!-- 分类筛选 -->
      <div style="display:flex;gap:8px;padding:10px 16px;background:#fff;overflow-x:auto;white-space:nowrap;border-bottom:1px solid #f5f5f5;">
        <span v-for="cat in categories" :key="cat.key" @click="switchCategory(cat.key)"
          style="padding:4px 12px;border-radius:12px;font-size:13px;cursor:pointer;"
          :style="currentCategory===cat.key ? 'background:#07c160;color:#fff;' : 'background:#f2f2f2;color:#666;'">
          {{ cat.label }}
        </span>
      </div>

      <!-- 发布商品弹层 -->
      <div v-if="showItemModal" class="modal-overlay" @click.self="showItemModal=false">
        <div class="modal-content">
          <div class="modal-header">发布商品</div>
          <input class="form-input" v-model="itemForm.title" placeholder="标题" style="margin-bottom:8px;">
          <input class="form-input" v-model.number="itemForm.price" placeholder="价格（分）" type="number" style="margin-bottom:8px;">
          <select v-model="itemForm.condition" style="width:100%;padding:8px;margin-bottom:8px;border:1px solid #ddd;border-radius:4px;">
            <option v-for="c in conditions" :value="c.key">{{ c.label }}</option>
          </select>
          <select v-model="itemForm.category" style="width:100%;padding:8px;margin-bottom:8px;border:1px solid #ddd;border-radius:4px;">
            <option v-for="c in categories.filter(x=>x.key)" :value="c.key">{{ c.label }}</option>
          </select>
          <input class="form-input" v-model="itemForm.images" placeholder="图片URL，逗号分隔" style="margin-bottom:8px;">
          <button class="btn btn-primary" @click="submitItem" :disabled="itemLoading">发布</button>
        </div>
      </div>

      <!-- 商品列表 -->
      <div v-if="items.length === 0 && !loading" class="empty">暂无商品</div>
      <div v-for="item in items" :key="item.id" class="item-card" @click="viewItem(item)">
        <div class="item-image"><img :src="item.images?.[0] || 'https://via.placeholder.com/80?text=No+Img'" alt=""></div>
        <div class="item-info">
          <div class="item-title">{{ item.title }}</div>
          <div class="item-meta">
            <span class="item-price">¥{{ (item.price / 100).toFixed(2) }}</span>
            <span class="item-condition">{{ conditionText(item.condition) }}</span>
          </div>
          <div style="font-size:12px;color:#999;">{{ item.seller?.nickname || '用户' + item.seller_id }}</div>
        </div>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-if="hasMore && items.length > 0" style="text-align:center;padding:16px;color:#999;font-size:13px;" @click="loadMore">点击加载更多</div>

      <button class="fab" @click="showItemModal=true">+</button>

      <div class="tabbar">
        <div class="tabbar-item" @click="$router.push('/feed')"><span class="icon">📰</span>广场</div>
        <div class="tabbar-item active"><span class="icon">🛒</span>集市</div>
        <div class="tabbar-item" @click="$router.push('/chat')"><span class="icon">💬</span>聊天</div>
        <div class="tabbar-item" @click="$router.push('/profile')"><span class="icon">👤</span>我的</div>
      </div>
    </div>
  `,
  data() {
    return {
      communities: [],
      currentCommunity: null,
      categories,
      conditions,
      currentCategory: '',
      items: [],
      loading: false,
      hasMore: true,
      lastId: 0,
      showCommunitySwitch: false,
      showItemModal: false,
      itemForm: { title: '', price: 0, condition: 'like_new', category: 'other', images: '' },
      selectedFiles: [],
      previewImages: [],
      itemLoading: false,
    }
  },
  async mounted() {
    await this.loadCommunities()
    if (this.currentCommunity) this.loadItems()
  },
  methods: {
    async loadCommunities() {
      try {
        this.communities = await api.getMyCommunities()
        this.currentCommunity = this.communities.find(c => c.is_primary) || this.communities[0]
      } catch (e) { console.error(e) }
    },
    switchCommunity(c) {
      this.currentCommunity = c
      this.showCommunitySwitch = false
      this.items = []; this.lastId = 0; this.hasMore = true
      this.loadItems()
    },
    switchCategory(cat) {
      this.currentCategory = cat
      this.items = []; this.lastId = 0; this.hasMore = true
      this.loadItems()
    },
    async loadItems() {
      if (!this.currentCommunity || this.loading) return
      this.loading = true
      try {
        const list = await api.listItems(this.currentCommunity.id, this.currentCategory, this.lastId)
        if (list.length < 10) this.hasMore = false
        if (list.length > 0) this.lastId = list[list.length - 1].id
        this.items.push(...list)
      } catch (e) { console.error(e) }
      finally { this.loading = false }
    },
    loadMore() { this.loadItems() },
    async submitItem() {
      if (!this.itemForm.title || this.itemForm.price <= 0) { alert('请填写标题和价格'); return }
      this.itemLoading = true
      try {
        const images = this.itemForm.images.split(',').map(s => s.trim()).filter(Boolean)
        await api.createItem({ ...this.itemForm, community_id: this.currentCommunity.id, images })
        this.showItemModal = false
        this.itemForm = { title: '', price: 0, condition: 'like_new', category: 'other', images: '' }
        this.items = []; this.lastId = 0; this.hasMore = true
        this.loadItems()
      } catch (e) { alert(e.message) }
      finally { this.itemLoading = false }
    },
    viewItem(item) {
      alert(`商品：${item.title}\n价格：¥${(item.price/100).toFixed(2)}\n卖家：${item.seller?.nickname || '用户'+item.seller_id}\n\n点击“联系卖家”可通过聊天私信`)
    },
    conditionText(c) {
      const map = { new: '全新', like_new: '几乎全新', lightly_used: '轻微使用', heavily_used: '明显使用' }
      return map[c] || c
    }
  }
}
