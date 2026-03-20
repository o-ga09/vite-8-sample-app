const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:18080";

/**
 * Orval 8 のカスタムインスタンス。
 * 生成コードは customInstance<ResponseType>(url, fetchInit) の形で呼び出す。
 * レスポンスは { data, status, headers } 形式で返す。
 */
export const customInstance = async <T>(url: string, options?: RequestInit): Promise<T> => {
  const response = await fetch(`${BASE_URL}${url}`, options);
  const headers = response.headers;
  const status = response.status;

  if (!response.ok) {
    const errorBody = (await response.json().catch(() => ({}))) as Record<string, unknown>;
    throw new Error(typeof errorBody["message"] === "string" ? errorBody["message"] : `HTTP ${response.status}`);
  }

  if (status === 204) return { data: undefined, status, headers } as T;

  const data = (await response.json()) as unknown;
  return { data, status, headers } as T;
};
