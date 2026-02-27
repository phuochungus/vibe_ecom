import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { format } from 'date-fns'
import { Eye, Search } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { ordersApi } from '@/services/orders'
import {
  formatCurrency,
  formatOrderStatus,
  getOrderStatusVariant,
} from '@/lib/utils'

const PAGE_SIZE = 20

const PAYMENT_STATUS_LABEL: Record<string, string> = {
  UNPAID: 'Chưa TT',
  PAID: 'Đã TT',
  FAILED: 'TT Thất bại',
  REFUNDED: 'Hoàn tiền',
}

export default function OrderListPage() {
  const [status, setStatus] = useState('all')
  const [q, setQ] = useState('')
  const [page, setPage] = useState(1)

  const { data, isLoading } = useQuery({
    queryKey: ['admin-orders', status, page],
    queryFn: () =>
      ordersApi.adminList({
        status: status === 'all' ? undefined : status,
        page,
        page_size: PAGE_SIZE,
      }),
  })

  const filtered = q
    ? data?.items.filter(
        (o) =>
          o.order_code.toLowerCase().includes(q.toLowerCase()) ||
          o.shipping_recipient_name.toLowerCase().includes(q.toLowerCase())
      )
    : data?.items

  const totalPages = data?.pagination.total_pages ?? 1

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-2xl font-bold">Đơn hàng</h1>
        <p className="text-sm text-muted-foreground">
          Quản lý và xử lý đơn hàng ({data?.pagination.total ?? 0} đơn)
        </p>
      </div>

      {/* Filters */}
      <div className="flex gap-3">
        <div className="relative flex-1 max-w-sm">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Mã đơn / Tên khách..."
            value={q}
            onChange={(e) => setQ(e.target.value)}
            className="pl-9"
          />
        </div>
        <Select value={status} onValueChange={(v: string) => { setStatus(v); setPage(1) }}>
          <SelectTrigger className="w-52">
            <SelectValue placeholder="Trạng thái" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">Tất cả trạng thái</SelectItem>
            <SelectItem value="NEW">Mới tạo</SelectItem>
            <SelectItem value="PENDING_PAYMENT">Chờ thanh toán</SelectItem>
            <SelectItem value="PAID">Đã thanh toán</SelectItem>
            <SelectItem value="PROCESSING">Đang xử lý</SelectItem>
            <SelectItem value="SHIPPING">Đang giao</SelectItem>
            <SelectItem value="COMPLETED">Hoàn tất</SelectItem>
            <SelectItem value="CANCELLED">Đã hủy</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {/* Table */}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Mã đơn</TableHead>
              <TableHead>Khách hàng</TableHead>
              <TableHead className="text-right">Tổng tiền</TableHead>
              <TableHead>Trạng thái</TableHead>
              <TableHead>Thanh toán</TableHead>
              <TableHead>Ngày tạo</TableHead>
              <TableHead className="w-12"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading
              ? Array.from({ length: 8 }).map((_, i) => (
                  <TableRow key={i}>
                    {Array.from({ length: 7 }).map((_, j) => (
                      <TableCell key={j}>
                        <Skeleton className="h-5 w-full" />
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              : (filtered ?? []).map((order) => (
                  <TableRow key={order.id} className="group">
                    <TableCell className="font-mono text-xs font-medium">
                      {order.order_code}
                    </TableCell>
                    <TableCell>
                      <div>
                        <p className="font-medium text-sm">{order.shipping_recipient_name}</p>
                        <p className="text-xs text-muted-foreground">{order.shipping_phone}</p>
                      </div>
                    </TableCell>
                    <TableCell className="text-right font-semibold">
                      {formatCurrency(order.total_amount)}
                    </TableCell>
                    <TableCell>
                      <Badge variant={getOrderStatusVariant(order.order_status)}>
                        {formatOrderStatus(order.order_status)}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <span className="text-xs text-muted-foreground">
                        {PAYMENT_STATUS_LABEL[order.payment_status] ?? order.payment_status}
                      </span>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {format(new Date(order.placed_at), 'dd/MM/yyyy HH:mm')}
                    </TableCell>
                    <TableCell>
                      <Button variant="ghost" size="icon" asChild>
                        <Link to={`/orders/${order.id}`}>
                          <Eye className="h-4 w-4" />
                        </Link>
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
            {!isLoading && (filtered ?? []).length === 0 && (
              <TableRow>
                <TableCell colSpan={7} className="h-40 text-center text-muted-foreground">
                  Không có đơn hàng nào
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
            Trang {page} / {totalPages} ({data?.pagination.total} đơn)
          </p>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1}>
              Trước
            </Button>
            <Button variant="outline" size="sm" onClick={() => setPage((p) => Math.min(totalPages, p + 1))} disabled={page === totalPages}>
              Sau
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
