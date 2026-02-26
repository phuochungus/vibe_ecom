# PaymentTransaction


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **string** |  | [default to undefined]
**txn_type** | **string** |  | [default to undefined]
**provider** | **string** |  | [default to undefined]
**provider_txn_code** | **string** |  | [optional] [default to undefined]
**status** | [**PaymentTxnState**](PaymentTxnState.md) |  | [default to undefined]
**amount** | **string** |  | [default to undefined]
**currency_code** | **string** |  | [default to undefined]
**created_at** | **string** |  | [default to undefined]

## Example

```typescript
import { PaymentTransaction } from 'golf-store-api-client';

const instance: PaymentTransaction = {
    id,
    txn_type,
    provider,
    provider_txn_code,
    status,
    amount,
    currency_code,
    created_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
