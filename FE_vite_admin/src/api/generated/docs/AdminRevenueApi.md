# AdminRevenueApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminRevenueOrdersGet**](#adminrevenueordersget) | **GET** /admin/revenue/orders | Revenue orders breakdown|
|[**adminRevenueSummaryGet**](#adminrevenuesummaryget) | **GET** /admin/revenue/summary | Revenue summary|

# **adminRevenueOrdersGet**
> RevenueOrdersEnvelope adminRevenueOrdersGet()


### Example

```typescript
import {
    AdminRevenueApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminRevenueApi(configuration);

let from: string; // (optional) (default to undefined)
let to: string; // (optional) (default to undefined)
let page: number; // (optional) (default to 1)
let pageSize: number; // (optional) (default to 20)

const { status, data } = await apiInstance.adminRevenueOrdersGet(
    from,
    to,
    page,
    pageSize
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **from** | [**string**] |  | (optional) defaults to undefined|
| **to** | [**string**] |  | (optional) defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to 1|
| **pageSize** | [**number**] |  | (optional) defaults to 20|


### Return type

**RevenueOrdersEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Revenue orders list |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminRevenueSummaryGet**
> RevenueSummaryEnvelope adminRevenueSummaryGet()


### Example

```typescript
import {
    AdminRevenueApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminRevenueApi(configuration);

let from: string; // (optional) (default to undefined)
let to: string; // (optional) (default to undefined)
let groupBy: 'day' | 'month' | 'year'; // (optional) (default to undefined)

const { status, data } = await apiInstance.adminRevenueSummaryGet(
    from,
    to,
    groupBy
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **from** | [**string**] |  | (optional) defaults to undefined|
| **to** | [**string**] |  | (optional) defaults to undefined|
| **groupBy** | [**&#39;day&#39; | &#39;month&#39; | &#39;year&#39;**]**Array<&#39;day&#39; &#124; &#39;month&#39; &#124; &#39;year&#39;>** |  | (optional) defaults to undefined|


### Return type

**RevenueSummaryEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Revenue summary |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

