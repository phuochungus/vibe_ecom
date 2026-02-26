export const getEnvelopeData = <T>(payload: unknown): T | undefined => {
  if (!payload || typeof payload !== "object") {
    return undefined;
  }
  return (payload as { data?: T }).data;
};

export const getEnvelopeItems = <T>(payload: unknown): T[] => {
  const data = getEnvelopeData<{ items?: T[] }>(payload);
  return data?.items ?? [];
};

export const getEnvelopePagination = (payload: unknown) => {
  const data = getEnvelopeData<{ pagination?: { page?: number; page_size?: number; total_items?: number; total_pages?: number } }>(payload);
  return data?.pagination;
};
