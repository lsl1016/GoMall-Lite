import request from '@/utils/request'

export const getProductsAPI = (params = {}) => request.get('/products', { params })
export const getProductDetailAPI = (id) => request.get(`/products/${id}`)
