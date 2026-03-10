import { useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Row, Col, Spin, Empty, Pagination, Space, Typography } from 'antd'
import { useQuery } from '@tanstack/react-query'
import ProductCard from '@/components/product/product-card'
import ProductFilters from '@/components/product/product-filters'
import { productsApi } from '@/services/products'
import type { ProductListParams } from '@/services/products'

const { Title } = Typography

export default function ProductListPage() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [filters, setFilters] = useState<ProductListParams>({
    page: Number(searchParams.get('page')) || 1,
    page_size: 12,
    search: searchParams.get('search') || undefined,
    status: 'ACTIVE', // Only show active products
  })

  const { data, isLoading } = useQuery({
    queryKey: ['products', filters],
    queryFn: () => productsApi.list(filters),
  })

  const handleFilterChange = (newFilters: Partial<ProductListParams>) => {
    setFilters((prev) => ({ ...prev, ...newFilters, page: 1 }))
  }

  const handlePageChange = (page: number) => {
    setFilters((prev) => ({ ...prev, page }))
    setSearchParams({ page: String(page) })
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  return (
    <div style={{ padding: 24, maxWidth: 1400, margin: '0 auto' }}>
      <Title level={2}>Sản phẩm</Title>

      <Row gutter={[24, 24]}>
        {/* Filters */}
        <Col xs={24} md={6}>
          <ProductFilters onFilterChange={handleFilterChange} />
        </Col>

        {/* Products */}
        <Col xs={24} md={18}>
          {isLoading ? (
            <div style={{ textAlign: 'center', padding: 48 }}>
              <Spin size="large" />
            </div>
          ) : !data || data.items.length === 0 ? (
            <Empty description="Không tìm thấy sản phẩm" />
          ) : (
            <Space direction="vertical" style={{ width: '100%' }} size="large">
              <Row gutter={[16, 16]}>
                {data.items.map((product) => (
                  <Col xs={24} sm={12} lg={8} key={product.id}>
                    <ProductCard product={product} />
                  </Col>
                ))}
              </Row>

              {data.pagination.total_pages > 1 && (
                <div style={{ textAlign: 'center', marginTop: 24 }}>
                  <Pagination
                    current={data.pagination.page}
                    pageSize={data.pagination.page_size}
                    total={data.pagination.total}
                    onChange={handlePageChange}
                    showSizeChanger={false}
                    showTotal={(total) => `Tổng ${total} sản phẩm`}
                  />
                </div>
              )}
            </Space>
          )}
        </Col>
      </Row>
    </div>
  )
}
