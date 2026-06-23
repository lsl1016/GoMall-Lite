<script setup>
import { useRoute, useRouter } from 'vue-router'
import NavBar from '@/components/NavBar.vue'

const route = useRoute()
const router = useRouter()
</script>

<template>
  <div class="page">
    <NavBar />
    <main>
      <section class="result-card card">
        <div class="success-icon">✓</div>
        <h1>下单成功</h1>
        <p class="muted">您的订单已提交，请尽快完成支付</p>
        <div class="result-info">
          <p>订单号：<strong>{{ route.query.orderNo || '-' }}</strong></p>
          <p>应付金额：<span class="total-price">￥{{ route.query.amount || 0 }}</span></p>
        </div>
        <div class="result-actions">
          <button class="primary-btn" @click="router.push('/orders')">查看订单</button>
          <button class="ghost-btn" @click="router.push('/products')">返回首页</button>
        </div>
      </section>
    </main>
  </div>
</template>
