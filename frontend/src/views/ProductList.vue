<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import ProductCard from '@/components/ProductCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import { getProductsAPI } from '@/api/product'

const router = useRouter()
const products = ref([])
const loading = ref(false)
const errorMsg = ref('')
const categories = [
  { name: '手机数码', icon: '📱' },
  { name: '电脑办公', icon: '💻' },
  { name: '家用电器', icon: '🔌' },
  { name: '服饰鞋包', icon: '👕' },
  { name: '美妆个护', icon: '🧴' },
  { name: '食品生鲜', icon: '🍊' }
]

async function fetchProducts() {
  try {
    loading.value = true
    errorMsg.value = ''
    products.value = await getProductsAPI()
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

onMounted(fetchProducts)
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <section class="hero">
        <div>
          <h1>精选好物 轻松购物</h1>
          <p>品质生活，从这里开始</p>
          <button class="primary-btn" @click="router.push('/category')">立即购物</button>
        </div>
        <img src="/images/cart-hero.svg" alt="购物车" />
      </section>

      <section class="category-row">
        <button v-for="cat in categories" :key="cat.name" class="category-item" @click="router.push({ path: '/category', query: { category: cat.name } })">
          <div class="category-icon">{{ cat.icon }}</div>
          <div>{{ cat.name }}</div>
        </button>
      </section>

      <section class="section-title">
        <h2>推荐商品</h2>
        <RouterLink class="muted" to="/category">查看全部 ></RouterLink>
      </section>

      <div v-if="loading" class="loading card">商品加载中...</div>
      <div v-else-if="errorMsg" class="error-text card" style="padding: 30px">{{ errorMsg }}</div>
      <EmptyState v-else-if="!products.length" title="暂无商品" />
      <div v-else class="product-grid">
        <ProductCard v-for="item in products.slice(0, 8)" :key="item.id" :product="item" />
      </div>
    </main>
  </div>
</template>
