# FE API Layer

## Mục tiêu
- Typed API client generate từ `docs/FE/openapi.yaml`.
- Một entrypoint dùng chung cho FE qua `api-gateway`.
- Quản lý token + chuẩn hóa lỗi API ở một nơi.

## Cấu trúc
- `generated/`: code auto-generate từ OpenAPI (không sửa tay).
- `runtime.ts`: axios instance, `X-Request-Id`, base URL.
- `client.ts`: gom các service instances (`apiClient.*`).
- `token-store.ts`: lưu/đọc access token, refresh token.
- `errors.ts`: chuẩn hóa lỗi HTTP thành `ApiClientError`.

## Generate lại client
```powershell
cd FE
powershell -NoProfile -ExecutionPolicy Bypass -File .\scripts\generate-api-client.ps1
```

## Cách dùng nhanh
```ts
import { apiClient, apiSession } from "./api";

apiSession.setAccessToken("<jwt>");

const res = await apiClient.products.productsGet();
console.log(res.data.data.items);
```

## Base URL
- Mặc định: `/api/v1`.
- Có thể set runtime:
```ts
import { setApiBaseUrl } from "./api";
setApiBaseUrl("https://api.example.com/api/v1");
```
