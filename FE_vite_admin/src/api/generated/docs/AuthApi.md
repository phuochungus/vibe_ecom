# AuthApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**authLoginPost**](#authloginpost) | **POST** /auth/login | Login|
|[**authLogoutPost**](#authlogoutpost) | **POST** /auth/logout | Logout|
|[**authRefreshPost**](#authrefreshpost) | **POST** /auth/refresh | Refresh token|

# **authLoginPost**
> LoginResponseEnvelope authLoginPost(loginRequest)


### Example

```typescript
import {
    AuthApi,
    Configuration,
    LoginRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let loginRequest: LoginRequest; //

const { status, data } = await apiInstance.authLoginPost(
    loginRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **loginRequest** | **LoginRequest**|  | |


### Return type

**LoginResponseEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Login success |  -  |
|**401** | Unauthorized |  -  |
|**423** | Account locked |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authLogoutPost**
> authLogoutPost()


### Example

```typescript
import {
    AuthApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.authLogoutPost();
```

### Parameters
This endpoint does not have any parameters.


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
|**204** | Logout success |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authRefreshPost**
> RefreshTokenResponseEnvelope authRefreshPost(refreshTokenRequest)


### Example

```typescript
import {
    AuthApi,
    Configuration,
    RefreshTokenRequest
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let refreshTokenRequest: RefreshTokenRequest; //

const { status, data } = await apiInstance.authRefreshPost(
    refreshTokenRequest
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **refreshTokenRequest** | **RefreshTokenRequest**|  | |


### Return type

**RefreshTokenResponseEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Refresh success |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

