import axios from "axios";

import { Configuration } from "./generated/configuration";
import { toApiClientError } from "./errors";
import { tokenStore } from "./token-store";

declare global {
  interface Window {
    __API_BASE_URL__?: string;
  }
}

const defaultBaseUrl =
  typeof window !== "undefined" && window.__API_BASE_URL__
    ? window.__API_BASE_URL__
    : "/api/v1";

const buildRequestId = (): string => {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }
  return `req_${Date.now()}_${Math.random().toString(16).slice(2)}`;
};

export const httpClient = axios.create({
  baseURL: defaultBaseUrl,
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

httpClient.interceptors.request.use((config) => {
  config.headers = config.headers ?? {};

  if (!config.headers["X-Request-Id"]) {
    config.headers["X-Request-Id"] = buildRequestId();
  }

  return config;
});

httpClient.interceptors.response.use(
  (response) => response,
  (error) => Promise.reject(toApiClientError(error)),
);

export const setApiBaseUrl = (baseUrl: string): void => {
  httpClient.defaults.baseURL = baseUrl.replace(/\/+$/, "");
};

export const getApiBaseUrl = (): string => {
  return httpClient.defaults.baseURL ?? "/api/v1";
};

export const openApiConfiguration = new Configuration({
  accessToken: () => tokenStore.getAccessToken() ?? "",
});
