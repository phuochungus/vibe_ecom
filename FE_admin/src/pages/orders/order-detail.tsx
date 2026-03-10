import { useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { format } from "date-fns";
import { ArrowLeft, Loader2, CheckCircle2, Clock, Package, Truck, XCircle, Circle } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import {
    AlertDialog,
    AlertDialogAction,
    AlertDialogCancel,
    AlertDialogContent,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { ordersApi } from "@/services/orders";
import { formatCurrency, formatOrderStatus, getOrderStatusVariant, ORDER_TRANSITIONS } from "@/lib/utils";
import type { OrderStatus } from "@/types";

const STATUS_ICONS: Record<OrderStatus, React.ElementType> = {
    NEW: Circle,
    PENDING_PAYMENT: Clock,
    PAID: CheckCircle2,
    PROCESSING: Package,
    SHIPPING: Truck,
    COMPLETED: CheckCircle2,
    CANCELLED: XCircle,
};

export default function OrderDetailPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const [nextStatus, setNextStatus] = useState<OrderStatus | "">("");
    const [statusReason, setStatusReason] = useState("");

    const { data: order, isLoading } = useQuery({
        queryKey: ["admin-order", id],
        queryFn: () => ordersApi.adminGet(id!),
        enabled: !!id,
    });

    const { data: tracking, isLoading: loadingTracking } = useQuery({
        queryKey: ["order-tracking", id],
        queryFn: () => ordersApi.getTracking(id!),
        enabled: !!id,
    });

    const statusMutation = useMutation({
        mutationFn: () => ordersApi.updateStatus(id!, nextStatus as OrderStatus, statusReason || undefined),
        onSuccess: () => {
            toast.success(`Đã chuyển trạng thái sang ${formatOrderStatus(nextStatus as OrderStatus)}`);
            queryClient.invalidateQueries({ queryKey: ["admin-order", id] });
            queryClient.invalidateQueries({ queryKey: ["order-tracking", id] });
            queryClient.invalidateQueries({ queryKey: ["admin-orders"] });
            setNextStatus("");
            setStatusReason("");
        },
        onError: () => toast.error("Chuyển trạng thái thất bại, vui lòng thử lại"),
    });

    const validNextStatuses = order ? ORDER_TRANSITIONS[order.order_status] : [];
    const trackingTimeline =
        tracking?.timeline?.map((event) => ({
            status: event.to_status ?? event.status ?? order?.order_status ?? "NEW",
            actor: event.changed_by_type ?? event.source_type ?? "SYSTEM",
            note: event.change_reason ?? event.description,
            occurred_at: event.occurred_at,
        })) ?? [];

    if (isLoading) {
        return (
            <div className="max-w-3xl space-y-4">
                <Skeleton className="h-10 w-48" />
                <Skeleton className="h-64 w-full" />
                <Skeleton className="h-40 w-full" />
            </div>
        );
    }

    if (!order) return null;

    return (
        <div className="max-w-3xl space-y-6">
            {/* Header */}
            <div className="flex items-center gap-3">
                <Button variant="ghost" size="icon" onClick={() => navigate("/orders")}>
                    <ArrowLeft className="h-4 w-4" />
                </Button>
                <div className="flex-1">
                    <div className="flex items-center gap-3">
                        <h1 className="text-xl font-bold font-mono">{order.order_code}</h1>
                        <Badge variant={getOrderStatusVariant(order.order_status)}>
                            {formatOrderStatus(order.order_status)}
                        </Badge>
                    </div>
                    <p className="text-xs text-muted-foreground">
                        Đặt lúc {format(new Date(order.placed_at), "HH:mm dd/MM/yyyy")}
                    </p>
                </div>

                {/* Status transition */}
                {validNextStatuses.length > 0 && (
                    <AlertDialog>
                        <div className="flex items-center gap-2">
                            <Select value={nextStatus} onValueChange={(v: string) => setNextStatus(v as OrderStatus)}>
                                <SelectTrigger className="w-48" data-testid="order-status-select">
                                    <SelectValue placeholder="Chuyển trạng thái" />
                                </SelectTrigger>
                                <SelectContent>
                                    {validNextStatuses.map((s) => (
                                        <SelectItem key={s} value={s}>
                                            {formatOrderStatus(s)}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                            <AlertDialogTrigger asChild>
                                <Button disabled={!nextStatus} size="sm" data-testid="order-status-submit">
                                    Cập nhật
                                </Button>
                            </AlertDialogTrigger>
                        </div>
                        <AlertDialogContent>
                            <AlertDialogHeader>
                                <AlertDialogTitle>Xác nhận chuyển trạng thái</AlertDialogTitle>
                                <AlertDialogDescription>
                                    Chuyển đơn <strong>{order.order_code}</strong> từ{" "}
                                    <strong>{formatOrderStatus(order.order_status)}</strong> sang{" "}
                                    <strong>{nextStatus ? formatOrderStatus(nextStatus as OrderStatus) : ""}</strong>?
                                </AlertDialogDescription>
                            </AlertDialogHeader>
                            <div className="space-y-2">
                                <p className="text-sm font-medium">Lý do cập nhật</p>
                                <Textarea
                                    value={statusReason}
                                    onChange={(event) => setStatusReason(event.target.value)}
                                    placeholder="verified payment / handed to carrier / ghi chú vận hành..."
                                    data-testid="order-status-reason"
                                />
                            </div>
                            <AlertDialogFooter>
                                <AlertDialogCancel>Hủy</AlertDialogCancel>
                                <AlertDialogAction
                                    onClick={() => statusMutation.mutate()}
                                    disabled={statusMutation.isPending}
                                >
                                    {statusMutation.isPending ? (
                                        <Loader2 className="h-4 w-4 animate-spin" />
                                    ) : (
                                        "Xác nhận"
                                    )}
                                </AlertDialogAction>
                            </AlertDialogFooter>
                        </AlertDialogContent>
                    </AlertDialog>
                )}
            </div>

            {/* Order items */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Sản phẩm đặt hàng</CardTitle>
                </CardHeader>
                <CardContent className="space-y-3">
                    {order.items?.map((item) => (
                        <div key={item.id} className="flex items-center justify-between text-sm">
                            <div>
                                <p className="font-medium">{item.product_name_snapshot ?? item.name ?? "Sản phẩm"}</p>
                                <p className="text-xs text-muted-foreground">
                                    SKU: {item.product_sku_snapshot ?? item.sku ?? "—"}
                                </p>
                            </div>
                            <div className="text-right">
                                <p>
                                    {formatCurrency(item.unit_price)} × {item.quantity}
                                </p>
                                <p className="font-semibold">{formatCurrency(item.line_total)}</p>
                            </div>
                        </div>
                    )) ?? <p className="text-sm text-muted-foreground">Không có dữ liệu sản phẩm</p>}
                    <Separator />
                    <div className="space-y-1 text-sm">
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">Tạm tính</span>
                            <span>{formatCurrency(order.subtotal_amount)}</span>
                        </div>
                        <div className="flex justify-between">
                            <span className="text-muted-foreground">Phí vận chuyển</span>
                            <span>{formatCurrency(order.shipping_fee)}</span>
                        </div>
                        {parseFloat(order.discount_amount) > 0 && (
                            <div className="flex justify-between text-green-600">
                                <span>Giảm giá</span>
                                <span>-{formatCurrency(order.discount_amount)}</span>
                            </div>
                        )}
                        <Separator />
                        <div className="flex justify-between font-bold text-base">
                            <span>Tổng cộng</span>
                            <span>{formatCurrency(order.total_amount)}</span>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* Shipping info */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Thông tin giao hàng</CardTitle>
                </CardHeader>
                <CardContent className="text-sm space-y-1">
                    <p className="font-medium">{order.shipping_recipient_name}</p>
                    <p className="text-muted-foreground">{order.shipping_phone}</p>
                    <p className="text-muted-foreground">
                        {[
                            order.shipping_line1,
                            order.shipping_line2,
                            order.shipping_ward,
                            order.shipping_district,
                            order.shipping_city,
                        ]
                            .filter(Boolean)
                            .join(", ")}
                    </p>
                    {order.customer_note && (
                        <p className="mt-2 text-amber-600 dark:text-amber-400">📌 Ghi chú: {order.customer_note}</p>
                    )}
                </CardContent>
            </Card>

            {/* Status timeline */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-base">Lịch sử trạng thái</CardTitle>
                </CardHeader>
                <CardContent className="pt-2">
                    {loadingTracking ? (
                        <div className="space-y-3">
                            {Array.from({ length: 3 }).map((_, i) => (
                                <Skeleton key={i} className="h-10 w-full" />
                            ))}
                        </div>
                    ) : (
                        <div className="relative pl-8 pt-1 space-y-4">
                            <div className="absolute left-3 top-2 bottom-2 w-px bg-border/70" />
                            {trackingTimeline.map((event, idx) => {
                                const Icon = STATUS_ICONS[event.status] ?? Circle;
                                return (
                                    <div key={idx} className="relative px-4 py-2">
                                        <div className="absolute -left-7 top-2.5 bg-background px-1">
                                            <Icon className="h-4 w-4 text-primary" />
                                        </div>
                                        <div className="space-y-1">
                                            <p className="font-medium leading-none tracking-tight">{formatOrderStatus(event.status)}</p>
                                            <p className="text-xs leading-6 text-muted-foreground">
                                                {format(new Date(event.occurred_at), "HH:mm dd/MM/yyyy")} ·{" "}
                                                {event.actor}
                                                {event.note ? ` · ${event.note}` : ""}
                                            </p>
                                        </div>
                                    </div>
                                );
                            }) ?? <p className="text-muted-foreground text-sm">Không có lịch sử</p>}
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
