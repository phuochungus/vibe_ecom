# OrdersApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**ordersGet**](#ordersget) | **GET** /orders | List my orders|
|[**ordersOrderIdCancelPost**](#ordersorderidcancelpost) | **POST** /orders/{order_id}/cancel | Cancel order|
|[**ordersOrderIdGet**](#ordersorderidget) | **GET** /orders/{order_id} | Order detail|
|[**ordersOrderIdTrackingGet**](#ordersorderidtrackingget) | **GET** /orders/{order_id}/tracking | Order tracking timeline|
|[**ordersPost**](#orderspost) | **POST** /orders | Create order|

# **ordersGet**
> OrderListEnvelope ordersGet()


### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let status: OrderStatus; // (optional) (default to undefined)
let from: string; // (optional) (default to undefined)
let to: string; // (optional) (default to undefined)
let page: number; // (optional) (default to 1)
let pageSize: number; // (optional) (default to 20)

const { status, data } = await apiInstance.ordersGet(
    status,
    from,
    to,
    page,
    pageSize
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **status** | **OrderStatus** |  | (optional) defaults to undefined|
| **from** | [**string**] |  | (optional) defaults to undefined|
| **to** | [**string**] |  | (optional) defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to 1|
| **pageSize** | [**number**] |  | (optional) defaults to 20|


### Return type

**OrderListEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Order list |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersOrderIdCancelPost**
> OrderSummaryEnvelope ordersOrderIdCancelPost(cancelOrderRequest)


### Example

```typescript
import {
    OrdersApi,
    Configuration,
    CancelOrderRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let orderId: string; // (default to undefined)
let cancelOrderRequest: CancelOrderRequest; //

const { status, data } = await apiInstance.ordersOrderIdCancelPost(
    orderId,
    cancelOrderRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **cancelOrderRequest** | **CancelOrderRequest**|  | |
| **orderId** | [**string**] |  | defaults to undefined|


### Return type

**OrderSummaryEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Order canceled |  -  |
|**401** | Unauthorized |  -  |
|**404** | Resource not found |  -  |
|**409** | Order cannot be canceled |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersOrderIdGet**
> OrderDetailEnvelope ordersOrderIdGet()


### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let orderId: string; // (default to undefined)

const { status, data } = await apiInstance.ordersOrderIdGet(
    orderId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orderId** | [**string**] |  | defaults to undefined|


### Return type

**OrderDetailEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Order detail |  -  |
|**401** | Unauthorized |  -  |
|**404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersOrderIdTrackingGet**
> OrderTrackingEnvelope ordersOrderIdTrackingGet()


### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let orderId: string; // (default to undefined)

const { status, data } = await apiInstance.ordersOrderIdTrackingGet(
    orderId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orderId** | [**string**] |  | defaults to undefined|


### Return type

**OrderTrackingEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Tracking data |  -  |
|**401** | Unauthorized |  -  |
|**404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersPost**
> OrderSummaryEnvelope ordersPost(createOrderRequest)


### Example

```typescript
import {
    OrdersApi,
    Configuration,
    CreateOrderRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let idempotencyKey: string; // (default to undefined)
let createOrderRequest: CreateOrderRequest; //

const { status, data } = await apiInstance.ordersPost(
    idempotencyKey,
    createOrderRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createOrderRequest** | **CreateOrderRequest**|  | |
| **idempotencyKey** | [**string**] |  | defaults to undefined|


### Return type

**OrderSummaryEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Order created |  -  |
|**400** | Validation error |  -  |
|**401** | Unauthorized |  -  |
|**409** | Product out of stock |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

