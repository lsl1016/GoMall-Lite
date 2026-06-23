import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { loginAPI, registerAPI, getUserInfoAPI } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('gomall_token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('gomall_user') || 'null'))
  const isLogin = computed(() => Boolean(token.value))

  async function login(form) {
    const data = await loginAPI(form)
    token.value = data.token
    userInfo.value = data.user
    localStorage.setItem('gomall_token', data.token)
    localStorage.setItem('gomall_user', JSON.stringify(data.user))
    return data
  }

  async function register(form) {
    return registerAPI(form)
  }

  async function fetchUserInfo() {
    if (!token.value) return null
    const data = await getUserInfoAPI()
    userInfo.value = data
    localStorage.setItem('gomall_user', JSON.stringify(data))
    return data
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('gomall_token')
    localStorage.removeItem('gomall_user')
  }

  return { token, userInfo, isLogin, login, register, logout, fetchUserInfo }
})
