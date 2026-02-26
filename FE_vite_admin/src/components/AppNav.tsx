import { Link, NavLink, useNavigate } from "react-router-dom";

import { useAuth } from "../state/auth";

export const AppNav = () => {
  const navigate = useNavigate();
  const { user, logout, isAuthenticated } = useAuth();

  const handleLogout = async () => {
    await logout();
    navigate("/login", { replace: true });
  };

  return (
    <header className="app-nav">
      <div className="app-nav__brand">
        <Link to="/admin/products">Golf Store Admin</Link>
      </div>
      {isAuthenticated && user?.role === "ADMIN" ? (
        <nav className="app-nav__links">
          <NavLink to="/admin/products">Products</NavLink>
          <NavLink to="/admin/orders">Orders</NavLink>
          <NavLink to="/admin/revenue">Revenue</NavLink>
        </nav>
      ) : (
        <div />
      )}
      <div className="app-nav__actions">
        {isAuthenticated ? (
          <>
            <span className="muted">{user?.fullName}</span>
            <button type="button" onClick={handleLogout}>
              Logout
            </button>
          </>
        ) : (
          <Link className="button-link" to="/login">
            Login
          </Link>
        )}
      </div>
    </header>
  );
};
