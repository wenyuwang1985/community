<template>
  <view class="card" @click="$emit('click')">
    <slot></slot>
  </view>
</template>

<script>
export default {
  name: 'Card'
}
</script>

<style scoped>
.card {
  background: #fff;
  padding: 12px 16px;
  margin-bottom: 8px;
  border-radius: 8px;
}
</style>