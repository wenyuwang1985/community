// 本地存储工具函数

const TOKEN_KEY = 'token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const USER_ID_KEY = 'user_id'
const CURRENT_COMMUNITY_KEY = 'current_community'

export const storage = {
  // Token 管理
  getToken() {
    return uni.getStorageSync(TOKEN_KEY) || ''
  },
  setToken(token) {
    uni.setStorageSync(TOKEN_KEY, token)
  },
  removeToken() {
    uni.removeStorageSync(TOKEN_KEY)
  },

  // Refresh Token 管理
  getRefreshToken() {
    return uni.getStorageSync(REFRESH_TOKEN_KEY) || ''
  },
  setRefreshToken(token) {
    uni.setStorageSync(REFRESH_TOKEN_KEY, token)
  },
  removeRefreshToken() {
    uni.removeStorageSync(REFRESH_TOKEN_KEY)
  },

  // 用户 ID 管理
  getUserId() {
    return uni.getStorageSync(USER_ID_KEY) || 0
  },
  setUserId(id) {
    uni.setStorageSync(USER_ID_KEY, id)
  },
  removeUserId() {
    uni.removeStorageSync(USER_ID_KEY)
  },

  // 当前社区管理
  getCurrentCommunity() {
    return uni.getStorageSync(CURRENT_COMMUNITY_KEY) || null
  },
  setCurrentCommunity(community) {
    uni.setStorageSync(CURRENT_COMMUNITY_KEY, community)
  },
  removeCurrentCommunity() {
    uni.removeStorageSync(CURRENT_COMMUNITY_KEY)
  },

  // 清除所有登录相关数据
  clearAuth() {
    this.removeToken()
    this.removeRefreshToken()
    this.removeUserId()
    this.removeCurrentCommunity()
  },

  // 通用存储
  get(key) {
    return uni.getStorageSync(key)
  },
  set(key, value) {
    uni.setStorageSync(key, value)
  },
  remove(key) {
    uni.removeStorageSync(key)
  }
}