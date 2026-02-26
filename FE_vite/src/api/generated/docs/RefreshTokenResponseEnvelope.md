# RefreshTokenResponseEnvelope


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**success** | **boolean** |  | [default to undefined]
**request_id** | **string** |  | [default to undefined]
**timestamp** | **string** |  | [default to undefined]
**meta** | **{ [key: string]: any; }** |  | [optional] [default to undefined]
**data** | [**RefreshTokenResponseEnvelopeAllOfData**](RefreshTokenResponseEnvelopeAllOfData.md) |  | [default to undefined]

## Example

```typescript
import { RefreshTokenResponseEnvelope } from 'golf-store-api-client';

const instance: RefreshTokenResponseEnvelope = {
    success,
    request_id,
    timestamp,
    meta,
    data,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
