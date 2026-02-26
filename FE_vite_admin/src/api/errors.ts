import axios, { AxiosError } from "axios";

type ErrorEnvelope = {
  error?: {
    code?: string;
    message?: string;
    details?: unknown;
  };
  request_id?: string;
  timestamp?: string;
};

export class ApiClientError extends Error {
  readonly status: number;
  readonly code: string;
  readonly requestId?: string;
  readonly details?: unknown;
  readonly timestamp?: string;
  readonly original?: unknown;

  constructor(params: {
    message: string;
    status: number;
    code: string;
    requestId?: string;
    details?: unknown;
    timestamp?: string;
    original?: unknown;
  }) {
    super(params.message);
    this.name = "ApiClientError";
    this.status = params.status;
    this.code = params.code;
    this.requestId = params.requestId;
    this.details = params.details;
    this.timestamp = params.timestamp;
    this.original = params.original;
  }
}

export const toApiClientError = (error: unknown): ApiClientError => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<ErrorEnvelope>;
    const payload = axiosError.response?.data;
    const status = axiosError.response?.status ?? 0;
    const message = payload?.error?.message ?? axiosError.message ?? "Unknown API error";
    const code = payload?.error?.code ?? "INTERNAL_ERROR";

    return new ApiClientError({
      message,
      status,
      code,
      requestId: payload?.request_id,
      details: payload?.error?.details,
      timestamp: payload?.timestamp,
      original: error,
    });
  }

  if (error instanceof Error) {
    return new ApiClientError({
      message: error.message,
      status: 0,
      code: "UNKNOWN_ERROR",
      original: error,
    });
  }

  return new ApiClientError({
    message: "Unknown API error",
    status: 0,
    code: "UNKNOWN_ERROR",
    original: error,
  });
};
