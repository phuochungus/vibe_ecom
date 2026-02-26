# OrderTracking


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**order_id** | **string** |  | [default to undefined]
**current_status** | [**OrderStatus**](OrderStatus.md) |  | [default to undefined]
**timeline** | [**Array&lt;TrackingEvent&gt;**](TrackingEvent.md) |  | [default to undefined]

## Example

```typescript
import { OrderTracking } from 'golf-store-api-client';

const instance: OrderTracking = {
    order_id,
    current_status,
    timeline,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
