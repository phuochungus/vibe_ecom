import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { format, subDays, startOfMonth, endOfMonth } from 'date-fns'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Legend,
} from 'recharts'
import { TrendingUp, TrendingDown, CheckCircle, Percent } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Skeleton } from '@/components/ui/skeleton'
import { revenueApi } from '@/services/revenue'
import {
  formatCurrency,
  formatOrderStatus,
  getOrderStatusVariant,
} from '@/lib/utils'

type Period = '7d' | '30d' | 'month'

function getPeriodRange(period: Period): { from: string; to: string } {
  const now = new Date()
  let from: Date
  let to: Date = now
  if (period === '7d') {
    from = subDays(now, 6)
  } else if (period === '30d') {
    from = subDays(now, 29)
  } else {
    from = startOfMonth(now)
    to = endOfMonth(now)
  }
  return {
    from: format(from, "yyyy-MM-dd'T'00:00:00'Z'"),
    to: format(to, "yyyy-MM-dd'T'23:59:59'Z'"),
  }
}

interface SummaryCardProps {
  title: string
  value: string
  icon: React.ElementType
  variant?: 'default' | 'destructive'
  sub?: string
  loading?: boolean
}

function SummaryCard({ title, value, icon: Icon, variant = 'default', sub, loading }: SummaryCardProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">{title}</CardTitle>
        <div className={`rounded-md p-2 ${variant === 'destructive' ? 'bg-destructive/10' : 'bg-primary/10'}`}>
          <Icon className={`h-4 w-4 ${variant === 'destructive' ? 'text-destructive' : 'text-primary'}`} />
        </div>
      </CardHeader>
      <CardContent>
        {loading ? (
          <Skeleton className="h-8 w-28" />
        ) : (
          <div className={`text-2xl font-bold ${variant === 'destructive' ? 'text-destructive' : ''}`}>{value}</div>
        )}
        {sub && <p className="mt-1 text-xs text-muted-foreground">{sub}</p>}
      </CardContent>
    </Card>
  )
}

