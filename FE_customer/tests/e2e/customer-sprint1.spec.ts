import { expect, test } from '@playwright/test'

const user = {
  id: 'user-1',
  full_name: 'Nguyen Van A',
  email: 'user@golf.local',
  phone: '0901234567',
  role: 'USER',
  status: 'ACTIVE',
}

const products = {
  items: [
    {
      id: 'prod-1',
      sku: 'CLUB-001',
      name: 'Driver X',
      description: 'Driver cho nguoi choi golf',
      price: '1590000',
      stock: 5,
      status: 'ACTIVE',
      image_url: 'http://127.0.0.1:8080/assets/seed-images/driver-pro-x.svg',
      created_at: '2026-03-10T00:00:00Z',
      updated_at: '2026-03-10T00:00:00Z',
    },
  ],
  pagination: { page: 1, page_size: 12, total: 1, total_pages: 1 },
}

const orderSummary = {
  id: 'ord-1',
  order_code: 'ORD-20260310-0001',
  user_id: 'user-1',
  order_status: 'PENDING_PAYMENT',
  payment_status: 'UNPAID',
  currency_code: 'VND',
  subtotal_amount: '1590000',
  discount_amount: '0',
  shipping_fee: '30000',
  total_amount: '1620000',
  payment_due_at: '2026-03-10T10:30:00Z',
  shipping_recipient_name: 'Nguyen Van A',
  shipping_phone: '0901234567',
  shipping_line1: '123 Nguyen Hue',
  shipping_line2: '',
  shipping_ward: 'Ben Nghe',
  shipping_district: 'District 1',
  shipping_city: 'Ho Chi Minh',
  shipping_province: 'Ho Chi Minh',
  shipping_country_code: 'VN',
  placed_at: '2026-03-10T10:00:00Z',
  created_at: '2026-03-10T10:00:00Z',
}

const orderDetail = {
  ...orderSummary,
  customer_note: 'Giao gio hanh chinh',
  items: [
    {
      id: 'item-1',
      product_id: 'prod-1',
      sku: 'CLUB-001',
      name: 'Driver X',
      unit_price: '1590000',
      quantity: 1,
      line_total: '1590000',
    },
  ],
}

const tracking = {
  order_id: 'ord-1',
  current_status: 'PENDING_PAYMENT',
  timeline: [
    {
      status: 'PENDING_PAYMENT',
      source_type: 'SYSTEM',
      occurred_at: '2026-03-10T10:00:00Z',
      description: 'Order created',
    },
    {
      status: 'PAID',
      source_type: 'PAYMENT_GATEWAY',
      occurred_at: '2026-03-10T10:05:00Z',
      description: 'Payment confirmed',
    },
  ],
}

const notifications = {
  items: [
    {
      id: 'noti-1',
      user_id: 'user-1',
      channel: 'IN_APP',
      event_type: 'ORDER_CREATED',
      event_key: 'ORDER_CREATED:ord-1',
      title: 'Don hang moi',
      content: 'Don hang ORD-20260310-0001 da duoc tao',
      status: 'SENT',
      read: false,
      created_at: '2026-03-10T10:00:00Z',
    },
  ],
  pagination: { page: 1, page_size: 50, total: 1, total_pages: 1 },
}

test.beforeEach(async ({ page }) => {
  await page.route('**/api/v1/auth/login', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: {
          access_token: 'token',
          refresh_token: 'refresh',
          token_type: 'Bearer',
          expires_in: 3600,
          user,
        },
      }),
    })
  })

  await page.route('**/api/v1/products**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: products }),
    })
  })

  await page.route('**/api/v1/orders**', async (route) => {
    if (route.request().method() === 'POST') {
      await route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, data: orderSummary }),
      })
      return
    }

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: {
          items: [orderSummary],
          pagination: { page: 1, page_size: 10, total: 1, total_pages: 1 },
        },
      }),
    })
  })

  await page.route('**/api/v1/orders/ord-1/payments', async (route) => {
    await route.fulfill({
      status: 201,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: {
          payment_id: 'pay-1',
          status: 'PENDING',
          checkout_url: 'http://127.0.0.1:3001/checkout/return?orderId=ord-1&status=PAID',
        },
      }),
    })
  })

  await page.route('**/api/v1/orders/ord-1/tracking', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: tracking }),
    })
  })

  await page.route('**/api/v1/orders/ord-1/cancel', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: {
          ...orderSummary,
          order_status: 'CANCELLED',
        },
      }),
    })
  })

  await page.route('**/api/v1/orders/ord-1', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: orderDetail }),
    })
  })

  await page.route('**/api/v1/notifications/read-all', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: { updated_count: 1 } }),
    })
  })

  await page.route('**/api/v1/notifications/*/read', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        success: true,
        data: { ...notifications.items[0], read: true },
      }),
    })
  })

  await page.route('**/api/v1/notifications**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: notifications }),
    })
  })
})

async function login(page: Parameters<typeof test>[0]['page']) {
  await page.goto('/login')
  await page.getByTestId('customer-login-email').fill('user@golf.local')
  await page.getByTestId('customer-login-password').fill('123456')
  await page.getByTestId('customer-login-submit').click()
  await expect(page.getByTestId('header-cart-link')).toBeVisible()
}

test('customer can log in and browse products', async ({ page }) => {
  await login(page)
  await page.goto('/products')
  await expect(page.getByText('Driver X')).toBeVisible()
  await page.getByTestId('product-add-prod-1').click()
  await page.getByTestId('header-cart-link').click()
  await expect(page.getByTestId('cart-title')).toBeVisible()
  await expect(page.getByText('Driver X')).toBeVisible()
})

test('customer can complete COD checkout', async ({ page }) => {
  await login(page)
  await page.goto('/products')
  await page.getByTestId('product-add-prod-1').click()
  await page.getByTestId('header-cart-link').click()
  await page.getByTestId('cart-checkout-button').click()

  await page.getByTestId('address-recipient-name').fill('Nguyen Van A')
  await page.getByTestId('address-phone').fill('0901234567')
  await page.getByTestId('address-line1').fill('123 Nguyen Hue')
  await page.getByTestId('address-district').fill('District 1')
  await page.getByTestId('address-city').fill('Ho Chi Minh')
  await page.getByTestId('checkout-next-button').click()
  await page.getByTestId('checkout-next-button').click()

  await expect(page.getByTestId('payment-return-result')).toBeVisible()
})

test('customer can view order detail and tracking timeline', async ({ page }) => {
  await login(page)
  await page.goto('/orders')
  await expect(page.getByText('ORD-20260310-0001')).toBeVisible()
  await page.getByRole('button', { name: /Xem chi tiết/i }).click()
  await expect(page.getByText('Payment confirmed')).toBeVisible()
  await expect(page.getByTestId('retry-payment-button')).toBeVisible()
})

test('customer can view notifications and mark all as read', async ({ page }) => {
  await login(page)
  await page.goto('/notifications')
  await expect(page.getByTestId('notifications-title')).toBeVisible()
  await expect(page.getByTestId('notification-item-noti-1')).toBeVisible()
  await page.getByTestId('notifications-mark-all').click()
  await expect(page.getByText('Don hang moi')).toBeVisible()
})
