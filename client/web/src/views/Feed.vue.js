import { api, uploadFiles } from '../api.js'

export default {
  template: `
    <div class="page">
      <div class="header">
        <span>广场 - {{ currentCommunity?.name || '选择社区' }}</span>
        <span style="position:absolute;right:16px;top:12px;font-size:14px;" @click="showCommunitySwitch=true">切换</span>
      </div>

      <!-- 社区选择弹层 -->
      <div v-if="showCommunitySwitch" class="modal-overlay" @click.self="showCommunitySwitch=false">
        <div class="modal-content">
          <div class="modal-header">选择社区</div>
          <div v-for="c in communities" :key="c.id" class="form-group" style="display:flex;justify-content:space-between;align-items:center;cursor:pointer;" @click="switchCommunity(c)">
            <span>{{ c.name }} {{ c.is_primary ? '(主)' : '' }}</span>
            <span v-if="c.id === currentCommunity?.id" style="color:#07c160;">✓</span>
          </div>
          <button class="btn btn-default" style="margin-top:12px;" @click="$router.push('/communities')">管理我的社区</button>
        </div>
      </div>

      <!-- 发帖弹层 -->
      <div v-if="showPostModal" class="modal-overlay" @click.self="showPostModal=false">
        <div class="modal-content">
          <div class="modal-header">发布动态</div>
          <select v-model="postTag" style="width:100%;padding:8px;margin-bottom:8px;border:1px solid #ddd;border-radius:4px;">
            <option value="share">分享</option>
            <option value="help">求助</option>
            <option value="notice">通知</option>
            <option value="qa">问答</option>
          </select>
          <textarea v-model="postContent" placeholder="说点什么..." style="width:100%;min-height:120px;padding:10px;border:1px solid #ddd;border-radius:4px;resize:none;"></textarea>
          <div style="margin-top:8px;font-size:12px;color:#999;">图片功能开发中，请输入图片URL（逗号分隔）</div>
          <input v-model="postImages" placeholder="图片URL，逗号分隔" style="width:100%;padding:8px;margin-top:8px;border:1px solid #ddd;border-radius:4px;">
          <button class="btn btn-primary" style="margin-top:16px;" @click="submitPost" :disabled="postLoading">发布</button>
        </div>
      </div>

      <!-- Feed 列表 -->
      <div v-if="posts.length === 0 && !loading" class="empty">暂无动态，点击右下角发布第一条</div>
      <div v-for="post in posts" :key="post.id" class="card">
        <div class="card-header">
          <div class="avatar"><img :src="post.author?.avatar_url || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + post.author?.id" alt=""></div>
          <div class="info">
            <div class="nickname">{{ post.author?.nickname || '用户' + post.author?.id }}</div>
            <div class="time">{{ formatTime(post.created_at) }}</div>
          </div>
        </div>
        <div style="margin-bottom:6px;">
          <span class="tag" :class="'tag-' + post.tag">{{ tagText(post.tag) }}</span>
        </div>
        <div class="content-text">{{ post.content }}</div>
        <div v-if="post.images && post.images.length" class="content-images">
          <img v-for="(img, idx) in post.images.filter(Boolean)" :key="idx" :src="img" alt="">
        </div>
        <div class="actions">
          <span class="action-item" :class="{active: post.is_liked}" @click="toggleLike(post)">
            {{ post.is_liked ? '❤️' : '🤍' }} {{ post.like_count }}
          </span>
          <span class="action-item" @click="showComments(post)">
            💬 {{ post.comment_count }}
          </span>
        </div>

        <!-- 评论展开 -->
        <div v-if="post.showComment" style="margin-top:10px;background:#f9f9f9;padding:10px;border-radius:6px;">
          <div v-for="c in post.comments || []" :key="c.id" class="comment">
            <span class="comment-author">{{ c.author?.nickname || '用户' + c.user_id }}</span>
            <div class="comment-text">{{ c.content }}</div>
          </div>
          <div v-if="!post.comments || post.comments.length===0" style="font-size:12px;color:#999;padding:4px 0;">暂无评论</div>
          <div style="display:flex;gap:6px;margin-top:8px;">
            <input v-model="post.newComment" placeholder="写评论..." style="flex:1;padding:6px 10px;border:1px solid #ddd;border-radius:4px;font-size:13px;">
            <button style="padding:6px 12px;background:#07c160;color:#fff;border-radius:4px;font-size:13px;" @click="submitComment(post)">发送</button>
          </div>
        </div>
      </div>

      <div v-if="loading" class="loading">加载中...</div>
      <div v-if="hasMore && posts.length > 0" style="text-align:center;padding:16px;color:#999;font-size:13px;" @click="loadMore">点击加载更多</div>

      <button class="fab" @click="showPostModal=true">+</button>

      <div class="tabbar">
        <div class="tabbar-item active"><span class="icon">📰</span>广场</div>
        <div class="tabbar-item" @click="$router.push('/market')"><span class="icon">🛒</span>集市</div>
        <div class="tabbar-item" @click="$router.push('/chat')"><span class="icon">💬</span>聊天</div>
        <div class="tabbar-item" @click="$router.push('/profile')"><span class="icon">👤</span>我的</div>
      </div>
    </div>
  `,
  data() {
    return {
      communities: [],
      currentCommunity: null,
      posts: [],
      loading: false,
      hasMore: true,
      lastId: 0,
      showCommunitySwitch: false,
      showPostModal: false,
      postTag: 'share',
      postContent: '',
      selectedFiles: [],
      previewImages: [],
      postLoading: false,
    }
  },
  async mounted() {
    await this.loadCommunities()
    if (this.currentCommunity) this.loadPosts()
  },
  methods: {
    async loadCommunities() {
      try {
        this.communities = await api.getMyCommunities()
        this.currentCommunity = this.communities.find(c => c.is_primary) || this.communities[0]
      } catch (e) {
        console.error(e)
      }
    },
    switchCommunity(c) {
      this.currentCommunity = c
      this.showCommunitySwitch = false
      this.posts = []
      this.lastId = 0
      this.hasMore = true
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
      } catch (e) {
        console.error(e)
      } finally {
        this.loading = false
      }
    },
    loadMore() {
      this.loadPosts()
    },
    onFileChange(e) {
      const files = Array.from(e.target.files || [])
      if (files.length > 9) { alert('最多选择9张图片'); return }
      this.selectedFiles = files
      this.previewImages = files.map(f => URL.createObjectURL(f))
    },
    async submitPost() {
      if (!this.postContent.trim()) return
      this.postLoading = true
      try {
        let images = []
        if (this.selectedFiles.length > 0) {
          images = await uploadFiles(this.selectedFiles)
        }
        await api.createPost(this.currentCommunity.id, this.postTag, this.postContent, images)
        this.showPostModal = false
        this.postContent = ''; this.selectedFiles = []; this.previewImages = []; this.postTag = 'share'
        if (this.$refs.fileInput) this.$refs.fileInput.value = ''
        this.posts = []; this.lastId = 0; this.hasMore = true
        this.loadPosts()
      } catch (e) {
        alert(e.message)
      } finally {
        this.postLoading = false
      }
    },
    async toggleLike(post) {
      try {
        if (post.is_liked) {
          await api.unlikePost(post.id)
          post.is_liked = false
          post.like_count--
        } else {
          await api.likePost(post.id)
          post.is_liked = true
          post.like_count++
        }
      } catch (e) {
        alert(e.message)
      }
    },
    async showComments(post) {
      post.showComment = !post.showComment
      if (post.showComment && !post.commentsLoaded) {
        try {
          post.comments = await api.listComments(post.id)
          post.commentsLoaded = true
        } catch (e) {
          console.error(e)
        }
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
        alert(e.message)
      }
    },
    formatTime(t) {
      const d = new Date(t)
      return `${d.getMonth()+1}月${d.getDate()}日 ${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`
    },
    tagText(tag) {
      const map = { help: '求助', share: '分享', notice: '通知', qa: '问答' }
      return map[tag] || tag
    }
  }
}
