# OrderTrackingEnvelope


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**success** | **boolean** |  | [default to undefined]
**request_id** | **string** |  | [default to undefined]
**timestamp** | **string** |  | [default to undefined]
**meta** | **{ [key: string]: any; }** |  | [optional] [default to undefined]
**data** | [**OrderTracking**](OrderTracking.md) |  | [optional] [default to undefined]

## Example

```typescript
import { OrderTrackingEnvelope } from 'golf-store-api-client';

const instance: OrderTrackingEnvelope = {
    success,
    request_id,
    timestamp,
    meta,
    data,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
