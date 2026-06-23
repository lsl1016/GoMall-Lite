import request from '@/utils/request'

export const getAddressesAPI = () => request.get('/addresses')
export const addAddressAPI = (data) => request.post('/addresses', data)
export const updateAddressAPI = (id, data) => request.put(`/addresses/${id}`, data)
export const deleteAddressAPI = (id) => request.delete(`/addresses/${id}`)
export const setDefaultAddressAPI = (id) => request.put(`/addresses/${id}/default`)
