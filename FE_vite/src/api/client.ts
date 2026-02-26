import {
  AdminOrdersApi,
  AdminProductsApi,
  AdminRevenueApi,
  AuthApi,
  NotificationsApi,
  OrdersApi,
  PaymentsApi,
  ProductsApi,
  ProfileApi,
  WebhooksApi,
} from "./generated/api";
import { openApiConfiguration, httpClient } from "./runtime";
import { tokenStore } from "./token-store";

export const apiClient = {
  auth: new AuthApi(openApiConfiguration, undefined, httpClient),
  profile: new ProfileApi(openApiConfiguration, undefined, httpClient),
  products: new ProductsApi(openApiConfiguration, undefined, httpClient),
  orders: new OrdersApi(openApiConfiguration, undefined, httpClient),
  payments: new PaymentsApi(openApiConfiguration, undefined, httpClient),
  notifications: new NotificationsApi(openApiConfiguration, undefined, httpClient),
  adminProducts: new AdminProductsApi(openApiConfiguration, undefined, httpClient),
  adminOrders: new AdminOrdersApi(openApiConfiguration, undefined, httpClient),
  adminRevenue: new AdminRevenueApi(openApiConfiguration, undefined, httpClient),
  webhooks: new WebhooksApi(openApiConfiguration, undefined, httpClient),
};

export const apiSession = {
  getAccessToken: tokenStore.getAccessToken,
  setAccessToken: tokenStore.setAccessToken,
  getRefreshToken: tokenStore.getRefreshToken,
  setRefreshToken: tokenStore.setRefreshToken,
  clear: tokenStore.clear,
};
