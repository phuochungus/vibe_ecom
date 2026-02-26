# AdminProductsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**adminProductsPost**](#adminproductspost) | **POST** /admin/products | Create product (admin)|
|[**adminProductsProductIdDelete**](#adminproductsproductiddelete) | **DELETE** /admin/products/{product_id} | Delete product (soft delete)|
|[**adminProductsProductIdPatch**](#adminproductsproductidpatch) | **PATCH** /admin/products/{product_id} | Update product (admin)|

# **adminProductsPost**
> ProductDetailEnvelope adminProductsPost(adminCreateProductRequest)


### Example

```typescript
import {
    AdminProductsApi,
    Configuration,
    AdminCreateProductRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminProductsApi(configuration);

let adminCreateProductRequest: AdminCreateProductRequest; //

const { status, data } = await apiInstance.adminProductsPost(
    adminCreateProductRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **adminCreateProductRequest** | **AdminCreateProductRequest**|  | |


### Return type

**ProductDetailEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Product created |  -  |
|**400** | Validation error |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminProductsProductIdDelete**
> adminProductsProductIdDelete()


### Example

```typescript
import {
    AdminProductsApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminProductsApi(configuration);

let productId: string; // (default to undefined)

const { status, data } = await apiInstance.adminProductsProductIdDelete(
    productId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productId** | [**string**] |  | defaults to undefined|


### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**204** | Product deleted |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **adminProductsProductIdPatch**
> ProductDetailEnvelope adminProductsProductIdPatch(adminUpdateProductRequest)


### Example

```typescript
import {
    AdminProductsApi,
    Configuration,
    AdminUpdateProductRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AdminProductsApi(configuration);

let productId: string; // (default to undefined)
let adminUpdateProductRequest: AdminUpdateProductRequest; //

const { status, data } = await apiInstance.adminProductsProductIdPatch(
    productId,
    adminUpdateProductRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **adminUpdateProductRequest** | **AdminUpdateProductRequest**|  | |
| **productId** | [**string**] |  | defaults to undefined|


### Return type

**ProductDetailEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product updated |  -  |
|**400** | Validation error |  -  |
|**401** | Unauthorized |  -  |
|**403** | Forbidden |  -  |
|**404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

