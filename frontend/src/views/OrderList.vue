<script setup>
import { computed, onMounted, ref } from 'vue'
import NavBar from '@/components/NavBar.vue'
import EmptyState from '@/components/EmptyState.vue'
import { cancelOrderAPI, getOrdersAPI, payOrderAPI } from '@/api/order'

const orders = ref([])
const loading = ref(false)
const errorMsg = ref('')
const activeStatus = ref('全部')
const tabs = ['全部', '待支付', '已支付', '待收货', '已完成', '已取消']

const displayOrders = computed(() => {
  if (activeStatus.value === '全部') return orders.value
  return orders.value.filter((item) => item.status === activeStatus.value)
})

async function fetchOrders() {
  try {
    loading.value = true
    errorMsg.value = ''
    orders.value = await getOrdersAPI()
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

async function pay(id) {
  await payOrderAPI(id)
  await fetchOrders()
}

async function cancel(id) {
  await cancelOrderAPI(id)
  await fetchOrders()
}

function statusClass(status) {
  if (status === '待支付') return 'warning'
  if (status === '已支付' || status === '已完成') return 'success'
  if (status === '已取消') return 'muted'
  return ''
}

onMounted(fetchOrders)
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <section class="card order-card-wrap">
        <h1 style="font-size: 34px; margin-top: 0">订单中心</h1>
        <div class="order-tabs">
          <button v-for="tab in tabs" :key="tab" :class="{ active: activeStatus === tab }" @click="activeStatus = tab">{{ tab }}</button>
        </div>
        <div v-if="loading" class="loading">订单加载中...</div>
        <p v-else-if="errorMsg" class="error-text">{{ errorMsg }}</p>
        <EmptyState v-else-if="!displayOrders.length" title="暂无订单" />
        <div v-else>
          <div v-for="order in displayOrders" :key="order.id" class="order-item-card">
            <div class="order-head">
              <strong>订单号：{{ order.orderNo }}</strong>
              <span class="muted">创建时间：{{ order.createdAt }}</span>
              <span class="order-status" :class="statusClass(order.status)">{{ order.status }}</span>
            </div>
            <div class="order-body">
              <div class="order-products">
                <div v-for="item in order.items" :key="item.id" class="order-product">
                  <img :src="item.productImage" :alt="item.productName" />
                  <div>
                    <strong>{{ item.productName }}</strong>
                    <p class="muted">x{{ item.count }}</p>
                  </div>
                </div>
              </div>
              <div class="order-actions">
                <div>合计：<span class="price">￥{{ order.totalAmount }}</span></div>
                <button v-if="order.status === '待支付'" class="primary-btn small" @click="pay(order.id)">去支付</button>
                <button v-if="order.status === '待支付'" class="ghost-btn small" @click="cancel(order.id)">取消订单</button>
                <button v-else class="ghost-btn small">查看详情</button>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>
