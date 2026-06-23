import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  { path: '/', redirect: '/products' },
  { path: '/login', component: () => import('@/views/Login.vue') },
  { path: '/register', component: () => import('@/views/Register.vue') },
  { path: '/products', component: () => import('@/views/ProductList.vue') },
  { path: '/category', component: () => import('@/views/ProductCategory.vue') },
  { path: '/products/:id', component: () => import('@/views/ProductDetail.vue') },
  { path: '/cart', component: () => import('@/views/Cart.vue'), meta: { requiresAuth: true } },
  { path: '/checkout', component: () => import('@/views/Checkout.vue'), meta: { requiresAuth: true } },
  { path: '/order-success', component: () => import('@/views/OrderSuccess.vue'), meta: { requiresAuth: true } },
  { path: '/orders', component: () => import('@/views/OrderList.vue'), meta: { requiresAuth: true } },
  { path: '/address', component: () => import('@/views/Address.vue'), meta: { requiresAuth: true } },
  { path: '/profile', component: () => import('@/views/UserCenter.vue'), meta: { requiresAuth: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 })
})

router.beforeEach((to) => {
  const userStore = useUserStore()
  if (to.meta.requiresAuth && !userStore.isLogin) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
})

export default router
