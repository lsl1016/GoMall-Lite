<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import EmptyState from '@/components/EmptyState.vue'
import { useCartStore } from '@/stores/cart'
import { getAddressesAPI } from '@/api/address'
import { createOrderAPI } from '@/api/order'

const router = useRouter()
const cartStore = useCartStore()
const addresses = ref([])
const selectedAddressId = ref(null)
const remark = ref('')
const loading = ref(false)
const errorMsg = ref('')
const selectedItems = computed(() => cartStore.selectedList)
const selectedAddress = computed(() => addresses.value.find((item) => item.id === selectedAddressId.value))

async function init() {
  try {
    loading.value = true
    errorMsg.value = ''
    await cartStore.fetchCart()
    addresses.value = await getAddressesAPI()
    const def = addresses.value.find((item) => item.isDefault) || addresses.value[0]
    selectedAddressId.value = def?.id || null
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

async function submit() {
  if (!selectedAddressId.value) {
    errorMsg.value = '请选择收货地址'
    return
  }
  if (!selectedItems.value.length) {
    router.push('/cart')
    return
  }
  try {
    loading.value = true
    const order = await createOrderAPI({
      addressId: selectedAddressId.value,
      remark: remark.value,
      items: selectedItems.value.map((item) => ({ cartId: item.id, productId: item.productId, count: item.count }))
    })
    await cartStore.fetchCart()
    router.push({ path: '/order-success', query: { orderNo: order.orderNo, amount: order.totalAmount } })
  } catch (error) {
    errorMsg.value = error.message
  } finally {
    loading.value = false
  }
}

onMounted(init)
</script>

<template>
  <div class="page">
    <NavBar />
    <main class="container main">
      <div class="breadcrumb">首页 / 订单确认</div>
      <div v-if="loading" class="loading card">处理中...</div>
      <p v-else-if="errorMsg" class="error-text card" style="padding: 18px 24px">{{ errorMsg }}</p>
      <EmptyState v-if="!loading && !selectedItems.length" title="没有选择结算商品" button-text="返回购物车" @action="router.push('/cart')" />
      <div v-else class="checkout-grid">
        <div class="checkout-left">
          <section class="checkout-card card">
            <h2>收货地址</h2>
            <div v-if="selectedAddress" class="address-preview">
              <div class="address-main">
                <div class="address-pin">📍</div>
                <div>
                  <strong style="font-size: 20px">{{ selectedAddress.receiver }}</strong>
                  <span style="margin-left: 18px">{{ selectedAddress.phone }}</span>
                  <span v-if="selectedAddress.isDefault" class="tag" style="margin-left: 12px">默认地址</span>
                  <p class="muted">{{ selectedAddress.province }} {{ selectedAddress.city }} {{ selectedAddress.district }} {{ selectedAddress.detail }}</p>
                </div>
              </div>
              <select v-model="selectedAddressId" class="select" style="width: 180px">
                <option v-for="addr in addresses" :key="addr.id" :value="addr.id">{{ addr.receiver }} - {{ addr.city }}</option>
              </select>
            </div>
            <EmptyState v-else title="暂无收货地址" button-text="去新增地址" @action="router.push('/address')" />
          </section>

          <section class="checkout-card card">
            <h2>商品清单</h2>
            <table class="simple-table">
              <thead><tr><th>商品</th><th>单价</th><th>数量</th><th>小计</th></tr></thead>
              <tbody>
                <tr v-for="item in selectedItems" :key="item.id">
                  <td><div class="cart-product"><img :src="item.image" /><strong>{{ item.name }}</strong></div></td>
                  <td>￥{{ item.price }}</td>
                  <td>{{ item.count }}</td>
                  <td><strong>￥{{ item.price * item.count }}</strong></td>
                </tr>
              </tbody>
            </table>
          </section>

          <section class="checkout-card card">
            <h2>订单备注</h2>
            <textarea v-model="remark" class="textarea" maxlength="200" placeholder="选填，请先和商家协商一致"></textarea>
            <p class="muted" style="text-align: right">{{ remark.length }}/200</p>
          </section>
        </div>

        <aside class="summary-card card">
          <h2 style="margin-top: 0">金额汇总</h2>
          <div class="summary-line"><span>商品总数</span><strong>{{ cartStore.selectedCount }} 件</strong></div>
          <div class="summary-line"><span>商品金额</span><strong>￥{{ cartStore.totalPrice }}</strong></div>
          <div class="summary-line"><span>运费</span><strong>￥0</strong></div>
          <div class="summary-total"><span>应付金额</span><span class="total-price">￥{{ cartStore.totalPrice }}</span></div>
          <button class="primary-btn" style="width: 100%; margin-top: 22px" :disabled="loading" @click="submit">提交订单</button>
        </aside>
      </div>
    </main>
  </div>
</template>
