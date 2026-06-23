import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { addCartAPI, clearCartAPI, getCartAPI, removeCartAPI, updateCartAPI } from '@/api/cart'
import { useUserStore } from './user'

export const useCartStore = defineStore('cart', () => {
  const cartList = ref([])
  const loading = ref(false)

  const selectedList = computed(() => cartList.value.filter((item) => item.checked))
  const selectedCount = computed(() => selectedList.value.reduce((sum, item) => sum + item.count, 0))
  const selectedKinds = computed(() => selectedList.value.length)
  const totalPrice = computed(() => selectedList.value.reduce((sum, item) => sum + item.price * item.count, 0))
  const allChecked = computed(() => cartList.value.length > 0 && cartList.value.every((item) => item.checked))
  const badgeCount = computed(() => cartList.value.reduce((sum, item) => sum + item.count, 0))

  async function fetchCart() {
    const userStore = useUserStore()
    if (!userStore.isLogin) {
      cartList.value = []
      return []
    }
    loading.value = true
    try {
      cartList.value = await getCartAPI()
      return cartList.value
    } finally {
      loading.value = false
    }
  }

  async function addToCart(productId, count = 1) {
    cartList.value = await addCartAPI({ productId, count })
  }

  async function updateCartItem(id, data) {
    cartList.value = await updateCartAPI(id, data)
  }

  async function removeCartItem(id) {
    cartList.value = await removeCartAPI(id)
  }

  async function clearCart() {
    await clearCartAPI()
    cartList.value = []
  }

  async function toggleAll(checked) {
    await Promise.all(cartList.value.map((item) => updateCartAPI(item.id, { checked })))
    await fetchCart()
  }

  return {
    cartList,
    loading,
    selectedList,
    selectedCount,
    selectedKinds,
    totalPrice,
    allChecked,
    badgeCount,
    fetchCart,
    addToCart,
    updateCartItem,
    removeCartItem,
    clearCart,
    toggleAll
  }
})
