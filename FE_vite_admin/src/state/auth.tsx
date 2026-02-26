import { createContext, useContext, useEffect, useMemo, useState } from "react";

import { apiClient, apiSession } from "../api";

type AppRole = "USER" | "ADMIN";

type AppUser = {
  id: string;
  role: AppRole;
  fullName: string;
  email: string;
  phone: string;
};

type LoginPayload = {
  identifier: string;
  password: string;
};

type AuthContextValue = {
  user: AppUser | null;
  isAuthenticated: boolean;
  login: (payload: LoginPayload) => Promise<void>;
  logout: () => Promise<void>;
};

const AUTH_USER_KEY = "golf_store_user";

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<AppUser | null>(null);

  useEffect(() => {
    const raw = window.localStorage.getItem(AUTH_USER_KEY);
    if (raw) {
      try {
        setUser(JSON.parse(raw) as AppUser);
      } catch {
        window.localStorage.removeItem(AUTH_USER_KEY);
      }
    }
  }, []);

  const login = async ({ identifier, password }: LoginPayload) => {
    const response = await apiClient.auth.authLoginPost({ identifier, password });
    const payload = response.data as any;
    const data = payload?.data;

    if (!data?.access_token || !data?.user) {
      throw new Error("Login response is missing token or user data");
    }

    const normalizedUser: AppUser = {
      id: data.user.id,
      role: data.user.role,
      fullName: data.user.full_name,
      email: data.user.email,
      phone: data.user.phone,
    };

    apiSession.setAccessToken(data.access_token);
    if (data.refresh_token) {
      apiSession.setRefreshToken(data.refresh_token);
    }

    setUser(normalizedUser);
    window.localStorage.setItem(AUTH_USER_KEY, JSON.stringify(normalizedUser));
  };

  const logout = async () => {
    try {
      await apiClient.auth.authLogoutPost();
    } catch {
      // Ignore logout errors and clear local session anyway.
    }

    apiSession.clear();
    setUser(null);
    window.localStorage.removeItem(AUTH_USER_KEY);
  };

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      isAuthenticated: !!user && !!apiSession.getAccessToken(),
      login,
      logout,
    }),
    [user],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = (): AuthContextValue => {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
};
