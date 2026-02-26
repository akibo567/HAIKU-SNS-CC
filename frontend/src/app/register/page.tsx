"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { register } from "@/lib/api/auth";
import { useAuth } from "@/hooks/useAuth";

export default function RegisterPage() {
  const router = useRouter();
  const { login: storeLogin } = useAuth();
  const [form, setForm] = useState({
    username: "",
    email: "",
    password: "",
    displayName: "",
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      const res = await register(form);
      storeLogin(res.data.accessToken, res.data.user);
      router.push("/");
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "登録に失敗しました");
    } finally {
      setLoading(false);
    }
  }

  const fields = [
    { name: "displayName" as const, label: "表示名", type: "text", placeholder: "松尾芭蕉" },
    { name: "username" as const, label: "ユーザー名", type: "text", placeholder: "basho" },
    { name: "email" as const, label: "メールアドレス", type: "email", placeholder: "basho@example.com" },
    { name: "password" as const, label: "パスワード", type: "password", placeholder: "8文字以上" },
  ];

  return (
    <div className="max-w-sm mx-auto">
      <h1 className="text-2xl font-bold text-center mb-6">新規登録</h1>

      <form onSubmit={handleSubmit} className="space-y-4">
        {fields.map(({ name, label, type, placeholder }) => (
          <div key={name}>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {label}
            </label>
            <input
              type={type}
              name={name}
              value={form[name]}
              onChange={handleChange}
              placeholder={placeholder}
              required
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
            />
          </div>
        ))}

        {error && (
          <p className="text-sm text-red-600 bg-red-50 p-3 rounded-md">{error}</p>
        )}

        <button
          type="submit"
          disabled={loading}
          className="w-full py-3 bg-indigo-600 text-white font-semibold rounded-md hover:bg-indigo-700 disabled:opacity-50 transition-colors"
        >
          {loading ? "登録中..." : "はじめる"}
        </button>
      </form>

      <p className="mt-4 text-center text-sm text-gray-500">
        既にアカウントをお持ちの方は{" "}
        <Link href="/login" className="text-indigo-600 hover:underline">
          ログイン
        </Link>
      </p>
    </div>
  );
}
