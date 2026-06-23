<script setup>
import { useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

function logout() {
  userStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <section class="profile-card card">
        <div class="profile-info">
          <div class="profile-avatar">{{ userStore.userInfo?.nickname?.slice(0, 1) || 'U' }}</div>
          <div>
            <h1 style="margin: 0 0 8px">{{ userStore.userInfo?.nickname || '用户' }}</h1>
            <p class="muted">用户名：{{ userStore.userInfo?.username }}</p>
          </div>
        </div>
        <button class="ghost-btn" @click="logout">退出登录</button>
      </section>

      <section class="service-grid">
        <RouterLink class="service-card card" to="/orders"><div class="service-icon">📋</div><strong>我的订单</strong></RouterLink>
        <RouterLink class="service-card card" to="/cart"><div class="service-icon">🛒</div><strong>我的购物车</strong></RouterLink>
        <RouterLink class="service-card card" to="/address"><div class="service-icon">📍</div><strong>地址管理</strong></RouterLink>
        <div class="service-card card"><div class="service-icon">⚙️</div><strong>账号设置</strong></div>
      </section>
    </main>
  </div>
</template>
