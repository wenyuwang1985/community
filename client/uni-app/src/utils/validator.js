// 表单验证工具函数

export function validatePhone(phone) {
  const regex = /^1[3-9]\d{9}$/
  return regex.test(phone)
}

export function validatePassword(password) {
  return password && password.length >= 6
}

export function validateNickname(nickname) {
  return nickname && nickname.length >= 1 && nickname.length <= 50
}

export function validateUrl(url) {
  if (!url) return true
  const regex = /^https?:\/\/.+/
  return regex.test(url)
}

export function validatePrice(price) {
  return !isNaN(price) && price >= 0 && price <= 10000000
}

export function validateTitle(title) {
  return title && title.length >= 1 && title.length <= 200
}

export function validateContent(content) {
  return content && content.length >= 1 && content.length <= 5000
}

// 综合验证
export function validateRegisterForm(data) {
  const errors = []

  if (!validatePhone(data.phone)) {
    errors.push('请输入正确的手机号')
  }

  if (!validatePassword(data.password)) {
    errors.push('密码至少6位')
  }

  return errors
}

export function validateLoginForm(data) {
  const errors = []

  if (!data.phone) {
    errors.push('请输入手机号')
  }

  if (!data.password) {
    errors.push('请输入密码')
  }

  return errors
}

export function validatePostForm(data) {
  const errors = []

  if (!validateContent(data.content)) {
    errors.push('内容长度应在1-5000字之间')
  }

  return errors
}

export function validateItemForm(data) {
  const errors = []

  if (!validateTitle(data.title)) {
    errors.push('标题长度应在1-200字之间')
  }

  if (!validatePrice(data.price)) {
    errors.push('请输入有效的价格')
  }

  return errors
}