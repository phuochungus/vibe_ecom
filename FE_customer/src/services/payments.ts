import { api } from '@/lib/api'

export interface PaymentInitiateResponse {
  checkout_url: string
  payment_id: string
}

export const paymentsApi = {
  initiate: async (
    orderId: string,
    payload: {
      provider: 'PAYOS'
      return_url: string
      cancel_url: string
    },
    idempotencyKey: string
  ) => {
    const { data } = await api.post<{ success: boolean; data: PaymentInitiateResponse }>(
      `/orders/${orderId}/payments`,
      payload,
      {
        headers: {
          'Idempotency-Key': idempotencyKey,
        },
      }
    )
    return data.data
  },
}
