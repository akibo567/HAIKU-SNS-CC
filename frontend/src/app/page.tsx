"use client";

import { useCallback } from "react";
import useSWRInfinite from "swr/infinite";
import { HaikuCard } from "@/components/haiku/HaikuCard";
import { useAuth } from "@/hooks/useAuth";
import { fetchTimeline } from "@/lib/api/haiku";
import type { HaikuPost } from "@/types/haiku";

type PageData = {
  data: HaikuPost[];
  meta: { cursor: string; hasNext: boolean };
};

export default function TimelinePage() {
  const { accessToken } = useAuth();

  const getKey = (pageIndex: number, prev: PageData | null) => {
    if (prev && !prev.meta.hasNext) return null;
    const cursor = prev?.meta.cursor;
    return ["timeline", cursor, accessToken] as const;
  };

  const { data, size, setSize, isLoading, mutate } = useSWRInfinite(
    getKey,
    ([, cursor]: readonly [string, string | undefined, string | null]) =>
      fetchTimeline(cursor ?? undefined, accessToken ?? undefined),
    { revalidateFirstPage: true }
  );

  const posts = data?.flatMap((d) => d.data) ?? [];
  const hasMore = data?.[data.length - 1]?.meta.hasNext ?? false;

  const handleDelete = useCallback(() => {
    mutate();
  }, [mutate]);

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold text-gray-900">タイムライン</h1>

      {isLoading && posts.length === 0 && (
        <div className="text-center py-12 text-gray-400">読み込み中...</div>
      )}

      {!isLoading && posts.length === 0 && (
        <div className="text-center py-12 text-gray-400">
          まだ俳句がありません。最初の一句を詠みましょう！
        </div>
      )}

      {posts.map((post) => (
        <HaikuCard key={post.id} post={post} onDelete={handleDelete} />
      ))}

      {hasMore && (
        <button
          onClick={() => setSize(size + 1)}
          className="w-full py-3 text-sm text-gray-500 hover:text-gray-700 transition-colors"
        >
          もっと見る
        </button>
      )}
    </div>
  );
}
