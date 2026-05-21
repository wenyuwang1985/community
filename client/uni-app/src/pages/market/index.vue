<template>
  <view class="page">
    <view class="community-bar" @click="showPicker = true">
      <text>{{ currentCommunity?.name || '选择社区' }}</text>
      <text style="margin-left:4px;">▼</text>
    </view>

    <view v-if="showPicker" class="picker-mask" @click.self="showPicker = false">
      <view class="picker-content">
        <view v-for="c in communities" :key="c.id" class="picker-item" @click="switchCommunity(c)">
          <text>{{ c.name }}</text>
          <text v-if="c.id === currentCommunity?.id" style="color:#07c160;">✓</text>
        </view>
      </view>
    </view>

    <!-- 分类筛选 -->
    <scroll-view scroll-x style="white-space:nowrap;padding:10px 16px;background:#fff;border-bottom:1px solid #f5f5f5;">
      <text v-for="cat in categories" :key="cat.key" @click="switchCategory(cat.key)"
        style="display:inline-block;padding:4px 12px;border-radius:12px;font-size:13px;margin-right:8px;"
        :style="currentCategory===cat.key ? 'background:#07c160;color:#fff;' : 'background:#f2f2f2;color:#666;'">
        {{ cat.label }}
      </text>
    </scroll-view>

    <!-- 发布弹窗 -->
    <view v-if="showItem" class="picker-mask" @click.self="showItem = false">
      <view class="picker-content">
        <input v-model="itemForm.title" placeholder="标题" style="padding:8px 0;border-bottom:1px solid #eee;margin-bottom:8px;" />
        <input v-model.number="itemForm.price" placeholder="价格（分）" type="number" style="padding:8px 0;border-bottom:1px solid #eee;margin-bottom:8px;" />
        <picker mode="selector" :range="conditionLabels" :value="conditionIndex" @change="onConditionChange">
          <view style="padding:8px 0;border-bottom:1px solid #eee;margin-bottom:8px;">新旧：{{ conditionLabels[conditionIndex] }}</view>
        </picker>
        <picker mode="selector" :range="categoryLabels" :value="categoryIndex" @change="onCategoryChange">
          <view style="padding:8px 0;border-bottom:1px solid #eee;margin-bottom:8px;">分类：{{ categoryLabels[categoryIndex] }}</view>
        </picker>
        <input v-model="itemForm.images" placeholder="图片URL，逗号分隔" style="padding:8px 0;border-bottom:1px solid #eee;margin-bottom:8px;" />
        <button class="btn-primary" @click="submitItem">发布</button>
      </view>
    </view>

    <!-- 列表 -->
    <view v-if="items.length === 0 && !loading" class="empty">暂无商品</view>
    <view v-for="item in items" :key="item.id" class="item-card">
      <image class="item-image" :src="item.images?.[0] || 'https://via.placeholder.com/80?text=No+Img'" mode="aspectFill" />
      <view class="item-info">
        <text class="item-title">{{ item.title }}</text>
        <view class="item-meta">
          <text class="item-price">¥{{ (item.price / 100).toFixed(2) }}</text>
          <text class="item-condition">{{ conditionText(item.condition) }}</text>
        </view>
        <text style="font-size:12px;color:#999;">{{ item.seller?.nickname || '用户' + item.seller_id }}</text>
      </view>
    </view>

    <view v-if="loading" class="loading">加载中...</view>
    <view v-if="hasMore && items.length > 0" style="text-align:center;padding:16px;color:#999;font-size:13px;" @click="loadMore">点击加载更多</view>

    <button class="fab" @click="showItem = true">+</button>
  </view>
</template>

<script>
import { api } from '../../utils/api.js'

const categories = [
  { key: '', label: '全部' }, { key: 'appliance', label: '家电' },
  { key: 'furniture', label: '家具' }, { key: 'book', label: '书籍' },
  { key: 'baby', label: '母婴' }, { key: 'sports', label: '运动' },
  { key: 'other', label: '其他' }
]
const conditions = [
  { key: 'new', label: '全新' }, { key: 'like_new', label: '几乎全新' },
  { key: 'lightly_used', label: '轻微使用' }, { key: 'heavily_used', label: '明显使用' }
]

