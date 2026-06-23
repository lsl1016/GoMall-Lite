<script setup>
import { useRouter } from 'vue-router'
import { useCartStore } from '@/stores/cart'
import { useUserStore } from '@/stores/user'

const props = defineProps({
  product: { type: Object, required: true }
})

const router = useRouter()
const cartStore = useCartStore()
const userStore = useUserStore()

function detail() {
  router.push(`/products/${props.product.id}`)
}

async function add() {
  if (!userStore.isLogin) {
    router.push('/login')
    return
  }
  await cartStore.addToCart(props.product.id, 1)
}
</script>

<template>
  <div class="product-card" @click="detail">
    <div class="product-img-wrap">
      <img :src="product.image" :alt="product.name" />
    </div>
    <div class="product-card-body">
      <h3>{{ product.name }}</h3>
      <p class="muted">库存 {{ product.stock }} 件</p>
      <div class="product-bottom">
        <strong class="price">￥{{ product.price }}</strong>
        <button class="icon-btn" title="加入购物车" @click.stop="add">🛒</button>
      </div>
    </div>
  </div>
</template>
