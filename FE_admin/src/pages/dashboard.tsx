import { useQuery } from '@tanstack/react-query'
import { format, subDays } from 'date-fns'
import {
  TrendingUp,
  ShoppingCart,
  Package,
  CheckCircle,
  ArrowUpRight,
} from 'lucide-react'
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { revenueApi } from '@/services/revenue'
import { ordersApi } from '@/services/orders'
import { formatCurrency, formatOrderStatus, getOrderStatusVariant } from '@/lib/utils'

function StatCard({
  title,
  value,
  sub,
  icon: Icon,
  loading,
}: {
  title: string
  value: string
  sub?: string
  icon: React.ElementType
  loading?: boolean
}) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <div className="rounded-md bg-primary/10 p-2">
          <Icon className="h-4 w-4 text-primary" />
        </div>
      </CardHeader>
      <CardContent>
        {loading ? (
          <Skeleton className="h-8 w-28" />
        ) : (
          <div className="text-2xl font-bold">{value}</div>
        )}
        {sub && <p className="mt-1 text-xs text-muted-foreground">{sub}</p>}
      </CardContent>
    </Card>
  )
}

export default function DashboardPage() {
  const from = format(subDays(new Date(), 29), "yyyy-MM-dd'T'00:00:00'Z'")
  const to = format(new Date(), "yyyy-MM-dd'T'23:59:59'Z'")

  const { data: summary, isLoading: loadingSummary } = useQuery({
    queryKey: ['revenue-summary', from, to],
    queryFn: () => revenueApi.summary({ from, to, group_by: 'day' }),
  })

  const { data: latestOrders, isLoading: loadingOrders } = useQuery({
    queryKey: ['admin-orders-latest'],
    queryFn: () => ordersApi.adminList({ page: 1, page_size: 5 }),
  })

  const successRate = summary
    ? `${(parseFloat(summary.payment_success_rate) * 100).toFixed(1)}%`
    : '—'

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Dashboard</h1>
        <p className="text-sm text-muted-foreground">Tổng quan 30 ngày gần nhất</p>
      </div>

      {/* KPI cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Doanh thu thuần"
          value={summary ? formatCurrency(summary.net_revenue) : '—'}
          icon={TrendingUp}
          loading={loadingSummary}
        />
        <StatCard
          title="Đơn hoàn tất"
          value={summary ? summary.completed_orders.toString() : '—'}
          icon={CheckCircle}
          loading={loadingSummary}
        />
        <StatCard
          title="Tỷ lệ thanh toán thành công"
          value={successRate}
          icon={ArrowUpRight}
          loading={loadingSummary}
        />
        <StatCard
          title="Tổng đơn hàng (trang hiện tại)"
          value={latestOrders ? latestOrders.pagination.total.toString() : '—'}
          icon={ShoppingCart}
          loading={loadingOrders}
        />
      </div>

      {/* Revenue placeholder chart */}
      <Card>
        <CardHeader>
          <CardTitle>Doanh thu</CardTitle>
          <CardDescription>30 ngày gần nhất</CardDescription>
        </CardHeader>
        <CardContent>
          {loadingSummary ? (
            <Skeleton className="h-56 w-full" />
          ) : (
            <ResponsiveContainer width="100%" height={220}>
              <AreaChart
                data={[
                  { date: 'Hiện tại', net: parseFloat(summary?.net_revenue ?? '0'), gross: parseFloat(summary?.gross_revenue ?? '0') },
                ]}
              >
                <defs>
                  <linearGradient id="colorNet" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
                <XAxis dataKey="date" className="text-xs fill-muted-foreground" />
                <YAxis className="text-xs fill-muted-foreground" tickFormatter={(v) => `${(v / 1e6).toFixed(0)}M`} />
                <Tooltip formatter={(v: number | string | undefined) => v !== undefined ? formatCurrency(v.toString()) : '—'} />
                <Area
                  type="monotone"
                  dataKey="net"
                  stroke="hsl(var(--primary))"
                  fill="url(#colorNet)"
                  strokeWidth={2}
                  name="Doanh thu thuần"
                />
              </AreaChart>
            </ResponsiveContainer>
          )}
        </CardContent>
      </Card>

      {/* Recent orders */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Đơn hàng gần nhất</CardTitle>
            <CardDescription>5 đơn hàng mới nhất</CardDescription>
          </div>
          <Package className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          {loadingOrders ? (
            <div className="space-y-3">
              {Array.from({ length: 5 }).map((_, i) => (
                <Skeleton key={i} className="h-10 w-full" />
              ))}
            </div>
          ) : (
            <div className="space-y-3">
              {latestOrders?.items.map((order) => (
                <div
                  key={order.id}
                  className="flex items-center justify-between rounded-md border p-3 text-sm"
                >
                  <div>
                    <p className="font-medium">{order.order_code}</p>
                    <p className="text-xs text-muted-foreground">{order.shipping_recipient_name}</p>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="font-semibold">{formatCurrency(order.total_amount)}</span>
                    <Badge variant={getOrderStatusVariant(order.order_status)}>
                      {formatOrderStatus(order.order_status)}
                    </Badge>
                  </div>
                </div>
              ))}
              {(!latestOrders?.items || latestOrders.items.length === 0) && (
                <p className="text-center text-sm text-muted-foreground py-8">Chưa có đơn hàng nào</p>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
