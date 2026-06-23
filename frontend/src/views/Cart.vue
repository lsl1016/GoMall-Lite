<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import EmptyState from '@/components/EmptyState.vue'
import { useCartStore } from '@/stores/cart'

const router = useRouter()
const cartStore = useCartStore()
const errorMsg = ref('')

async function init() {
  try {
    errorMsg.value = ''
    await cartStore.fetchCart()
  } catch (error) {
    errorMsg.value = error.message
  }
}

async function changeCount(item, delta) {
  const next = Math.max(1, item.count + delta)
  await cartStore.updateCartItem(item.id, { count: next })
}

async function toggleItem(item) {
  await cartStore.updateCartItem(item.id, { checked: !item.checked })
}

function checkout() {
  if (!cartStore.selectedList.length) return
  router.push('/checkout')
}

onMounted(init)
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <div class="breadcrumb">首页 / 购物车</div>
      <section class="card table-card">
        <h1 style="font-size: 34px; margin-top: 0">购物车</h1>
        <div v-if="cartStore.loading" class="loading">购物车加载中...</div>
        <p v-else-if="errorMsg" class="error-text">{{ errorMsg }}</p>
        <EmptyState v-else-if="!cartStore.cartList.length" title="购物车还是空的" button-text="去逛逛" @action="router.push('/products')" />
        <template v-else>
          <div class="table-scroll">
            <table class="cart-table">
              <thead>
                <tr>
                  <th style="width: 60px"><input class="checkbox" type="checkbox" :checked="cartStore.allChecked" @change="cartStore.toggleAll($event.target.checked)" /></th>
                  <th>商品</th>
                  <th>单价</th>
                  <th>数量</th>
                  <th>小计</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in cartStore.cartList" :key="item.id">
                  <td><input class="checkbox" type="checkbox" :checked="item.checked" @change="toggleItem(item)" /></td>
                  <td>
                    <div class="cart-product">
                      <img :src="item.image" :alt="item.name" />
                      <div><strong>{{ item.name }}</strong><br><span class="muted">默认规格</span></div>
                    </div>
                  </td>
                  <td>￥{{ item.price }}</td>
                  <td>
                    <div class="stepper"><button @click="changeCount(item, -1)">−</button><span>{{ item.count }}</span><button @click="changeCount(item, 1)">+</button></div>
                  </td>
                  <td><strong class="price">￥{{ item.price * item.count }}</strong></td>
                  <td><button class="text-danger" @click="cartStore.removeCartItem(item.id)">删除</button></td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="cart-footer">
            <div class="cart-footer-left">
              <label><input class="checkbox" type="checkbox" :checked="cartStore.allChecked" @change="cartStore.toggleAll($event.target.checked)" /> 全选</label>
              <button class="text-danger" @click="cartStore.clearCart">清空购物车</button>
              <span>已选 <strong class="danger">{{ cartStore.selectedKinds }}</strong> 件商品</span>
            </div>
            <div class="cart-footer-right">
              <span>合计：</span>
              <span class="total-price">￥{{ cartStore.totalPrice }}</span>
              <button class="primary-btn" :disabled="!cartStore.selectedList.length" @click="checkout">去结算</button>
            </div>
          </div>
        </template>
      </section>
    </main>
  </div>
</template>
