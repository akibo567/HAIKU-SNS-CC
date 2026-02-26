"use client";

import { use } from "react";
import useSWR from "swr";
import useSWRInfinite from "swr/infinite";
import { fetchUser } from "@/lib/api/user";
import { fetchUserPosts } from "@/lib/api/haiku";
import { HaikuCard } from "@/components/haiku/HaikuCard";
import { useAuth } from "@/hooks/useAuth";
import type { HaikuPost } from "@/types/haiku";

type PageData = {
  data: HaikuPost[];
  meta: { cursor: string; hasNext: boolean };
};

type Props = {
  params: Promise<{ username: string }>;
};

export default function ProfilePage({ params }: Props) {
  const { username } = use(params);
  const { accessToken } = useAuth();

  const { data: userData, error: userError } = useSWR(
    ["user", username],
    () => fetchUser(username)
  );

  const getKey = (pageIndex: number, prev: PageData | null) => {
    if (prev && !prev.meta.hasNext) return null;
    return ["userPosts", username, prev?.meta.cursor, accessToken] as const;
  };

  const { data: postsData, size, setSize } = useSWRInfinite(
    getKey,
    ([, , cursor]: readonly [string, string, string | undefined, string | null]) =>
      fetchUserPosts(username, cursor, accessToken ?? undefined)
  );

  if (userError) {
    return (
      <div className="text-center py-12 text-gray-400">
        ユーザーが見つかりません
      </div>
    );
  }

  if (!userData) {
    return <div className="text-center py-12 text-gray-400">読み込み中...</div>;
  }

  const user = userData.data;
  const posts = postsData?.flatMap((d) => d.data) ?? [];
  const hasMore = postsData?.[postsData.length - 1]?.meta.hasNext ?? false;

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-xl border border-gray-200 p-6 shadow-sm">
        <h1 className="text-xl font-bold text-gray-900">{user.displayName}</h1>
        <p className="text-sm text-gray-500">@{user.username}</p>
        {user.bio && (
          <p className="mt-2 text-sm text-gray-700">{user.bio}</p>
        )}
      </div>

      <div className="space-y-4">
        <h2 className="font-semibold text-gray-700">投稿した俳句</h2>

        {posts.length === 0 && (
          <p className="text-center py-8 text-gray-400">
            まだ俳句を詠んでいません
          </p>
        )}

        {posts.map((post) => (
          <HaikuCard key={post.id} post={post} />
        ))}

        {hasMore && (
          <button
            onClick={() => setSize(size + 1)}
            className="w-full py-3 text-sm text-gray-500 hover:text-gray-700"
          >
            もっと見る
          </button>
        )}
      </div>
    </div>
  );
}
