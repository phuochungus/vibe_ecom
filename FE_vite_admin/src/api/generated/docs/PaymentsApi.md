# PaymentsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**ordersOrderIdPaymentsGet**](#ordersorderidpaymentsget) | **GET** /orders/{order_id}/payments | List payments by order|
|[**ordersOrderIdPaymentsPost**](#ordersorderidpaymentspost) | **POST** /orders/{order_id}/payments | Create payment for order|

# **ordersOrderIdPaymentsGet**
> PaymentListEnvelope ordersOrderIdPaymentsGet()


### Example

```typescript
import {
    PaymentsApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new PaymentsApi(configuration);

let orderId: string; // (default to undefined)

const { status, data } = await apiInstance.ordersOrderIdPaymentsGet(
    orderId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **orderId** | [**string**] |  | defaults to undefined|


### Return type

**PaymentListEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Payment list |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersOrderIdPaymentsPost**
> CreatePaymentEnvelope ordersOrderIdPaymentsPost(createPaymentRequest)


### Example

```typescript
import {
    PaymentsApi,
    Configuration,
    CreatePaymentRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new PaymentsApi(configuration);

let orderId: string; // (default to undefined)
let idempotencyKey: string; // (default to undefined)
let createPaymentRequest: CreatePaymentRequest; //

const { status, data } = await apiInstance.ordersOrderIdPaymentsPost(
    orderId,
    idempotencyKey,
    createPaymentRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **createPaymentRequest** | **CreatePaymentRequest**|  | |
| **orderId** | [**string**] |  | defaults to undefined|
| **idempotencyKey** | [**string**] |  | defaults to undefined|


### Return type

**CreatePaymentEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Payment created |  -  |
|**401** | Unauthorized |  -  |
|**404** | Resource not found |  -  |
|**409** | Duplicate payment request |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