export default {
  data() {
    return {
      communities: [], currentCommunity: null,
      categories, currentCategory: '',
      items: [], loading: false, hasMore: true, lastId: 0,
      showPicker: false, showItem: false,
      itemForm: { title: '', price: 0, condition: 'like_new', category: 'other', images: '' },
      conditionIndex: 1, categoryIndex: 5,
      conditionLabels: conditions.map(c => c.label),
      categoryLabels: categories.filter(c => c.key).map(c => c.label)
    }
  },
  onShow() {
    this.loadCommunities().then(() => {
      if (this.currentCommunity && this.items.length === 0) this.loadItems()
    })
  },
  methods: {
    async loadCommunities() {
      try {
        this.communities = await api.getMyCommunities()
        this.currentCommunity = this.communities.find(c => c.is_primary) || this.communities[0]
      } catch (e) { console.error(e) }
    },
    switchCommunity(c) { this.currentCommunity = c; this.showPicker = false; this.resetAndLoad() },
    switchCategory(cat) { this.currentCategory = cat; this.resetAndLoad() },
    resetAndLoad() { this.items = []; this.lastId = 0; this.hasMore = true; this.loadItems() },
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
    onConditionChange(e) { this.conditionIndex = e.detail.value; this.itemForm.condition = conditions[this.conditionIndex].key },
    onCategoryChange(e) { this.categoryIndex = e.detail.value; this.itemForm.category = categories.filter(c=>c.key)[this.categoryIndex].key },
    async submitItem() {
      if (!this.itemForm.title || this.itemForm.price <= 0) { uni.showToast({ title: '请填写标题和价格', icon: 'none' }); return }
      try {
        const images = this.itemForm.images.split(',').map(s => s.trim()).filter(Boolean)
        await api.createItem({ ...this.itemForm, community_id: this.currentCommunity.id, images })
        this.showItem = false
        this.itemForm = { title: '', price: 0, condition: 'like_new', category: 'other', images: '' }
        this.resetAndLoad()
      } catch (e) { uni.showToast({ title: e.message, icon: 'none' }) }
    },
    conditionText(c) {
      const map = { new: '全新', like_new: '几乎全新', lightly_used: '轻微使用', heavily_used: '明显使用' }
      return map[c] || c
    }
  }
}
</script>

<style scoped>
.page { padding-bottom: 20px; }
.community-bar { padding: 10px 16px; background: #fff; border-bottom: 1px solid #eee; display: flex; align-items: center; font-size: 15px; }
.picker-mask { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 200; display: flex; align-items: flex-end; }
.picker-content { background: #fff; width: 100%; border-radius: 12px 12px 0 0; padding: 16px; max-height: 60vh; overflow-y: auto; }
.picker-item { padding: 12px 0; border-bottom: 1px solid #f5f5f5; display: flex; justify-content: space-between; font-size: 15px; }
.item-card { display: flex; padding: 12px 16px; background: #fff; border-bottom: 1px solid #f5f5f5; gap: 10px; }
.item-image { width: 80px; height: 80px; border-radius: 6px; background: #f5f5f5; flex-shrink: 0; }
.item-info { flex: 1; display: flex; flex-direction: column; justify-content: space-between; }
.item-title { font-size: 15px; color: #333; line-height: 1.4; overflow: hidden; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; }
.item-meta { display: flex; justify-content: space-between; align-items: center; margin-top: 6px; }
.item-price { font-size: 16px; color: #fa5151; font-weight: 600; }
.item-condition { font-size: 12px; color: #999; }
.empty { text-align: center; padding: 60px 20px; color: #999; font-size: 14px; }
.loading { text-align: center; padding: 20px; color: #999; }
.fab { position: fixed; bottom: 30px; right: 16px; width: 50px; height: 50px; border-radius: 50%; background: #07c160; color: #fff; font-size: 28px; display: flex; align-items: center; justify-content: center; box-shadow: 0 2px 8px rgba(0,0,0,0.15); z-index: 99; }
.btn-primary { background: #07c160; color: #fff; border-radius: 6px; font-size: 16px; }
</style>
