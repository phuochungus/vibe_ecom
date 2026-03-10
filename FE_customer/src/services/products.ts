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
    const query = params
      ? {
          q: params.search,
          status: params.status,
          min_price: params.min_price,
          max_price: params.max_price,
          page: params.page,
          page_size: params.page_size,
          sort:
            params.sort === 'price_asc' || params.sort === 'price_desc'
              ? 'price'
              : params.sort === 'newest' || params.sort === 'oldest'
                ? 'created_at'
                : undefined,
          order:
            params.sort === 'price_desc' || params.sort === 'oldest'
              ? 'desc'
              : params.sort === 'price_asc' || params.sort === 'newest'
                ? 'asc'
                : undefined,
        }
      : undefined

    const { data } = await api.get<{ success: boolean; data: PaginatedResponse<Product> }>('/products', { params: query })
    return data.data
  },

  getById: async (id: string) => {
    const { data } = await api.get<{ success: boolean; data: Product }>(`/products/${id}`)
    return data.data
  },
}
