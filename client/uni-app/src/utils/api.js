const BASE = 'http://localhost:8080/api/v1'

function getToken() {
  return uni.getStorageSync('token') || ''
}

function request(method, path, body) {
  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE + path,
      method,
      header: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + getToken()
      },
      data: body,
      success: (res) => {
        const data = res.data
        if (data.code !== 0) {
          reject(new Error(data.msg || '请求失败'))
        } else {
          resolve(data.data)
        }
      },
      fail: (err) => reject(new Error(err.errMsg || '网络错误'))
    })
  })
}

export const api = {
  login: (phone, password) => request('POST', '/auth/login', { phone, password }),
  register: (phone, password) => request('POST', '/auth/register', { phone, password }),
  refresh: (refreshToken) => request('POST', '/auth/refresh', { refresh_token: refreshToken }),

  getProfile: () => request('GET', '/user/profile'),
  updateProfile: (nickname, avatar_url) => request('PUT', '/user/profile', { nickname, avatar_url }),
  getUser: (id) => request('GET', `/users/${id}`),

  searchCommunities: (q) => request('GET', '/communities/search?q=' + encodeURIComponent(q)),
  subscribeCommunity: (id) => request('POST', `/communities/${id}/subscribe`),
  unsubscribeCommunity: (id) => request('DELETE', `/communities/${id}/subscribe`),
  getMyCommunities: () => request('GET', '/user/communities'),
  setPrimaryCommunity: (id) => request('PUT', `/user/communities/${id}/primary`),

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

  createItem: (body) => request('POST', '/items', body),
  listItems: (community_id, category, last_id, limit = 10) => {
    let url = `/items?community_id=${community_id}`
    if (category) url += `&category=${category}`
    if (last_id) url += `&last_id=${last_id}`
    url += `&limit=${limit}`
    return request('GET', url)
  },
  getItem: (id) => request('GET', `/items/${id}`),
  markSold: (id) => request('PUT', `/items/${id}/sold`),
  markOff: (id) => request('PUT', `/items/${id}/off`),

  getConversations: () => request('GET', '/conversations'),
  createPrivateConversation: (target_user_id) => request('POST', '/conversations/private', { target_user_id }),
  getMessages: (id, last_id) => {
    let url = `/conversations/${id}/messages`
    if (last_id) url += `?last_id=${last_id}`
    return request('GET', url)
  }
}
