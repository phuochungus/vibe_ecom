import { api } from '@/lib/api'

export interface PaymentInitiateResponse {
  checkout_url: string
  payment_id: string
}

export const paymentsApi = {
  initiate: async (orderId: string) => {
    const { data } = await api.post<{ success: boolean; data: PaymentInitiateResponse }>(
      `/orders/${orderId}/payments`
    )
    return data.data
  },
}
