<template>
  <view class="page">
    <!-- 社区切换 -->
    <view class="community-bar" @click="showPicker = true">
      <text>{{ currentCommunity?.name || '选择社区' }}</text>
      <text style="margin-left:4px;">▼</text>
    </view>

    <!-- 社区选择弹窗 -->
    <view v-if="showPicker" class="picker-mask" @click.self="showPicker = false">
      <view class="picker-content">
        <view v-for="c in communities" :key="c.id" class="picker-item" @click="switchCommunity(c)">
          <text>{{ c.name }} {{ c.is_primary ? '(主)' : '' }}</text>
          <text v-if="c.id === currentCommunity?.id" style="color:#07c160;">✓</text>
        </view>
        <button style="margin-top:12px;background:#f2f2f2;color:#333;font-size:14px;" @click="goCommunities">管理我的社区</button>
      </view>
    </view>

    <!-- 发帖弹窗 -->
    <view v-if="showPost" class="picker-mask" @click.self="showPost = false">
      <view class="picker-content">
        <picker mode="selector" :range="tagOptions" :value="tagIndex" @change="onTagChange">
          <view style="padding:8px 0;border-bottom:1px solid #eee;margin-bottom:8px;">
            标签：{{ tagOptions[tagIndex] }}
          </view>
        </picker>
        <textarea v-model="postContent" placeholder="说点什么..." style="width:100%;height:100px;border:1px solid #eee;padding:8px;" />
        <button style="padding:8px 0;color:#576b95;font-size:14px;" @click="chooseImages">选择图片（最多9张）</button>
        <view v-if="previewImages.length" style="display:flex;gap:4px;margin-top:8px;flex-wrap:wrap;">
          <image v-for="(img, idx) in previewImages" :key="idx" :src="img" style="width:60px;height:60px;border-radius:4px;" mode="aspectFill" />
        </view>
        <button class="btn-primary" style="margin-top:16px;" @click="submitPost" :disabled="postLoading">{{ postLoading ? '发布中...' : '发布' }}</button>
      </view>
    </view>

    <!-- Feed 列表 -->
    <view v-if="posts.length === 0 && !loading" class="empty">暂无动态</view>
    <view v-for="post in posts" :key="post.id" class="card">
      <view class="card-header">
        <image class="avatar" :src="post.author?.avatar_url || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + post.author?.id" mode="aspectFill" />
        <view class="info">
          <text class="nickname">{{ post.author?.nickname || '用户' + post.author?.id }}</text>
          <text class="time">{{ formatTime(post.created_at) }}</text>
        </view>
      </view>
      <view style="margin-bottom:6px;">
        <text class="tag" :class="'tag-' + post.tag">{{ tagText(post.tag) }}</text>
      </view>
      <text class="content">{{ post.content }}</text>
      <view v-if="post.images && post.images.length" class="images">
        <image v-for="(img, idx) in post.images.filter(Boolean)" :key="idx" :src="img" mode="aspectFill" />
      </view>
      <view class="actions">
        <text :class="{active: post.is_liked}" @click="toggleLike(post)">
          {{ post.is_liked ? '❤️' : '🤍' }} {{ post.like_count }}
        </text>
        <text @click="toggleComment(post)">
          💬 {{ post.comment_count }}
        </text>
      </view>

      <!-- 评论 -->
      <view v-if="post.showComment" class="comment-box">
        <view v-for="c in post.comments || []" :key="c.id" class="comment">
          <text style="color:#576b95;font-size:13px;">{{ c.author?.nickname || '用户' + c.user_id }}</text>
          <text style="font-size:14px;color:#333;margin-top:2px;display:block;">{{ c.content }}</text>
        </view>
        <view v-if="!post.comments || post.comments.length===0" style="font-size:12px;color:#999;padding:4px 0;">暂无评论</view>
        <view style="display:flex;gap:6px;margin-top:8px;">
          <input v-model="post.newComment" placeholder="写评论..." style="flex:1;padding:6px 10px;border:1px solid #ddd;border-radius:4px;font-size:13px;" />
          <button style="padding:6px 12px;background:#07c160;color:#fff;border-radius:4px;font-size:13px;" @click="submitComment(post)">发送</button>
        </view>
      </view>
    </view>

    <view v-if="loading" class="loading">加载中...</view>
    <view v-if="hasMore && posts.length > 0" style="text-align:center;padding:16px;color:#999;font-size:13px;" @click="loadMore">点击加载更多</view>

    <button class="fab" @click="showPost = true">+</button>
  </view>
</template>

<script>
import { api } from '../../utils/api.js'

const tagMap = { help: '求助', share: '分享', notice: '通知', qa: '问答' }
const tagKeys = ['share', 'help', 'notice', 'qa']
const tagLabels = ['分享', '求助', '通知', '问答']

