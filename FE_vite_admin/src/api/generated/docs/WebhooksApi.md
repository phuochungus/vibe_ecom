# WebhooksApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**webhooksPaymentsProviderPost**](#webhookspaymentsproviderpost) | **POST** /webhooks/payments/{provider} | Payment webhook entrypoint|

# **webhooksPaymentsProviderPost**
> AcceptedEnvelope webhooksPaymentsProviderPost(requestBody)

External provider calls this endpoint via API Gateway.

### Example

```typescript
import {
    WebhooksApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new WebhooksApi(configuration);

let provider: string; // (default to undefined)
let requestBody: { [key: string]: any; }; //

const { status, data } = await apiInstance.webhooksPaymentsProviderPost(
    provider,
    requestBody
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **requestBody** | **{ [key: string]: any; }**|  | |
| **provider** | [**string**] |  | defaults to undefined|


### Return type

**AcceptedEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**202** | Accepted for async processing |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

