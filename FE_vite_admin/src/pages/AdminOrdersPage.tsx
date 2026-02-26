import { useEffect, useState } from "react";

import { apiClient } from "../api";
import type { OrderStatus, OrderSummary } from "../api";
import { getEnvelopeItems } from "../utils/api";

const orderStatusOptions: OrderStatus[] = [
  "NEW",
  "PENDING_PAYMENT",
  "PAID",
  "PROCESSING",
  "SHIPPING",
  "COMPLETED",
  "CANCELLED",
];

export const AdminOrdersPage = () => {
  const [orders, setOrders] = useState<OrderSummary[]>([]);
  const [statusFilter, setStatusFilter] = useState<OrderStatus | "ALL">("ALL");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const fetchOrders = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await apiClient.adminOrders.adminOrdersGet(
        statusFilter === "ALL" ? undefined : statusFilter,
        undefined,
        undefined,
        1,
        100,
      );
      setOrders(getEnvelopeItems<OrderSummary>(response.data));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Không tải được danh sách đơn hàng.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchOrders();
  }, [statusFilter]);

  const updateStatus = async (orderId: string, toStatus: OrderStatus) => {
    setError(null);
    try {
      await apiClient.adminOrders.adminOrdersOrderIdStatusPatch(orderId, {
        to_status: toStatus,
        reason: "updated by admin dashboard",
      });
      await fetchOrders();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Không thể cập nhật trạng thái đơn.");
    }
  };

  return (
    <section className="card">
      <div className="section-header">
        <h2>Quản lý đơn hàng</h2>
        <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value as OrderStatus | "ALL")}>
          <option value="ALL">Tất cả</option>
          {orderStatusOptions.map((status) => (
            <option key={status} value={status}>
              {status}
            </option>
          ))}
        </select>
      </div>
      {loading && <p className="muted">Đang tải dữ liệu...</p>}
      {error && <p className="error">{error}</p>}
      <table>
        <thead>
          <tr>
            <th>Mã đơn</th>
            <th>Order Status</th>
            <th>Payment Status</th>
            <th>Tổng tiền</th>
            <th>Thời điểm đặt</th>
            <th>Cập nhật trạng thái</th>
          </tr>
        </thead>
        <tbody>
          {orders.map((order) => (
            <tr key={order.id}>
              <td>{order.order_code}</td>
              <td>{order.order_status}</td>
              <td>{order.payment_status}</td>
              <td>{order.total_amount}</td>
              <td>{new Date(order.placed_at).toLocaleString()}</td>
              <td>
                <select defaultValue={order.order_status} onChange={(event) => void updateStatus(order.id, event.target.value as OrderStatus)}>
                  {orderStatusOptions.map((status) => (
                    <option key={status} value={status}>
                      {status}
                    </option>
                  ))}
                </select>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
};
