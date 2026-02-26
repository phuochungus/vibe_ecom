# AdminOrdersApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminOrdersGet**](#adminordersget) | **GET** /admin/orders | List orders (admin)|
|[**adminOrdersOrderIdGet**](#adminordersorderidget) | **GET** /admin/orders/{order_id} | Order detail (admin)|
|[**adminOrdersOrderIdStatusPatch**](#adminordersorderidstatuspatch) | **PATCH** /admin/orders/{order_id}/status | Update order status (admin)|

# **adminOrdersGet**
> OrderListEnvelope adminOrdersGet()


### Example

```typescript
import {
    AdminOrdersApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminOrdersApi(configuration);

let status: OrderStatus; // (optional) (default to undefined)
let from: string; // (optional) (default to undefined)
let to: string; // (optional) (default to undefined)
let page: number; // (optional) (default to 1)
let pageSize: number; // (optional) (default to 20)

const { status, data } = await apiInstance.adminOrdersGet(
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
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOrdersOrderIdGet**
> OrderDetailEnvelope adminOrdersOrderIdGet()


### Example

```typescript
import {
    AdminOrdersApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminOrdersApi(configuration);

let orderId: string; // (default to undefined)

const { status, data } = await apiInstance.adminOrdersOrderIdGet(
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
|**403** | Forbidden |  -  |
|**404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminOrdersOrderIdStatusPatch**
> OrderSummaryEnvelope adminOrdersOrderIdStatusPatch(adminUpdateOrderStatusRequest)


### Example

```typescript
import {
    AdminOrdersApi,
    Configuration,
    AdminUpdateOrderStatusRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminOrdersApi(configuration);

let orderId: string; // (default to undefined)
let adminUpdateOrderStatusRequest: AdminUpdateOrderStatusRequest; //

const { status, data } = await apiInstance.adminOrdersOrderIdStatusPatch(
    orderId,
    adminUpdateOrderStatusRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **adminUpdateOrderStatusRequest** | **AdminUpdateOrderStatusRequest**|  | |
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
|**200** | Status updated |  -  |
|**400** | Validation error |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Resource not found |  -  |
|**409** | Conflict |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

