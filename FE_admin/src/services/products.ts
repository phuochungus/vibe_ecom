import { api } from '@/lib/api'
import type { Product, ProductCreatePayload, ProductUpdatePayload, PaginatedResponse } from '@/types'

export const productsApi = {
  list: (params?: {
    q?: string
    status?: string
    min_price?: string
    max_price?: string
    page?: number
    page_size?: number
    sort?: string
    order?: string
  }) =>
    api.get<{ data: PaginatedResponse<Product> }>('/products', { params }).then((r) => r.data.data),

  get: (id: string) =>
    api.get<{ data: Product }>(`/products/${id}`).then((r) => r.data.data),

  create: (payload: ProductCreatePayload) =>
    api.post<{ data: Product }>('/admin/products', payload).then((r) => r.data.data),

  update: (id: string, payload: ProductUpdatePayload) =>
    api.patch<{ data: Product }>(`/admin/products/${id}`, payload).then((r) => r.data.data),

  delete: (id: string) =>
    api.delete(`/admin/products/${id}`),
}
