const ACCESS_TOKEN_KEY = "golf_store_access_token";
const REFRESH_TOKEN_KEY = "golf_store_refresh_token";

let accessTokenMemory: string | null = null;
let refreshTokenMemory: string | null = null;

const hasLocalStorage = (): boolean => {
  return typeof window !== "undefined" && typeof window.localStorage !== "undefined";
};

export const tokenStore = {
  getAccessToken(): string | null {
    if (hasLocalStorage()) {
      return window.localStorage.getItem(ACCESS_TOKEN_KEY);
    }
    return accessTokenMemory;
  },

  setAccessToken(token: string): void {
    accessTokenMemory = token;
    if (hasLocalStorage()) {
      window.localStorage.setItem(ACCESS_TOKEN_KEY, token);
    }
  },

  getRefreshToken(): string | null {
    if (hasLocalStorage()) {
      return window.localStorage.getItem(REFRESH_TOKEN_KEY);
    }
    return refreshTokenMemory;
  },

  setRefreshToken(token: string): void {
    refreshTokenMemory = token;
    if (hasLocalStorage()) {
      window.localStorage.setItem(REFRESH_TOKEN_KEY, token);
    }
  },

  clear(): void {
    accessTokenMemory = null;
    refreshTokenMemory = null;
    if (hasLocalStorage()) {
      window.localStorage.removeItem(ACCESS_TOKEN_KEY);
      window.localStorage.removeItem(REFRESH_TOKEN_KEY);
    }
  },
};
