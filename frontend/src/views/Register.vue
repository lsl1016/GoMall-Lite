<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import NavBar from '@/components/NavBar.vue'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const errorMsg = ref('')
const form = reactive({ username: '', password: '', nickname: '' })

async function submit() {
  try {
    loading.value = true
    errorMsg.value = ''
    await userStore.register(form)
    router.push('/login')
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <NavBar />
    <main class="auth-main card">
      <section class="auth-visual"><img src="/images/cart-hero.svg" alt="购物车插画" /></section>
      <section class="auth-panel">
        <h1>创建账号</h1>
        <p class="subtitle">开始你的购物体验</p>
        <div class="form-row"><label>用户名</label><input v-model="form.username" class="input" placeholder="请输入用户名" /></div>
        <div class="form-row"><label>昵称</label><input v-model="form.nickname" class="input" placeholder="请输入昵称" /></div>
        <div class="form-row"><label>密码</label><input v-model="form.password" class="input" type="password" placeholder="请输入密码" @keyup.enter="submit" /></div>
        <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
        <button class="primary-btn" style="width: 100%" :disabled="loading" @click="submit">{{ loading ? '注册中...' : '注册' }}</button>
        <div class="auth-bottom">已有账号？ <RouterLink to="/login">去登录</RouterLink></div>
      </section>
    </main>
  </div>
</template>
