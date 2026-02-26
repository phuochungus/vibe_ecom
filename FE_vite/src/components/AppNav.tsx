import { Link, NavLink } from "react-router-dom";

import { useAuth } from "../state/auth";
import { useCart } from "../state/cart";

export const AppNav = () => {
  const { user, logout, isAuthenticated } = useAuth();
  const { items } = useCart();

  return (
    <header className="app-nav">
      <div className="app-nav__brand">
        <Link to="/">Golf Store</Link>
      </div>
      <nav className="app-nav__links">
        <NavLink to="/">Products</NavLink>
        <NavLink to="/cart">Cart ({items.length})</NavLink>
        {isAuthenticated && <NavLink to="/orders">Orders</NavLink>}
        {isAuthenticated && <NavLink to="/notifications">Notifications</NavLink>}
        {user?.role === "ADMIN" && <NavLink to="/admin/products">Admin Products</NavLink>}
        {user?.role === "ADMIN" && <NavLink to="/admin/orders">Admin Orders</NavLink>}
        {user?.role === "ADMIN" && <NavLink to="/admin/revenue">Admin Revenue</NavLink>}
      </nav>
      <div className="app-nav__actions">
        {isAuthenticated ? (
          <>
            <span className="muted">{user?.fullName}</span>
            <button type="button" onClick={logout}>
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
