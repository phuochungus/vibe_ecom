# AdminUpdateProductRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **string** |  | [optional] [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**price** | **string** |  | [optional] [default to undefined]
**stock** | **number** |  | [optional] [default to undefined]
**status** | [**ProductStatus**](ProductStatus.md) |  | [optional] [default to undefined]
**image_url** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { AdminUpdateProductRequest } from 'golf-store-api-client';

const instance: AdminUpdateProductRequest = {
    name,
    description,
    price,
    stock,
    status,
    image_url,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
