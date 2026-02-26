"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { useHaikuValidation } from "@/hooks/useHaikuValidation";
import { useAuth } from "@/hooks/useAuth";
import { createHaiku } from "@/lib/api/haiku";
import { HaikuCounter } from "./HaikuCounter";
import type { KuKey } from "@/lib/mora";

export function HaikuForm() {
  const router = useRouter();
  const { accessToken } = useAuth();
  const { ku1, ku2, ku3, setKu1, setKu2, setKu3, counts, isValid } =
    useHaikuValidation();
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!accessToken || !isValid) return;

    setSubmitting(true);
    setError(null);

    try {
      await createHaiku({ ku1, ku2, ku3 }, accessToken);
      router.push("/");
      router.refresh();
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "投稿に失敗しました");
    } finally {
      setSubmitting(false);
    }
  }

  const fields: {
    key: KuKey;
    label: string;
    value: string;
    setter: (v: string) => void;
    placeholder: string;
  }[] = [
    { key: "ku1", label: "上の句（5音）", value: ku1, setter: setKu1, placeholder: "はるのはな" },
    { key: "ku2", label: "中の句（7音）", value: ku2, setter: setKu2, placeholder: "やまにかすみが" },
    { key: "ku3", label: "下の句（5音）", value: ku3, setter: setKu3, placeholder: "たなびいて" },
  ];

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {fields.map(({ key, label, value, setter, placeholder }) => (
        <div key={key}>
          <div className="flex items-center justify-between mb-1">
            <label className="text-sm font-medium text-gray-700">{label}</label>
            <HaikuCounter kuKey={key} count={counts[key]} />
          </div>
          <input
            type="text"
            value={value}
            onChange={(e) => setter(e.target.value)}
            placeholder={placeholder}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-lg font-serif focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
      ))}

      {error && (
        <p className="text-sm text-red-600 bg-red-50 p-3 rounded-md">{error}</p>
      )}

      <button
        type="submit"
        disabled={!isValid || submitting}
        className="w-full py-3 px-4 bg-indigo-600 text-white font-semibold rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        {submitting ? "投稿中..." : "俳句を詠む"}
      </button>

      {!isValid && (ku1 || ku2 || ku3) && (
        <p className="text-xs text-center text-gray-500">
          5・7・5音になると投稿できます
        </p>
      )}
    </form>
  );
}
