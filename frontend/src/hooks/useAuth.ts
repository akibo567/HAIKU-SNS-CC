"use client";

import { useEffect } from "react";
import { useAuthStore } from "@/store/authStore";
import { refreshToken } from "@/lib/api/auth";

export function useAuth() {
  return useAuthStore();
}

export function useAutoRefresh() {
  const { setToken } = useAuthStore();

  useEffect(() => {
    refreshToken()
      .then((res) => setToken(res.data.accessToken))
      .catch(() => {
        // リフレッシュ失敗は正常（未ログイン）
      });
  }, [setToken]);
}
