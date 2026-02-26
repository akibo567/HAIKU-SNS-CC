import { create } from "zustand";
import type { User } from "@/types/user";

type AuthState = {
  user: User | null;
  accessToken: string | null;
  login: (token: string, user: User) => void;
  logout: () => void;
  setToken: (token: string) => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  accessToken: null,
  login: (token, user) => set({ accessToken: token, user }),
  logout: () => set({ accessToken: null, user: null }),
  setToken: (token) => set({ accessToken: token }),
}));
