<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useCartStore } from '@/stores/cart'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const cartStore = useCartStore()
const keyword = ref(route.query.keyword || '')

onMounted(() => {
  if (userStore.isLogin) cartStore.fetchCart().catch(() => {})
})

function search() {
  router.push({ path: '/category', query: { keyword: keyword.value } })
}

function logout() {
  userStore.logout()
  cartStore.cartList = []
  router.push('/login')
}
</script>

<template>
  <header class="nav-wrap">
    <div class="nav-inner">
      <RouterLink to="/products" class="brand">
        <span class="brand-icon">🛒</span>
        <span>ShopCart</span>
      </RouterLink>
      <nav class="nav-menu">
        <RouterLink to="/products">首页</RouterLink>
        <RouterLink to="/category">商品分类</RouterLink>
        <RouterLink to="/cart">购物车</RouterLink>
        <RouterLink to="/orders">订单中心</RouterLink>
        <RouterLink to="/profile">个人中心</RouterLink>
      </nav>
      <div class="nav-actions">
        <div class="search-box">
          <input v-model="keyword" placeholder="搜索商品..." @keyup.enter="search" />
          <button @click="search">⌕</button>
        </div>
        <RouterLink to="/cart" class="cart-badge">
          🛒
          <span v-if="cartStore.badgeCount">{{ cartStore.badgeCount }}</span>
        </RouterLink>
        <template v-if="userStore.isLogin">
          <div class="avatar">{{ userStore.userInfo?.nickname?.slice(0, 1) || 'U' }}</div>
          <button class="text-btn" @click="logout">退出</button>
        </template>
        <RouterLink v-else to="/login" class="login-link">登录</RouterLink>
      </div>
    </div>
  </header>
</template>
