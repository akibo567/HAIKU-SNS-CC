"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/hooks/useAuth";
import { logout } from "@/lib/api/auth";

export function Header() {
  const { user, accessToken, logout: storeLogout } = useAuth();
  const router = useRouter();

  async function handleLogout() {
    if (accessToken) {
      try {
        await logout(accessToken);
      } catch {
        // ignore
      }
    }
    storeLogout();
    router.push("/");
  }

  return (
    <header className="sticky top-0 z-50 bg-white border-b border-gray-200">
      <div className="max-w-2xl mx-auto px-4 h-14 flex items-center justify-between">
        <Link href="/" className="text-xl font-bold text-indigo-700 tracking-tight">
          Go Shichi Go!
        </Link>

        <nav className="flex items-center gap-4">
          {user ? (
            <>
              <Link
                href="/post/new"
                className="px-4 py-1.5 bg-indigo-600 text-white text-sm font-medium rounded-full hover:bg-indigo-700 transition-colors"
              >
                詠む
              </Link>
              <Link
                href={`/profile/${user.username}`}
                className="text-sm text-gray-700 hover:text-indigo-600 transition-colors"
              >
                {user.displayName}
              </Link>
              <button
                onClick={handleLogout}
                className="text-sm text-gray-500 hover:text-gray-700 transition-colors"
              >
                ログアウト
              </button>
            </>
          ) : (
            <>
              <Link
                href="/login"
                className="text-sm text-gray-700 hover:text-indigo-600 transition-colors"
              >
                ログイン
              </Link>
              <Link
                href="/register"
                className="px-4 py-1.5 bg-indigo-600 text-white text-sm font-medium rounded-full hover:bg-indigo-700 transition-colors"
              >
                はじめる
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
