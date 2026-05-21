const BASE = 'http://localhost:8080/api/v1'

function getToken() {
  return localStorage.getItem('token') || ''
}

async function request(method, path, body) {
  const opts = {
    method,
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + getToken()
    }
  }
  if (body) opts.body = JSON.stringify(body)
  const res = await fetch(BASE + path, opts)
  const data = await res.json().catch(() => ({}))
  if (data.code !== 0) {
    throw new Error(data.msg || '请求失败')
  }
  return data.data
}

export const api = {
  // 认证
  login: (phone, password) => request('POST', '/auth/login', { phone, password }),
  register: (phone, password) => request('POST', '/auth/register', { phone, password }),
  refresh: (refreshToken) => request('POST', '/auth/refresh', { refresh_token: refreshToken }),

  // 用户
  getProfile: () => request('GET', '/user/profile'),
  updateProfile: (nickname, avatar_url) => request('PUT', '/user/profile', { nickname, avatar_url }),
  getUser: (id) => request('GET', `/users/${id}`),

  // 社区
  searchCommunities: (q) => request('GET', '/communities/search?q=' + encodeURIComponent(q)),
  subscribeCommunity: (id) => request('POST', `/communities/${id}/subscribe`),
  unsubscribeCommunity: (id) => request('DELETE', `/communities/${id}/subscribe`),
  getMyCommunities: () => request('GET', '/user/communities'),
  setPrimaryCommunity: (id) => request('PUT', `/user/communities/${id}/primary`),

  // 动态
  createPost: (community_id, tag, content, images) => request('POST', '/posts', { community_id, tag, content, images }),
  listPosts: (community_id, last_id, limit = 10) => {
    let url = `/posts?community_id=${community_id}`
    if (last_id) url += `&last_id=${last_id}`
    url += `&limit=${limit}`
    return request('GET', url)
  },
  getPost: (id) => request('GET', `/posts/${id}`),
  deletePost: (id) => request('DELETE', `/posts/${id}`),
  likePost: (id) => request('POST', `/posts/${id}/like`),
  unlikePost: (id) => request('DELETE', `/posts/${id}/like`),
  createComment: (id, content) => request('POST', `/posts/${id}/comments`, { content }),
  listComments: (id) => request('GET', `/posts/${id}/comments`),

  // 集市
  createItem: (body) => request('POST', '/items', body),
  listItems: (community_id, category, last_id, limit = 10) => {
    let url = `/items?community_id=${community_id}`
    if (category) url += `&category=${category}`
    if (last_id) url += `&last_id=${last_id}`
    url += `&limit=${limit}`
    return request('GET', url)
  },
  getItem: (id) => request('GET', `/items/${id}`),
  updateItem: (id, body) => request('PUT', `/items/${id}`, body),
  markSold: (id) => request('PUT', `/items/${id}/sold`),
  markOff: (id) => request('PUT', `/items/${id}/off`),

  // 聊天
  getConversations: () => request('GET', '/conversations'),
  createPrivateConversation: (target_user_id) => request('POST', '/conversations/private', { target_user_id }),
  getMessages: (id, last_id) => {
    let url = `/conversations/${id}/messages`
    if (last_id) url += `?last_id=${last_id}`
    return request('GET', url)
  }
}
