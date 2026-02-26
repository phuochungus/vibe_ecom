# RevenueOrder


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**order_code** | **string** |  | [default to undefined]
**order_status** | [**OrderStatus**](OrderStatus.md) |  | [default to undefined]
**payment_status** | [**PaymentStatus**](PaymentStatus.md) |  | [default to undefined]
**subtotal_amount** | **string** |  | [default to undefined]
**discount_amount** | **string** |  | [default to undefined]
**shipping_fee** | **string** |  | [default to undefined]
**total_amount** | **string** |  | [default to undefined]
**payment_due_at** | **string** |  | [optional] [default to undefined]
**placed_at** | **string** |  | [default to undefined]

## Example

```typescript
import { RevenueOrder } from 'golf-store-api-client';

const instance: RevenueOrder = {
    id,
    order_code,
    order_status,
    payment_status,
    subtotal_amount,
    discount_amount,
    shipping_fee,
    total_amount,
    payment_due_at,
    placed_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
