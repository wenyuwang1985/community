// Toast 提示工具函数

export function showToast(title, icon = 'none', duration = 2000) {
  uni.showToast({
    title,
    icon,
    duration
  })
}

export function showSuccess(title) {
  showToast(title, 'success')
}

export function showError(title) {
  showToast(title, 'none')
}

export function showLoading(title = '加载中...') {
  uni.showLoading({
    title,
    mask: true
  })
}

export function hideLoading() {
  uni.hideLoading()
}

export function showModal(title, content, confirmText = '确定', cancelText = '取消') {
  return new Promise((resolve) => {
    uni.showModal({
      title,
      content,
      confirmText,
      cancelText,
      success: (res) => {
        resolve(res.confirm)
      },
      fail: () => {
        resolve(false)
      }
    })
  })
}

export function showActionSheet(itemList) {
  return new Promise((resolve) => {
    uni.showActionSheet({
      itemList,
      success: (res) => {
        resolve(res.tapIndex)
      },
      fail: () => {
        resolve(-1)
      }
    })
  })
}