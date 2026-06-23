import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 10000
})

request.interceptors.request.use((config) => {
  const token = localStorage.getItem('gomall_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => {
    const body = response.data
    if (body && body.code === 200) return body.data
    return Promise.reject(new Error(body?.message || '请求失败'))
  },
  (error) => {
    const message = error.response?.data?.message || error.message || '网络错误'
    if (error.response?.status === 401) {
      localStorage.removeItem('gomall_token')
      localStorage.removeItem('gomall_user')
      if (location.pathname !== '/login') location.href = '/login'
    }
    return Promise.reject(new Error(message))
  }
)

export default request
