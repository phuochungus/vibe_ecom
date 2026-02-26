import { Navigate, Route, Routes } from "react-router-dom";

import { AppNav } from "./components/AppNav";
import { RequireAuth } from "./components/RequireAuth";
import { LoginPage } from "./pages/LoginPage";
import { AdminOrdersPage } from "./pages/AdminOrdersPage";
import { AdminProductsPage } from "./pages/AdminProductsPage";
import { AdminRevenuePage } from "./pages/AdminRevenuePage";

import "./App.css";

function App() {
  return (
    <div className="app-shell">
      <AppNav />
      <main className="app-content">
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/admin/products"
            element={
              <RequireAuth role="ADMIN">
                <AdminProductsPage />
              </RequireAuth>
            }
          />
          <Route
            path="/admin/orders"
            element={
              <RequireAuth role="ADMIN">
                <AdminOrdersPage />
              </RequireAuth>
            }
          />
          <Route
            path="/admin/revenue"
            element={
              <RequireAuth role="ADMIN">
                <AdminRevenuePage />
              </RequireAuth>
            }
          />
          <Route path="/" element={<Navigate to="/admin/products" replace />} />
        </Routes>
      </main>
    </div>
  );
}

export default App;
