import { useState } from "react";
import type { FormEvent } from "react";
import { Navigate, useLocation, useNavigate } from "react-router-dom";

import { useAuth } from "../state/auth";

export const LoginPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, isAuthenticated, user } = useAuth();

  const [identifier, setIdentifier] = useState("admin@example.com");
  const [password, setPassword] = useState("123456");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  if (isAuthenticated && user?.role === "ADMIN") {
    return <Navigate to="/admin/products" replace />;
  }

  const onSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await login({ identifier, password });
      const nextPath = (location.state as { from?: string } | null)?.from ?? "/admin/products";
      navigate(nextPath, { replace: true });
    } catch (err) {
      setError(err instanceof Error ? err.message : "Đăng nhập thất bại.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <section className="card login-card">
      <h1>Admin Login</h1>
      <p className="muted">Dùng tài khoản ADMIN để quản lý sản phẩm, đơn hàng và doanh thu.</p>
      <form onSubmit={onSubmit} className="form-grid">
        <label>
          Email / Phone
          <input value={identifier} onChange={(event) => setIdentifier(event.target.value)} required />
        </label>
        <label>
          Mật khẩu
          <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} required />
        </label>
        {error && <p className="error">{error}</p>}
        <button type="submit" disabled={loading}>
          {loading ? "Đang đăng nhập..." : "Đăng nhập"}
        </button>
      </form>
    </section>
  );
};
