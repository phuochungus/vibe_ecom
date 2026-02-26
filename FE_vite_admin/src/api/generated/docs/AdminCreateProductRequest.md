# AdminCreateProductRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**sku** | **string** |  | [default to undefined]
**name** | **string** |  | [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**price** | **string** |  | [default to undefined]
**stock** | **number** |  | [default to undefined]
**status** | [**ProductStatus**](ProductStatus.md) |  | [default to undefined]
**image_url** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { AdminCreateProductRequest } from 'golf-store-api-client';

const instance: AdminCreateProductRequest = {
    sku,
    name,
    description,
    price,
    stock,
    status,
    image_url,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
