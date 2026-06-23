<script setup>
defineProps({
  title: { type: String, default: '暂无数据' },
  buttonText: { type: String, default: '' }
})
defineEmits(['action'])
</script>

<template>
  <div class="empty-state">
    <div class="empty-icon">🛒</div>
    <h3>{{ title }}</h3>
    <button v-if="buttonText" class="primary-btn small" @click="$emit('action')">{{ buttonText }}</button>
  </div>
</template>
