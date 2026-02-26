import type { ReactElement } from "react";
import { Navigate, useLocation } from "react-router-dom";

import { useAuth } from "../state/auth";

export const RequireAuth = ({
  children,
  role,
}: {
  children: ReactElement;
  role?: "ADMIN" | "USER";
}) => {
  const location = useLocation();
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  if (role && user?.role !== role) {
    return <Navigate to="/login" replace />;
  }

  return children;
};
