"use client";

import { useState } from "react";
import Link from "next/link";
import useSWR from "swr";
import { useAuth } from "@/hooks/useAuth";
import { useHaikuValidation } from "@/hooks/useHaikuValidation";
import { fetchReplies, createReply } from "@/lib/api/haiku";
import { HaikuCounter } from "./HaikuCounter";
import type { KuKey } from "@/lib/mora";
import type { Reply } from "@/types/haiku";

type Props = {
  postId: string;
};

export function ReplySection({ postId }: Props) {
  const { accessToken } = useAuth();
  const { data, mutate } = useSWR<{ data: Reply[] }>(
    `replies-${postId}`,
    () => fetchReplies(postId)
  );

  const { ku1, ku2, ku3, setKu1, setKu2, setKu3, counts, isValid } =
    useHaikuValidation();
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const replies = data?.data ?? [];

  const fields: {
    key: KuKey;
    label: string;
    value: string;
    setter: (v: string) => void;
    placeholder: string;
  }[] = [
    { key: "ku1", label: "上の句（5音）", value: ku1, setter: setKu1, placeholder: "はるのかぜ" },
    { key: "ku2", label: "中の句（7音）", value: ku2, setter: setKu2, placeholder: "こころおだやかに" },
    { key: "ku3", label: "下の句（5音）", value: ku3, setter: setKu3, placeholder: "こたえけり" },
  ];

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!accessToken || !isValid) return;

    setSubmitting(true);
    setError(null);
    try {
      const res = await createReply(postId, { ku1, ku2, ku3 }, accessToken);
      setKu1("");
      setKu2("");
      setKu3("");
      mutate(
        (prev) => ({ data: [...(prev?.data ?? []), res.data] }),
        { revalidate: false }
      );
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : "投稿に失敗しました");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="mt-8 space-y-6">
      <h2 className="text-base font-semibold text-gray-600 border-b pb-2">
        返句 {replies.length > 0 && <span className="text-gray-400 font-normal">（{replies.length}句）</span>}
      </h2>

      {/* 返句一覧 */}
      {replies.length > 0 ? (
        <ul className="space-y-3">
          {replies.map((reply) => (
            <li key={reply.id} className="bg-white rounded-xl border border-gray-200 p-4 shadow-sm">
              <div className="mb-3 space-y-1">
                <p className="text-lg font-serif text-gray-900 leading-relaxed">{reply.ku1}</p>
                <p className="text-lg font-serif text-gray-900 leading-relaxed pl-4">{reply.ku2}</p>
                <p className="text-lg font-serif text-gray-900 leading-relaxed">{reply.ku3}</p>
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <Link
                  href={`/profile/${reply.author.username}`}
                  className="font-medium text-gray-700 hover:text-indigo-600 transition-colors"
                >
                  {reply.author.displayName}
                </Link>
                <span>@{reply.author.username}</span>
                <span>·</span>
                <span>
                  {new Date(reply.createdAt).toLocaleDateString("ja-JP", {
                    year: "numeric",
                    month: "short",
                    day: "numeric",
                  })}
                </span>
              </div>
            </li>
          ))}
        </ul>
      ) : (
        <p className="text-sm text-gray-400 text-center py-4">まだ返句はありません</p>
      )}

      {/* 返句フォーム（ログイン時のみ表示） */}
      {accessToken && (
        <form onSubmit={handleSubmit} className="space-y-3 bg-gray-50 rounded-xl p-4 border border-gray-200">
          <p className="text-sm text-gray-500">この俳句に返句する</p>
          {fields.map(({ key, label, value, setter, placeholder }) => (
            <div key={key}>
              <div className="flex items-center justify-between mb-1">
                <label className="text-xs font-medium text-gray-600">{label}</label>
                <HaikuCounter kuKey={key} count={counts[key]} />
              </div>
              <input
                type="text"
                value={value}
                onChange={(e) => setter(e.target.value)}
                placeholder={placeholder}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-lg font-serif focus:outline-none focus:ring-2 focus:ring-indigo-500 bg-white"
              />
            </div>
          ))}
          {error && (
            <p className="text-sm text-red-600 bg-red-50 p-2 rounded-md">{error}</p>
          )}
          <button
            type="submit"
            disabled={!isValid || submitting}
            className="w-full py-2 px-4 bg-indigo-600 text-white font-semibold rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm"
          >
            {submitting ? "投稿中..." : "返句を詠む"}
          </button>
          {!isValid && (ku1 || ku2 || ku3) && (
            <p className="text-xs text-center text-gray-500">5・7・5音になると投稿できます</p>
          )}
        </form>
      )}
    </div>
  );
}
