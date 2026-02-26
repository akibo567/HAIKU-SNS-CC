"use client";

import { useAutoRefresh } from "@/hooks/useAuth";

export function AuthProvider({ children }: { children: React.ReactNode }) {
  useAutoRefresh();
  return <>{children}</>;
}
