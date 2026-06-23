<script setup>
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { getProductDetailAPI } from '@/api/product'
import { useCartStore } from '@/stores/cart'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const cartStore = useCartStore()
const userStore = useUserStore()
const product = ref(null)
const loading = ref(false)
const errorMsg = ref('')
const count = ref(1)
const color = ref('黑色')
const tab = ref('商品介绍')

async function fetchDetail() {
  try {
    loading.value = true
    errorMsg.value = ''
    product.value = await getProductDetailAPI(route.params.id)
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

function changeCount(delta) {
  count.value = Math.max(1, Math.min(product.value?.stock || 1, count.value + delta))
}

async function addToCart(goCart = false) {
  if (!userStore.isLogin) {
    router.push('/login')
    return
  }
  await cartStore.addToCart(product.value.id, count.value)
  if (goCart) router.push('/cart')
}

onMounted(fetchDetail)
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <div class="breadcrumb">首页 / 商品详情</div>
      <div v-if="loading" class="loading card">加载中...</div>
      <div v-else-if="errorMsg" class="error-text card" style="padding: 30px">{{ errorMsg }}</div>
      <template v-else-if="product">
        <section class="card detail-top">
          <div class="detail-img"><img :src="product.image" :alt="product.name" /></div>
          <div class="detail-info">
            <h1>{{ product.name }}</h1>
            <div class="big-price">￥{{ product.price }}</div>
            <p>库存： <strong>{{ product.stock }}</strong> 件</p>
            <p>分类： <RouterLink class="danger" :to="{ path: '/category', query: { category: product.category } }">{{ product.category }}</RouterLink></p>
            <div class="option-row">
              <span>颜色：</span>
              <button v-for="item in ['黑色', '白色', '蓝色']" :key="item" class="option" :class="{ active: color === item }" @click="color = item">{{ item }}</button>
            </div>
            <div class="option-row">
              <span>数量：</span>
              <div class="stepper"><button @click="changeCount(-1)">−</button><span>{{ count }}</span><button @click="changeCount(1)">+</button></div>
            </div>
            <div class="detail-actions">
              <button class="orange-btn" @click="addToCart(false)">🛒 加入购物车</button>
              <button class="primary-btn" @click="addToCart(true)">立即购买</button>
            </div>
          </div>
        </section>

        <section class="card tabs-card">
          <div class="tabs">
            <button v-for="item in ['商品介绍', '参数规格', '用户评价（128）']" :key="item" :class="{ active: tab === item }" @click="tab = item">{{ item }}</button>
          </div>
          <div class="tab-content">
            <p v-if="tab === '商品介绍'">{{ product.description }}</p>
            <p v-else-if="tab === '参数规格'">品牌官方品质，正品保障。商品分类：{{ product.category }}，库存：{{ product.stock }} 件。</p>
            <p v-else>用户评价整体良好，适合教学项目演示使用。</p>
            <div class="feature-row">
              <div class="feature"><span class="feature-icon">📦</span><div><strong>正品保障</strong><br><span class="muted">放心选购</span></div></div>
              <div class="feature"><span class="feature-icon">⚡</span><div><strong>快速发货</strong><br><span class="muted">高效配送</span></div></div>
              <div class="feature"><span class="feature-icon">💬</span><div><strong>售后服务</strong><br><span class="muted">响应及时</span></div></div>
              <div class="feature"><span class="feature-icon">🔒</span><div><strong>安全支付</strong><br><span class="muted">模拟流程</span></div></div>
            </div>
          </div>
        </section>
      </template>
    </main>
  </div>
</template>
