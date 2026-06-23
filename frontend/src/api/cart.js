import request from '@/utils/request'

export const getCartAPI = () => request.get('/cart')
export const addCartAPI = (data) => request.post('/cart', data)
export const updateCartAPI = (id, data) => request.put(`/cart/${id}`, data)
export const removeCartAPI = (id) => request.delete(`/cart/${id}`)
export const clearCartAPI = () => request.delete('/cart')
