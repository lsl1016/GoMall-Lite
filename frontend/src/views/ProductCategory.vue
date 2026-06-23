<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import ProductCard from '@/components/ProductCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import { getProductsAPI } from '@/api/product'

const route = useRoute()
const router = useRouter()
const categories = ['全部', '手机数码', '电脑办公', '家用电器', '服饰鞋包', '美妆个护', '食品生鲜']
const currentCategory = ref(route.query.category || '全部')
const sort = ref('综合排序')
const products = ref([])
const loading = ref(false)
const errorMsg = ref('')

async function fetchProducts() {
  try {
    loading.value = true
    errorMsg.value = ''
    const data = await getProductsAPI({ category: currentCategory.value === '全部' ? '' : currentCategory.value, keyword: route.query.keyword || '' })
    products.value = [...data]
    if (sort.value === '价格升序') products.value.sort((a, b) => a.price - b.price)
    if (sort.value === '价格降序') products.value.sort((a, b) => b.price - a.price)
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

function selectCategory(category) {
  currentCategory.value = category
  router.replace({ path: '/category', query: { ...route.query, category: category === '全部' ? undefined : category } })
}

watch(() => route.query, () => {
  currentCategory.value = route.query.category || '全部'
  fetchProducts()
}, { deep: true })
watch(sort, fetchProducts)
onMounted(fetchProducts)
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <div class="breadcrumb">首页 / 商品分类</div>
      <div class="category-layout">
        <aside class="side-card card">
          <button v-for="cat in categories" :key="cat" :class="{ active: currentCategory === cat }" @click="selectCategory(cat)">{{ cat }}</button>
        </aside>
        <section>
          <div class="toolbar card">
            <div class="filter-buttons">
              <button v-for="item in ['综合排序', '价格升序', '价格降序']" :key="item" :class="{ active: sort === item }" @click="sort = item">{{ item }}</button>
            </div>
            <span class="muted">共 {{ products.length }} 件商品</span>
          </div>
          <div v-if="loading" class="loading card">加载中...</div>
          <div v-else-if="errorMsg" class="error-text card" style="padding: 30px">{{ errorMsg }}</div>
          <EmptyState v-else-if="!products.length" title="没有找到商品" button-text="返回首页" @action="router.push('/products')" />
          <div v-else class="product-grid">
            <ProductCard v-for="item in products" :key="item.id" :product="item" />
          </div>
        </section>
      </div>
    </main>
  </div>
</template>
