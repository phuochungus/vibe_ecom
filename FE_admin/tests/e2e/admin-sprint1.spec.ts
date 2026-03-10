import { test, expect } from '@playwright/test'

const adminUser = {
  id: 'admin-1',
  full_name: 'Admin Golf',
  email: 'admin@golf.local',
  phone: '0909000000',
  role: 'ADMIN',
  status: 'ACTIVE',
}

const productList = {
  items: [
    {
      id: 'prod-1',
      sku: 'CLUB-001',
      name: 'Driver X',
      price: '1590000',
      stock: 12,
      status: 'ACTIVE',
      image_url: '',
      created_at: '2026-03-10T00:00:00Z',
      updated_at: '2026-03-10T00:00:00Z',
    },
  ],
  pagination: { page: 1, page_size: 20, total: 1, total_pages: 1 },
}

const adminOrders = {
  items: [
    {
      id: 'ord-1',
      order_code: 'ORD-20260310-0001',
      order_status: 'PAID',
      payment_status: 'PAID',
      subtotal_amount: '1590000',
      discount_amount: '0',
      shipping_fee: '30000',
      total_amount: '1620000',
      placed_at: '2026-03-10T10:00:00Z',
      shipping_recipient_name: 'Nguyen Van A',
      shipping_phone: '0901234567',
    },
  ],
  pagination: { page: 1, page_size: 20, total: 1, total_pages: 1 },
}

const orderDetail = {
  ...adminOrders.items[0],
  shipping_line1: '123 Nguyen Hue',
  shipping_line2: '',
  shipping_ward: 'Ben Nghe',
  shipping_district: 'District 1',
  shipping_city: 'Ho Chi Minh',
  items: [
    {
      id: 'item-1',
      product_id: 'prod-1',
      product_sku_snapshot: 'CLUB-001',
      product_name_snapshot: 'Driver X',
      unit_price: '1590000',
      quantity: 1,
      line_total: '1590000',
    },
  ],
}

const tracking = {
  order_id: 'ord-1',
  current_status: 'PAID',
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

const revenueSummary = {
  from: '2026-03-01T00:00:00Z',
  to: '2026-03-10T23:59:59Z',
  gross_revenue: '50000000',
  refund_amount: '2000000',
  net_revenue: '48000000',
  completed_orders: 32,
  payment_success_rate: '0.93',
}

const revenueOrders = {
  items: [
    {
      id: 'ord-1',
      order_code: 'ORD-20260310-0001',
      order_status: 'COMPLETED',
      total_amount: '1620000',
      placed_at: '2026-03-10T10:00:00Z',
      shipping_recipient_name: 'Nguyen Van A',
    },
  ],
  pagination: { page: 1, page_size: 15, total: 1, total_pages: 1 },
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
          user: adminUser,
        },
      }),
    })
  })

  await page.route('**/api/v1/products**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: productList }),
    })
  })

  await page.route('**/api/v1/admin/products', async (route) => {
    if (route.request().method() === 'POST') {
      const payload = JSON.parse(route.request().postData() || '{}')
      await route.fulfill({
        status: 201,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, data: { id: 'prod-2', created_at: '', updated_at: '', ...payload } }),
      })
      return
    }
    await route.fallback()
  })

  await page.route('**/api/v1/admin/orders**', async (route) => {
    const url = route.request().url()
    if (route.request().method() === 'PATCH' && url.endsWith('/api/v1/admin/orders/ord-1/status')) {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, data: { ...adminOrders.items[0], order_status: 'PROCESSING' } }),
      })
      return
    }
    if (route.request().method() === 'GET' && url.endsWith('/api/v1/admin/orders/ord-1')) {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, data: orderDetail }),
      })
      return
    }
    if (route.request().method() === 'GET') {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({ success: true, data: adminOrders }),
      })
      return
    }
    await route.fallback()
  })

  await page.route('**/api/v1/orders/ord-1/tracking', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: tracking }),
    })
  })

  await page.route('**/api/v1/admin/revenue/summary**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: revenueSummary }),
    })
  })

  await page.route('**/api/v1/admin/revenue/orders**', async (route) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ success: true, data: revenueOrders }),
    })
  })
})

async function login(page: Parameters<typeof test>[0]['page']) {
  await page.goto('/login')
  await page.getByTestId('login-identifier').fill('admin@golf.local')
  await page.getByTestId('login-password').fill('admin123')
  await page.getByTestId('login-submit').click()
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
}

test('admin can log in and see dashboard data', async ({ page }) => {
  await login(page)
  await expect(page.getByText('Admin Golf')).toBeVisible()
  await expect(page.getByText('ORD-20260310-0001')).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
  await expect(page.getByText('Đơn hoàn tất')).toBeVisible()
})

test('admin can create a product', async ({ page }) => {
  await login(page)
  await page.getByRole('link', { name: 'Sản phẩm' }).click()
  await page.getByRole('link', { name: 'Thêm sản phẩm' }).click()
  await page.getByLabel(/SKU/).fill('CLUB-002')
  await page.getByLabel(/Tên sản phẩm/).fill('Putter Elite')
  await page.getByLabel(/Giá bán/).fill('2100000')
  await page.getByLabel(/Tồn kho/).fill('5')
  await page.getByRole('button', { name: 'Tạo sản phẩm' }).click()
  await expect(page.getByRole('heading', { name: 'Sản phẩm' })).toBeVisible()
})

test('admin can review order timeline and update status with reason', async ({ page }) => {
  await login(page)
  await page.getByRole('link', { name: 'Đơn hàng' }).click()
  await page.locator('a[href="/orders/ord-1"]').click()
  await expect(page.getByText('Payment confirmed')).toBeVisible()
  await page.getByTestId('order-status-select').click()
  await page.getByRole('option', { name: 'Đang xử lý' }).click()
  await page.getByTestId('order-status-submit').click()
  await page.getByTestId('order-status-reason').fill('verified payment')
  await page.getByRole('button', { name: 'Xác nhận' }).click()
  await expect(page.getByRole('heading', { name: 'ORD-20260310-0001' })).toBeVisible()
})

test('admin can open revenue report', async ({ page }) => {
  await login(page)
  await page.getByRole('link', { name: 'Doanh thu' }).click()
  await expect(page.getByText('Tỷ lệ TT thành công')).toBeVisible()
  await expect(page.getByText('ORD-20260310-0001')).toBeVisible()
  await expect(page.getByText('Nguyen Van A')).toBeVisible()
})