export default {
  data() {
    return {
      communities: [],
      currentCommunity: null,
      posts: [],
      loading: false,
      hasMore: true,
      lastId: 0,
      showPicker: false,
      showPost: false,
      tagIndex: 0,
      tagOptions: tagLabels,
      postContent: '',
      selectedFiles: [],
      previewImages: []
    }
  },
  onShow() {
    this.loadCommunities().then(() => {
      if (this.currentCommunity && this.posts.length === 0) this.loadPosts()
    })
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
      this.showPicker = false
      this.posts = []; this.lastId = 0; this.hasMore = true
      this.loadPosts()
    },
    async loadPosts() {
      if (!this.currentCommunity || this.loading) return
      this.loading = true
      try {
        const list = await api.listPosts(this.currentCommunity.id, this.lastId)
        if (list.length < 10) this.hasMore = false
        if (list.length > 0) this.lastId = list[list.length - 1].id
        this.posts.push(...list)
      } catch (e) { console.error(e) }
      finally { this.loading = false }
    },
    loadMore() { this.loadPosts() },
    onTagChange(e) { this.tagIndex = e.detail.value },
    chooseImages() {
      uni.chooseImage({
        count: 9,
        sizeType: ['compressed'],
        success: (res) => {
          this.selectedFiles = res.tempFilePaths
          this.previewImages = res.tempFilePaths
        }
      })
    },
    async submitPost() {
      if (!this.postContent.trim()) return
      this.postLoading = true
      try {
        let images = []
        if (this.selectedFiles.length > 0) {
          images = await uploadFiles(this.selectedFiles)
        }
        await api.createPost(this.currentCommunity.id, tagKeys[this.tagIndex], this.postContent, images)
        this.showPost = false
        this.postContent = ''; this.selectedFiles = []; this.previewImages = []; this.tagIndex = 0
        this.posts = []; this.lastId = 0; this.hasMore = true
        this.loadPosts()
      } catch (e) {
        uni.showToast({ title: e.message, icon: 'none' })
      } finally {
        this.postLoading = false
      }
    },
    async toggleLike(post) {
      try {
        if (post.is_liked) {
          await api.unlikePost(post.id)
          post.is_liked = false; post.like_count--
        } else {
          await api.likePost(post.id)
          post.is_liked = true; post.like_count++
        }
      } catch (e) {
        uni.showToast({ title: e.message, icon: 'none' })
      }
    },
    async toggleComment(post) {
      post.showComment = !post.showComment
      if (post.showComment && !post.commentsLoaded) {
        try {
          post.comments = await api.listComments(post.id)
          post.commentsLoaded = true
        } catch (e) { console.error(e) }
      }
    },
    async submitComment(post) {
      if (!post.newComment?.trim()) return
      try {
        await api.createComment(post.id, post.newComment.trim())
        post.newComment = ''
        post.comments = await api.listComments(post.id)
        post.comment_count++
      } catch (e) {
        uni.showToast({ title: e.message, icon: 'none' })
      }
    },
    goCommunities() {
      this.showPicker = false
      uni.navigateTo({ url: '/src/pages/communities/index' })
    },
    formatTime(t) {
      const d = new Date(t)
      return `${d.getMonth()+1}月${d.getDate()}日 ${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`
    },
    tagText(tag) { return tagMap[tag] || tag }
  }
}
</script>

<style scoped>
.page { padding-bottom: 20px; }
.community-bar { padding: 10px 16px; background: #fff; border-bottom: 1px solid #eee; display: flex; align-items: center; font-size: 15px; color: #333; }
.picker-mask { position: fixed; inset: 0; background: rgba(0,0,0,0.5); z-index: 200; display: flex; align-items: flex-end; }
.picker-content { background: #fff; width: 100%; border-radius: 12px 12px 0 0; padding: 16px; max-height: 60vh; overflow-y: auto; }
.picker-item { padding: 12px 0; border-bottom: 1px solid #f5f5f5; display: flex; justify-content: space-between; font-size: 15px; }
.card { background: #fff; padding: 12px 16px; margin-bottom: 8px; }
.card-header { display: flex; align-items: center; margin-bottom: 8px; }
.avatar { width: 40px; height: 40px; border-radius: 50%; background: #ddd; }
.info { margin-left: 10px; }
.nickname { font-size: 15px; font-weight: 500; color: #333; display: block; }
.time { font-size: 12px; color: #999; }
.tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 12px; margin-right: 6px; }
.tag-help { background: #fff2f0; color: #fa5151; }
.tag-share { background: #e6f7ff; color: #1890ff; }
.tag-notice { background: #fff7e6; color: #fa8c16; }
.tag-qa { background: #f6ffed; color: #52c41a; }
.content { font-size: 15px; line-height: 1.6; color: #333; word-break: break-all; }
.images { display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px; margin-top: 8px; }
.images image { width: 100%; height: 100px; border-radius: 4px; }
.actions { display: flex; justify-content: flex-end; gap: 16px; margin-top: 10px; padding-top: 10px; border-top: 1px solid #f5f5f5; font-size: 13px; color: #999; }
.actions .active { color: #07c160; }
.comment-box { margin-top: 10px; background: #f9f9f9; padding: 10px; border-radius: 6px; }
.comment { padding: 6px 0; border-bottom: 1px solid #eee; }
.fab { position: fixed; bottom: 30px; right: 16px; width: 50px; height: 50px; border-radius: 50%; background: #07c160; color: #fff; font-size: 28px; display: flex; align-items: center; justify-content: center; box-shadow: 0 2px 8px rgba(0,0,0,0.15); z-index: 99; }
.empty { text-align: center; padding: 60px 20px; color: #999; font-size: 14px; }
.loading { text-align: center; padding: 20px; color: #999; }
.btn-primary { background: #07c160; color: #fff; border-radius: 6px; font-size: 16px; }
</style>
