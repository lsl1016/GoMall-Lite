import request from '@/utils/request'

export const createOrderAPI = (data) => request.post('/orders', data)
export const getOrdersAPI = () => request.get('/orders')
export const getOrderDetailAPI = (id) => request.get(`/orders/${id}`)
export const payOrderAPI = (id) => request.put(`/orders/${id}/pay`)
export const cancelOrderAPI = (id) => request.put(`/orders/${id}/cancel`)
