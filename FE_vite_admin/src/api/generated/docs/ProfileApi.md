# ProfileApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**meGet**](#meget) | **GET** /me | Current user profile|

# **meGet**
> UserProfileEnvelope meGet()


### Example

```typescript
import {
    ProfileApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new ProfileApi(configuration);

const { status, data } = await apiInstance.meGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**UserProfileEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Profile |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

