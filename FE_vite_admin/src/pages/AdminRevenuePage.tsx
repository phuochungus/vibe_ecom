import { useEffect, useState } from "react";

import { apiClient } from "../api";
import type { RevenueOrder, RevenueSummary } from "../api";
import type { AdminRevenueSummaryGetGroupByEnum } from "../api/generated/services/admin-revenue-api";
import { getEnvelopeData, getEnvelopeItems } from "../utils/api";

const defaultFrom = new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().slice(0, 10);
const defaultTo = new Date().toISOString().slice(0, 10);

export const AdminRevenuePage = () => {
  const [from, setFrom] = useState(defaultFrom);
  const [to, setTo] = useState(defaultTo);
  const [groupBy, setGroupBy] = useState<AdminRevenueSummaryGetGroupByEnum>("day");
  const [summary, setSummary] = useState<RevenueSummary | null>(null);
  const [orders, setOrders] = useState<RevenueOrder[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const fetchRevenue = async () => {
    setLoading(true);
    setError(null);
    try {
      const [summaryResponse, ordersResponse] = await Promise.all([
        apiClient.adminRevenue.adminRevenueSummaryGet(`${from}T00:00:00.000Z`, `${to}T23:59:59.000Z`, groupBy),
        apiClient.adminRevenue.adminRevenueOrdersGet(`${from}T00:00:00.000Z`, `${to}T23:59:59.000Z`, 1, 20),
      ]);

      setSummary(getEnvelopeData<RevenueSummary>(summaryResponse.data) ?? null);
      setOrders(getEnvelopeItems<RevenueOrder>(ordersResponse.data));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Không tải được dữ liệu doanh thu.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void fetchRevenue();
  }, []);

  return (
    <section className="page-grid">
      <article className="card">
        <div className="section-header">
          <h2>Báo cáo doanh thu</h2>
          <button onClick={() => void fetchRevenue()}>Làm mới</button>
        </div>
        <div className="filters-inline">
          <label>
            From
            <input type="date" value={from} onChange={(event) => setFrom(event.target.value)} />
          </label>
          <label>
            To
            <input type="date" value={to} onChange={(event) => setTo(event.target.value)} />
          </label>
          <label>
            Group by
            <select value={groupBy} onChange={(event) => setGroupBy(event.target.value as AdminRevenueSummaryGetGroupByEnum)}>
              <option value="day">day</option>
              <option value="month">month</option>
              <option value="year">year</option>
            </select>
          </label>
        </div>
        {loading && <p className="muted">Đang tải dữ liệu...</p>}
        {error && <p className="error">{error}</p>}
        {summary && (
          <div className="kpi-grid">
            <div className="kpi-card">
              <span>Gross Revenue</span>
              <strong>{summary.gross_revenue}</strong>
            </div>
            <div className="kpi-card">
              <span>Refund Amount</span>
              <strong>{summary.refund_amount}</strong>
            </div>
            <div className="kpi-card">
              <span>Net Revenue</span>
              <strong>{summary.net_revenue}</strong>
            </div>
            <div className="kpi-card">
              <span>Completed Orders</span>
              <strong>{summary.completed_orders}</strong>
            </div>
            <div className="kpi-card">
              <span>Payment Success Rate</span>
              <strong>{summary.payment_success_rate}</strong>
            </div>
          </div>
        )}
      </article>

      <article className="card">
        <h3>Đơn hàng dùng để đối soát doanh thu</h3>
        <table>
          <thead>
            <tr>
              <th>Order Code</th>
              <th>Order Status</th>
              <th>Payment Status</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order) => (
              <tr key={order.id}>
                <td>{order.order_code}</td>
                <td>{order.order_status}</td>
                <td>{order.payment_status}</td>
                <td>{order.total_amount}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </article>
    </section>
  );
};
