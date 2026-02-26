# Notification


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**channel** | **string** |  | [optional] [default to undefined]
**event_type** | **string** |  | [default to undefined]
**event_key** | **string** |  | [optional] [default to undefined]
**title** | **string** |  | [default to undefined]
**content** | **string** |  | [default to undefined]
**status** | [**NotificationStatus**](NotificationStatus.md) |  | [default to undefined]
**read** | **boolean** |  | [optional] [default to undefined]
**sent_at** | **string** |  | [optional] [default to undefined]
**created_at** | **string** |  | [default to undefined]

## Example

```typescript
import { Notification } from 'golf-store-api-client';

const instance: Notification = {
    id,
    channel,
    event_type,
    event_key,
    title,
    content,
    status,
    read,
    sent_at,
    created_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
