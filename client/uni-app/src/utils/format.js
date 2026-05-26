// 格式化工具函数

export function formatTime(timestamp) {
  const d = new Date(timestamp)
  const now = new Date()
  const diff = now - d

  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
  if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
  if (diff < 604800000) return Math.floor(diff / 86400000) + '天前'

  return `${d.getMonth() + 1}月${d.getDate()}日 ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}

export function formatPrice(price) {
  return '¥' + (price / 100).toFixed(2)
}

export function formatDistance(timestamp) {
  const d = new Date(timestamp)
  const now = new Date()
  const diff = now - d
  const days = Math.floor(diff / 86400000)

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return days + '天前'
  if (days < 30) return Math.floor(days / 7) + '周前'
  if (days < 365) return Math.floor(days / 30) + '月前'
  return Math.floor(days / 365) + '年前'
}

export function truncate(text, length = 100) {
  if (!text || text.length <= length) return text
  return text.substring(0, length) + '...'
}

export function debounce(func, wait) {
  let timeout
  return function(...args) {
    clearTimeout(timeout)
    timeout = setTimeout(() => func.apply(this, args), wait)
  }
}

export function throttle(func, wait) {
  let timeout
  return function(...args) {
    if (!timeout) {
      timeout = setTimeout(() => {
        timeout = null
        func.apply(this, args)
      }, wait)
    }
  }
}