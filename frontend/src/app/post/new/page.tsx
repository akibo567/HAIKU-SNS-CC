"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { HaikuForm } from "@/components/haiku/HaikuForm";
import { useAuth } from "@/hooks/useAuth";

export default function NewPostPage() {
  const { user, accessToken } = useAuth();
  const router = useRouter();

  useEffect(() => {
    // accessToken が確定して null の場合（未ログイン）はログインページへ
    if (accessToken === null && user === null) {
      router.push("/login");
    }
  }, [user, accessToken, router]);

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-2xl font-bold mb-6">俳句を詠む</h1>
      <div className="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
        <HaikuForm />
      </div>
      <p className="mt-4 text-xs text-center text-gray-400">
        上の句（5音）・中の句（7音）・下の句（5音）で一句詠みましょう
      </p>
    </div>
  );
}
