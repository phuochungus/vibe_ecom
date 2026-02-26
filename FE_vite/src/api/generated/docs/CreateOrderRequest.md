# CreateOrderRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**items** | [**Array&lt;OrderItemInput&gt;**](OrderItemInput.md) |  | [default to undefined]
**shipping_address** | [**ShippingAddressInput**](ShippingAddressInput.md) |  | [default to undefined]
**customer_note** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { CreateOrderRequest } from 'golf-store-api-client';

const instance: CreateOrderRequest = {
    items,
    shipping_address,
    customer_note,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
