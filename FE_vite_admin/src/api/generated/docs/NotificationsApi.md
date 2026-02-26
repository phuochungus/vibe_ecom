# NotificationsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**notificationsGet**](#notificationsget) | **GET** /notifications | List notifications|
|[**notificationsNotificationIdReadPatch**](#notificationsnotificationidreadpatch) | **PATCH** /notifications/{notification_id}/read | Mark notification as read|
|[**notificationsReadAllPatch**](#notificationsreadallpatch) | **PATCH** /notifications/read-all | Mark all notifications as read|

# **notificationsGet**
> NotificationListEnvelope notificationsGet()


### Example

```typescript
import {
    NotificationsApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new NotificationsApi(configuration);

let status: NotificationStatus; // (optional) (default to undefined)
let page: number; // (optional) (default to 1)
let pageSize: number; // (optional) (default to 20)

const { status, data } = await apiInstance.notificationsGet(
    status,
    page,
    pageSize
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **status** | **NotificationStatus** |  | (optional) defaults to undefined|
| **page** | [**number**] |  | (optional) defaults to 1|
| **pageSize** | [**number**] |  | (optional) defaults to 20|


### Return type

**NotificationListEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Notification list |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **notificationsNotificationIdReadPatch**
> NotificationEnvelope notificationsNotificationIdReadPatch()


### Example

```typescript
import {
    NotificationsApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new NotificationsApi(configuration);

let notificationId: string; // (default to undefined)

const { status, data } = await apiInstance.notificationsNotificationIdReadPatch(
    notificationId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **notificationId** | [**string**] |  | defaults to undefined|


### Return type

**NotificationEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Updated |  -  |
|**401** | Unauthorized |  -  |
|**404** | Resource not found |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **notificationsReadAllPatch**
> ReadAllNotificationsEnvelope notificationsReadAllPatch()


### Example

```typescript
import {
    NotificationsApi,
    Configuration
} from 'golf-store-api-client';

const configuration = new Configuration();
const apiInstance = new NotificationsApi(configuration);

const { status, data } = await apiInstance.notificationsReadAllPatch();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ReadAllNotificationsEnvelope**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Updated |  -  |
|**401** | Unauthorized |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

