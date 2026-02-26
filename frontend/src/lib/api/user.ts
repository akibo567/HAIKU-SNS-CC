import { apiRequest } from "./client";
import type { User } from "@/types/user";

export async function fetchUser(username: string) {
  return apiRequest<{ data: User }>(`/users/${username}`);
}

export async function updateProfile(
  params: { displayName: string; bio?: string },
  token: string
) {
  return apiRequest<{ data: User }>("/users/me", {
    method: "PUT",
    body: JSON.stringify(params),
    token,
  });
}
