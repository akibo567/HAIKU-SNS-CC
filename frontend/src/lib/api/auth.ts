import { apiRequest } from "./client";
import type { User } from "@/types/user";

export type AuthResponse = {
  data: {
    accessToken: string;
    user: User;
  };
};

export type RefreshResponse = {
  data: { accessToken: string };
};

export async function register(params: {
  username: string;
  email: string;
  password: string;
  displayName: string;
}): Promise<AuthResponse> {
  return apiRequest("/auth/register", {
    method: "POST",
    body: JSON.stringify(params),
  });
}

export async function login(params: {
  email: string;
  password: string;
}): Promise<AuthResponse> {
  return apiRequest("/auth/login", {
    method: "POST",
    body: JSON.stringify(params),
  });
}

export async function logout(token: string): Promise<void> {
  await apiRequest("/auth/logout", { method: "POST", token });
}

export async function refreshToken(): Promise<RefreshResponse> {
  return apiRequest("/auth/refresh", { method: "POST" });
}
