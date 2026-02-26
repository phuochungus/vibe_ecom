# OrderDetail


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
**items** | [**Array&lt;OrderItem&gt;**](OrderItem.md) |  | [optional] [default to undefined]
**payments** | [**Array&lt;PaymentTransaction&gt;**](PaymentTransaction.md) |  | [optional] [default to undefined]

## Example

```typescript
import { OrderDetail } from 'golf-store-api-client';

const instance: OrderDetail = {
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
    items,
    payments,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
