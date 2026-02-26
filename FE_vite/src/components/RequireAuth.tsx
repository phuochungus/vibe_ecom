import { Navigate } from "react-router-dom";

import { useAuth } from "../state/auth";

export const RequireAuth = ({
  children,
  role,
}: {
  children: JSX.Element;
  role?: "ADMIN" | "USER";
}) => {
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (role && user?.role !== role) {
    return <Navigate to="/" replace />;
  }

  return children;
};
