import { api } from '@/lib/api'
import type { Product, PaginatedResponse } from '@/types'

export interface ProductListParams {
  page?: number
  page_size?: number
  search?: string
  status?: string
  min_price?: number
  max_price?: number
  sort?: 'price_asc' | 'price_desc' | 'newest' | 'oldest'
}

export const productsApi = {
  list: async (params?: ProductListParams) => {
    const { data } = await api.get<{ success: boolean; data: PaginatedResponse<Product> }>('/products', { params })
    return data.data
  },

  getById: async (id: string) => {
    const { data } = await api.get<{ success: boolean; data: Product }>(`/products/${id}`)
    return data.data
  },
}
