## golf-store-api-client@1.0.0

This generator creates TypeScript/JavaScript client that utilizes [axios](https://github.com/axios/axios). The generated Node module can be used in the following environments:

Environment
* Node.js
* Webpack
* Browserify

Language level
* ES5 - you must have a Promises/A+ library installed
* ES6

Module system
* CommonJS
* ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition will be automatically resolved via `package.json`. ([Reference](https://www.typescriptlang.org/docs/handbook/declaration-files/consumption.html))

### Building

To build and compile the typescript sources to javascript use:
```
npm install
npm run build
```

### Publishing

First build the package then run `npm publish`

### Consuming

navigate to the folder of your consuming project and run one of the following commands.

_published:_

```
npm install golf-store-api-client@1.0.0 --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
```

### Documentation for API Endpoints

All URIs are relative to */api/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*AdminOrdersApi* | [**adminOrdersGet**](docs/AdminOrdersApi.md#adminordersget) | **GET** /admin/orders | List orders (admin)
*AdminOrdersApi* | [**adminOrdersOrderIdGet**](docs/AdminOrdersApi.md#adminordersorderidget) | **GET** /admin/orders/{order_id} | Order detail (admin)
*AdminOrdersApi* | [**adminOrdersOrderIdStatusPatch**](docs/AdminOrdersApi.md#adminordersorderidstatuspatch) | **PATCH** /admin/orders/{order_id}/status | Update order status (admin)
*AdminProductsApi* | [**adminProductsPost**](docs/AdminProductsApi.md#adminproductspost) | **POST** /admin/products | Create product (admin)
*AdminProductsApi* | [**adminProductsProductIdDelete**](docs/AdminProductsApi.md#adminproductsproductiddelete) | **DELETE** /admin/products/{product_id} | Delete product (soft delete)
*AdminProductsApi* | [**adminProductsProductIdPatch**](docs/AdminProductsApi.md#adminproductsproductidpatch) | **PATCH** /admin/products/{product_id} | Update product (admin)
*AdminRevenueApi* | [**adminRevenueOrdersGet**](docs/AdminRevenueApi.md#adminrevenueordersget) | **GET** /admin/revenue/orders | Revenue orders breakdown
*AdminRevenueApi* | [**adminRevenueSummaryGet**](docs/AdminRevenueApi.md#adminrevenuesummaryget) | **GET** /admin/revenue/summary | Revenue summary
*AuthApi* | [**authLoginPost**](docs/AuthApi.md#authloginpost) | **POST** /auth/login | Login
*AuthApi* | [**authLogoutPost**](docs/AuthApi.md#authlogoutpost) | **POST** /auth/logout | Logout
*AuthApi* | [**authRefreshPost**](docs/AuthApi.md#authrefreshpost) | **POST** /auth/refresh | Refresh token
*NotificationsApi* | [**notificationsGet**](docs/NotificationsApi.md#notificationsget) | **GET** /notifications | List notifications
*NotificationsApi* | [**notificationsNotificationIdReadPatch**](docs/NotificationsApi.md#notificationsnotificationidreadpatch) | **PATCH** /notifications/{notification_id}/read | Mark notification as read
*NotificationsApi* | [**notificationsReadAllPatch**](docs/NotificationsApi.md#notificationsreadallpatch) | **PATCH** /notifications/read-all | Mark all notifications as read
*OrdersApi* | [**ordersGet**](docs/OrdersApi.md#ordersget) | **GET** /orders | List my orders
*OrdersApi* | [**ordersOrderIdCancelPost**](docs/OrdersApi.md#ordersorderidcancelpost) | **POST** /orders/{order_id}/cancel | Cancel order
*OrdersApi* | [**ordersOrderIdGet**](docs/OrdersApi.md#ordersorderidget) | **GET** /orders/{order_id} | Order detail
*OrdersApi* | [**ordersOrderIdTrackingGet**](docs/OrdersApi.md#ordersorderidtrackingget) | **GET** /orders/{order_id}/tracking | Order tracking timeline
*OrdersApi* | [**ordersPost**](docs/OrdersApi.md#orderspost) | **POST** /orders | Create order
*PaymentsApi* | [**ordersOrderIdPaymentsGet**](docs/PaymentsApi.md#ordersorderidpaymentsget) | **GET** /orders/{order_id}/payments | List payments by order
*PaymentsApi* | [**ordersOrderIdPaymentsPost**](docs/PaymentsApi.md#ordersorderidpaymentspost) | **POST** /orders/{order_id}/payments | Create payment for order
*ProductsApi* | [**productsGet**](docs/ProductsApi.md#productsget) | **GET** /products | List products
*ProductsApi* | [**productsProductIdGet**](docs/ProductsApi.md#productsproductidget) | **GET** /products/{product_id} | Product detail
*ProfileApi* | [**meGet**](docs/ProfileApi.md#meget) | **GET** /me | Current user profile
*WebhooksApi* | [**webhooksPaymentsProviderPost**](docs/WebhooksApi.md#webhookspaymentsproviderpost) | **POST** /webhooks/payments/{provider} | Payment webhook entrypoint


### Documentation For Models

 - [AcceptedEnvelope](docs/AcceptedEnvelope.md)
 - [AcceptedEnvelopeAllOfData](docs/AcceptedEnvelopeAllOfData.md)
 - [AdminCreateProductRequest](docs/AdminCreateProductRequest.md)
 - [AdminUpdateOrderStatusRequest](docs/AdminUpdateOrderStatusRequest.md)
 - [AdminUpdateProductRequest](docs/AdminUpdateProductRequest.md)
 - [AuthTokenData](docs/AuthTokenData.md)
 - [CancelOrderRequest](docs/CancelOrderRequest.md)
 - [CreateOrderRequest](docs/CreateOrderRequest.md)
 - [CreatePaymentEnvelope](docs/CreatePaymentEnvelope.md)
 - [CreatePaymentRequest](docs/CreatePaymentRequest.md)
 - [CreatePaymentResult](docs/CreatePaymentResult.md)
 - [EnvelopeBase](docs/EnvelopeBase.md)
 - [ErrorEnvelope](docs/ErrorEnvelope.md)
 - [ErrorObject](docs/ErrorObject.md)
 - [LoginRequest](docs/LoginRequest.md)
 - [LoginResponseEnvelope](docs/LoginResponseEnvelope.md)
 - [Notification](docs/Notification.md)
 - [NotificationEnvelope](docs/NotificationEnvelope.md)
 - [NotificationListEnvelope](docs/NotificationListEnvelope.md)
 - [NotificationListEnvelopeAllOfData](docs/NotificationListEnvelopeAllOfData.md)
 - [NotificationStatus](docs/NotificationStatus.md)
 - [OrderDetail](docs/OrderDetail.md)
 - [OrderDetailEnvelope](docs/OrderDetailEnvelope.md)
 - [OrderItem](docs/OrderItem.md)
 - [OrderItemInput](docs/OrderItemInput.md)
 - [OrderListEnvelope](docs/OrderListEnvelope.md)
 - [OrderListEnvelopeAllOfData](docs/OrderListEnvelopeAllOfData.md)
 - [OrderStatus](docs/OrderStatus.md)
 - [OrderSummary](docs/OrderSummary.md)
 - [OrderSummaryEnvelope](docs/OrderSummaryEnvelope.md)
 - [OrderTracking](docs/OrderTracking.md)
 - [OrderTrackingEnvelope](docs/OrderTrackingEnvelope.md)
 - [Pagination](docs/Pagination.md)
 - [PaymentListEnvelope](docs/PaymentListEnvelope.md)
 - [PaymentListEnvelopeAllOfData](docs/PaymentListEnvelopeAllOfData.md)
 - [PaymentStatus](docs/PaymentStatus.md)
 - [PaymentTransaction](docs/PaymentTransaction.md)
 - [PaymentTxnState](docs/PaymentTxnState.md)
 - [Product](docs/Product.md)
 - [ProductDetailEnvelope](docs/ProductDetailEnvelope.md)
 - [ProductListEnvelope](docs/ProductListEnvelope.md)
 - [ProductListEnvelopeAllOfData](docs/ProductListEnvelopeAllOfData.md)
 - [ProductStatus](docs/ProductStatus.md)
 - [ReadAllNotificationsEnvelope](docs/ReadAllNotificationsEnvelope.md)
 - [ReadAllNotificationsEnvelopeAllOfData](docs/ReadAllNotificationsEnvelopeAllOfData.md)
 - [RefreshTokenRequest](docs/RefreshTokenRequest.md)
 - [RefreshTokenResponseEnvelope](docs/RefreshTokenResponseEnvelope.md)
 - [RefreshTokenResponseEnvelopeAllOfData](docs/RefreshTokenResponseEnvelopeAllOfData.md)
 - [RevenueOrder](docs/RevenueOrder.md)
 - [RevenueOrdersEnvelope](docs/RevenueOrdersEnvelope.md)
 - [RevenueOrdersEnvelopeAllOfData](docs/RevenueOrdersEnvelopeAllOfData.md)
 - [RevenueSummary](docs/RevenueSummary.md)
 - [RevenueSummaryEnvelope](docs/RevenueSummaryEnvelope.md)
 - [ShippingAddressInput](docs/ShippingAddressInput.md)
 - [TrackingEvent](docs/TrackingEvent.md)
 - [UserProfile](docs/UserProfile.md)
 - [UserProfileEnvelope](docs/UserProfileEnvelope.md)
 - [UserRole](docs/UserRole.md)
 - [UserStatus](docs/UserStatus.md)


<a id="documentation-for-authorization"></a>
## Documentation For Authorization


Authentication schemes defined for the API:
<a id="bearerAuth"></a>
### bearerAuth

- **Type**: Bearer authentication (JWT)

