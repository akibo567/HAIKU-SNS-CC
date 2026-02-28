// SSR（サーバーサイド）では Docker 内部 URL で直接バックエンドへ、
// ブラウザからは Nginx 経由でアクセス
const API_URL =
  typeof window === "undefined"
    ? "http://backend:8080/api"
    : (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost/api");

type RequestOptions = RequestInit & { token?: string };

export async function apiRequest<T>(
  path: string,
  options: RequestOptions = {}
): Promise<T> {
  const { token, ...init } = options;

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init.headers as Record<string, string>),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_URL}/v1${path}`, {
    ...init,
    headers,
    credentials: "include",
  });

  const data = await res.json();

  if (!res.ok) {
    const err = data?.error;
    throw new ApiError(
      err?.code ?? "UNKNOWN",
      err?.message ?? "エラーが発生しました",
      res.status
    );
  }

  return data;
}

export class ApiError extends Error {
  constructor(
    public code: string,
    message: string,
    public status: number
  ) {
    super(message);
    this.name = "ApiError";
  }
}
