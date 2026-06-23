import request from '@/utils/request'

export const loginAPI = (data) => request.post('/login', data)
export const registerAPI = (data) => request.post('/register', data)
export const getUserInfoAPI = () => request.get('/user/info')
