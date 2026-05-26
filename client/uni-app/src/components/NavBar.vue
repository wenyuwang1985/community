<template>
  <view class="nav-bar">
    <view class="nav-left" @click="goBack" v-if="showBack">
      <text class="back-icon">←</text>
    </view>
    <view class="nav-title">{{ title }}</view>
    <view class="nav-right">
      <slot name="right"></slot>
    </view>
  </view>
</template>

<script>
export default {
  props: {
    title: { type: String, default: '' },
    showBack: { type: Boolean, default: true }
  },
  methods: {
    goBack() {
      if (getCurrentPages().length > 1) {
        uni.navigateBack()
      } else {
        uni.switchTab({ url: '/src/pages/index/index' })
      }
    }
  }
}
</script>

<style scoped>
.nav-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 44px;
  background: #07c160;
  color: #fff;
  padding: 0 16px;
  font-size: 16px;
  font-weight: 500;
}

.nav-left {
  width: 50px;
  display: flex;
  align-items: center;
}

.back-icon {
  font-size: 20px;
  font-weight: 600;
}

.nav-title {
  flex: 1;
  text-align: center;
}

.nav-right {
  width: 50px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}
</style>