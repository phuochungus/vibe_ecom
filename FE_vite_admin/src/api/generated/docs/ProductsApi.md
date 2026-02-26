# ProductsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**productsGet**](#productsget) | **GET** /products | List products|
|[**productsProductIdGet**](#productsproductidget) | **GET** /products/{product_id} | Product detail|

# **productsGet**
> ProductListEnvelope productsGet()


### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let q: string; // (optional) (default to undefined)
let status: ProductStatus; // (optional) (default to undefined)
let minPrice: string; // (optional) (default to undefined)
let maxPrice: string; // (optional) (default to undefined)
let page: number; // (optional) (default to 1)
let pageSize: number; // (optional) (default to 20)
let sort: 'created_at' | 'price' | 'name'; // (optional) (default to undefined)
let order: 'asc' | 'desc'; // (optional) (default to undefined)

const { status, data } = await apiInstance.productsGet(
    q,
    status,
    minPrice,
    maxPrice,
    page,
    pageSize,
    sort,
    order
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **q** | [**string**] |  | (optional) defaults to undefined|
| **status** | **ProductStatus** |  | (optional) defaults to undefined|
| **minPrice** | [**string**] |  | (optional) defaults to undefined|
| **maxPrice** | [**string**] |  | (optional) defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to 1|
| **pageSize** | [**number**] |  | (optional) defaults to 20|
| **sort** | [**&#39;created_at&#39; | &#39;price&#39; | &#39;name&#39;**]**Array<&#39;created_at&#39; &#124; &#39;price&#39; &#124; &#39;name&#39;>** |  | (optional) defaults to undefined|
| **order** | [**&#39;asc&#39; | &#39;desc&#39;**]**Array<&#39;asc&#39; &#124; &#39;desc&#39;>** |  | (optional) defaults to undefined|


### Return type

**ProductListEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product list |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsProductIdGet**
> ProductDetailEnvelope productsProductIdGet()


### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let productId: string; // (default to undefined)

const { status, data } = await apiInstance.productsProductIdGet(
    productId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productId** | [**string**] |  | defaults to undefined|


### Return type

**ProductDetailEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product detail |  -  |
|**404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

