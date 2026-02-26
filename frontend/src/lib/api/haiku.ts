import { apiRequest } from "./client";
import type { HaikuPost, CreateHaikuInput } from "@/types/haiku";

type PostsResponse = {
  data: HaikuPost[];
  meta: { cursor: string; hasNext: boolean };
};

export async function fetchTimeline(cursor?: string, token?: string): Promise<PostsResponse> {
  const q = cursor ? `?cursor=${cursor}` : "";
  return apiRequest<PostsResponse>(`/posts${q}`, { token });
}

export async function fetchPost(id: string, token?: string) {
  return apiRequest<{ data: HaikuPost }>(`/posts/${id}`, { token });
}

export async function createHaiku(input: CreateHaikuInput, token: string) {
  return apiRequest<{ data: HaikuPost }>("/posts", {
    method: "POST",
    body: JSON.stringify(input),
    token,
  });
}

export async function deleteHaiku(id: string, token: string) {
  return apiRequest(`/posts/${id}`, { method: "DELETE", token });
}

export async function likeHaiku(id: string, token: string) {
  return apiRequest(`/posts/${id}/like`, { method: "POST", token });
}

export async function unlikeHaiku(id: string, token: string) {
  return apiRequest(`/posts/${id}/like`, { method: "DELETE", token });
}

export async function fetchUserPosts(
  username: string,
  cursor?: string,
  token?: string
): Promise<PostsResponse> {
  const q = cursor ? `?cursor=${cursor}` : "";
  return apiRequest<PostsResponse>(`/users/${username}/posts${q}`, { token });
}
