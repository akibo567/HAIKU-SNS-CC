import type { User } from "./user";

export type HaikuPost = {
  id: string;
  ku1: string;
  ku2: string;
  ku3: string;
  likeCount: number;
  likedByMe: boolean;
  createdAt: string;
  author: Pick<User, "id" | "username" | "displayName">;
};

export type CreateHaikuInput = {
  ku1: string;
  ku2: string;
  ku3: string;
};

export type PaginatedResponse<T> = {
  data: T[];
  meta: {
    cursor: string;
    hasNext: boolean;
  };
};