export default function RevenuePage() {
  const [period, setPeriod] = useState<Period>('30d')
  const [ordersPage, setOrdersPage] = useState(1)

  const { from, to } = getPeriodRange(period)

  const { data: summary, isLoading: loadingSummary } = useQuery({
    queryKey: ['revenue-summary', from, to],
    queryFn: () => revenueApi.summary({ from, to, group_by: 'day' }),
  })

  const { data: ordersData, isLoading: loadingOrders } = useQuery({
    queryKey: ['revenue-orders', from, to, ordersPage],
    queryFn: () => revenueApi.orders({ from, to, page: ordersPage, page_size: 15 }),
  })

  const successRate = summary
    ? `${(parseFloat(summary.payment_success_rate) * 100).toFixed(1)}%`
    : '—'

  const totalPages = ordersData?.pagination.total_pages ?? 1

  // For chart: we only have summary, so show 2-bar comparison
  const chartData = summary
    ? [
        {
          name: 'Kỳ hiện tại',
          'Doanh thu gộp': parseFloat(summary.gross_revenue),
          'Doanh thu thuần': parseFloat(summary.net_revenue),
          'Hoàn tiền': parseFloat(summary.refund_amount),
        },
      ]
    : []

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Báo cáo doanh thu</h1>
          <p className="text-sm text-muted-foreground">
            {format(new Date(from), 'dd/MM/yyyy')} – {format(new Date(to), 'dd/MM/yyyy')}
          </p>
        </div>
        <Tabs value={period} onValueChange={(v: string) => { setPeriod(v as Period); setOrdersPage(1) }}>
          <TabsList>
            <TabsTrigger value="7d">7 ngày</TabsTrigger>
            <TabsTrigger value="30d">30 ngày</TabsTrigger>
            <TabsTrigger value="month">Tháng này</TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* Summary KPIs */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <SummaryCard
          title="Doanh thu gộp"
          value={summary ? formatCurrency(summary.gross_revenue) : '—'}
          icon={TrendingUp}
          loading={loadingSummary}
          sub={`${summary?.completed_orders ?? 0} đơn hoàn tất`}
        />
        <SummaryCard
          title="Doanh thu thuần"
          value={summary ? formatCurrency(summary.net_revenue) : '—'}
          icon={TrendingUp}
          loading={loadingSummary}
          sub="Sau khi trừ hoàn tiền"
        />
        <SummaryCard
          title="Hoàn tiền"
          value={summary ? formatCurrency(summary.refund_amount) : '—'}
          icon={TrendingDown}
          variant="destructive"
          loading={loadingSummary}
        />
        <SummaryCard
          title="Tỷ lệ TT thành công"
          value={successRate}
          icon={Percent}
          loading={loadingSummary}
          sub={summary ? `${summary.completed_orders} đơn hoàn tất` : ''}
        />
      </div>

      {/* Chart */}
      <Card>
        <CardHeader>
          <CardTitle>So sánh doanh thu</CardTitle>
          <CardDescription>Doanh thu gộp / thuần / hoàn tiền</CardDescription>
        </CardHeader>
        <CardContent>
          {loadingSummary ? (
            <Skeleton className="h-56 w-full" />
          ) : (
            <ResponsiveContainer width="100%" height={220}>
              <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" className="stroke-border" />
                <XAxis dataKey="name" className="text-xs fill-muted-foreground" />
                <YAxis
                  className="text-xs fill-muted-foreground"
                  tickFormatter={(v) => `${(v / 1e6).toFixed(0)}M`}
                />
                <Tooltip formatter={(v: number | string | undefined) => v !== undefined ? formatCurrency(v.toString()) : '—'} />
                <Legend />
                <Bar dataKey="Doanh thu gộp" fill="hsl(var(--chart-1))" radius={[4, 4, 0, 0]} />
                <Bar dataKey="Doanh thu thuần" fill="hsl(var(--chart-2))" radius={[4, 4, 0, 0]} />
                <Bar dataKey="Hoàn tiền" fill="hsl(var(--destructive))" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          )}
        </CardContent>
      </Card>

      {/* Revenue orders table */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle className="text-base">Đơn hàng hoàn tất</CardTitle>
            <CardDescription>
              {ordersData?.pagination.total ?? 0} đơn trong kỳ · Truy vết tới mã đơn (BR-REV-03)
            </CardDescription>
          </div>
          <CheckCircle className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Mã đơn</TableHead>
                  <TableHead>Khách hàng</TableHead>
                  <TableHead className="text-right">Tổng tiền</TableHead>
                  <TableHead>Trạng thái</TableHead>
                  <TableHead>Ngày hoàn tất</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loadingOrders
                  ? Array.from({ length: 5 }).map((_, i) => (
                      <TableRow key={i}>
                        {Array.from({ length: 5 }).map((_, j) => (
                          <TableCell key={j}>
                            <Skeleton className="h-4 w-full" />
                          </TableCell>
                        ))}
                      </TableRow>
                    ))
                  : ordersData?.items.map((order) => (
                      <TableRow key={order.id}>
                        <TableCell className="font-mono text-xs font-medium">
                          {order.order_code}
                        </TableCell>
                        <TableCell className="text-sm">{order.shipping_recipient_name}</TableCell>
                        <TableCell className="text-right font-semibold text-sm">
                          {formatCurrency(order.total_amount)}
                        </TableCell>
                        <TableCell>
                          <Badge variant={getOrderStatusVariant(order.order_status)}>
                            {formatOrderStatus(order.order_status)}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-xs text-muted-foreground">
                          {format(new Date(order.placed_at), 'dd/MM/yyyy')}
                        </TableCell>
                      </TableRow>
                    ))}
                {!loadingOrders && (ordersData?.items.length === 0) && (
                  <TableRow>
                    <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                      Không có đơn hàng trong kỳ
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between">
              <p className="text-sm text-muted-foreground">
                Trang {ordersPage} / {totalPages}
              </p>
              <div className="flex gap-2">
                <button
                  className="text-sm underline disabled:opacity-40"
                  onClick={() => setOrdersPage((p) => Math.max(1, p - 1))}
                  disabled={ordersPage === 1}
                >
                  Trước
                </button>
                <button
                  className="text-sm underline disabled:opacity-40"
                  onClick={() => setOrdersPage((p) => Math.min(totalPages, p + 1))}
                  disabled={ordersPage === totalPages}
                >
                  Sau
                </button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
