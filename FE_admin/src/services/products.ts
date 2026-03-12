import { api } from '@/lib/api'
import type {
  PaginatedResponse,
  Product,
  ProductCreatePayload,
  ProductImageUploadResponse,
  ProductUpdatePayload,
} from '@/types'

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
    api.get<{ data: PaginatedResponse<Product> }>('/admin/products', { params }).then((r) => r.data.data),

  get: (id: string) =>
    api.get<{ data: Product }>(`/admin/products/${id}`).then((r) => r.data.data),

  create: (payload: ProductCreatePayload) =>
    api.post<{ data: Product }>('/admin/products', payload).then((r) => r.data.data),

  update: (id: string, payload: ProductUpdatePayload) =>
    api.patch<{ data: Product }>(`/admin/products/${id}`, payload).then((r) => r.data.data),

  uploadImage: (file: File) => {
    const formData = new FormData()
    formData.append('file', file)

    return api
      .post<{ data: ProductImageUploadResponse }>('/admin/products/upload-image', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      })
      .then((r) => r.data.data)
  },

  delete: (id: string) =>
    api.delete(`/admin/products/${id}`),
}
