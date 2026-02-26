# ErrorEnvelope


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**success** | **boolean** |  | [default to undefined]
**request_id** | **string** |  | [default to undefined]
**timestamp** | **string** |  | [default to undefined]
**meta** | **{ [key: string]: any; }** |  | [optional] [default to undefined]
**error** | [**ErrorObject**](ErrorObject.md) |  | [default to undefined]

## Example

```typescript
import { ErrorEnvelope } from 'golf-store-api-client';

const instance: ErrorEnvelope = {
    success,
    request_id,
    timestamp,
    meta,
    error,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
