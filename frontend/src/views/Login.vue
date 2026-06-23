<script setup>
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import NavBar from '@/components/NavBar.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const loading = ref(false)
const errorMsg = ref('')
const form = reactive({ username: 'admin', password: '123456' })

async function submit() {
  try {
    loading.value = true
    errorMsg.value = ''
    await userStore.login(form)
    router.push(route.query.redirect || '/products')
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
      <section class="auth-visual">
        <img src="/images/cart-hero.svg" alt="购物车插画" />
      </section>
      <section class="auth-panel">
        <h1>欢迎回来</h1>
        <p class="subtitle">登录你的账号</p>
        <div class="form-row">
          <label>用户名</label>
          <input v-model="form.username" class="input" placeholder="请输入用户名" @keyup.enter="submit" />
        </div>
        <div class="form-row">
          <label>密码</label>
          <input v-model="form.password" class="input" type="password" placeholder="请输入密码" @keyup.enter="submit" />
        </div>
        <div class="auth-extra">
          <label><input type="checkbox" /> 记住我</label>
          <a class="danger" href="javascript:;">忘记密码?</a>
        </div>
        <p v-if="errorMsg" class="error-text">{{ errorMsg }}</p>
        <button class="primary-btn" style="width: 100%" :disabled="loading" @click="submit">
          {{ loading ? '登录中...' : '登录' }}
        </button>
        <div class="auth-bottom">还没有账号？ <RouterLink to="/register">立即注册</RouterLink></div>
      </section>
    </main>
  </div>
</template>

