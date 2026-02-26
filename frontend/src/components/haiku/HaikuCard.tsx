"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/hooks/useAuth";
import { likeHaiku, unlikeHaiku, deleteHaiku } from "@/lib/api/haiku";
import type { HaikuPost } from "@/types/haiku";

type Props = {
  post: HaikuPost;
  onDelete?: (id: string) => void;
};

export function HaikuCard({ post, onDelete }: Props) {
  const { accessToken, user } = useAuth();
  const [liked, setLiked] = useState(post.likedByMe);
  const [likeCount, setLikeCount] = useState(post.likeCount);
  const [loading, setLoading] = useState(false);

  const isOwner = user?.id === post.author.id;

  async function handleLike() {
    if (!accessToken || loading) return;
    setLoading(true);
    try {
      if (liked) {
        await unlikeHaiku(post.id, accessToken);
        setLiked(false);
        setLikeCount((c) => Math.max(c - 1, 0));
      } else {
        await likeHaiku(post.id, accessToken);
        setLiked(true);
        setLikeCount((c) => c + 1);
      }
    } catch {
      // 楽観的更新を元に戻す
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete() {
    if (!accessToken || !confirm("この俳句を削除しますか？")) return;
    try {
      await deleteHaiku(post.id, accessToken);
      onDelete?.(post.id);
    } catch {
      alert("削除に失敗しました");
    }
  }

  const dateStr = new Date(post.createdAt).toLocaleDateString("ja-JP", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });

  return (
    <article className="bg-white rounded-xl border border-gray-200 p-5 shadow-sm hover:shadow-md transition-shadow">
      <Link href={`/post/${post.id}`}>
        <div className="mb-4 space-y-1 cursor-pointer">
          <p className="text-xl font-serif text-gray-900 leading-relaxed">
            {post.ku1}
          </p>
          <p className="text-xl font-serif text-gray-900 leading-relaxed pl-4">
            {post.ku2}
          </p>
          <p className="text-xl font-serif text-gray-900 leading-relaxed">
            {post.ku3}
          </p>
        </div>
      </Link>

      <div className="flex items-center justify-between text-sm text-gray-500">
        <div className="flex items-center gap-2 flex-wrap">
          <Link
            href={`/profile/${post.author.username}`}
            className="font-medium text-gray-700 hover:text-indigo-600 transition-colors"
          >
            {post.author.displayName}
          </Link>
          <span>@{post.author.username}</span>
          <span>·</span>
          <span>{dateStr}</span>
        </div>

        <div className="flex items-center gap-3">
          {isOwner && (
            <button
              onClick={handleDelete}
              className="text-gray-400 hover:text-red-500 transition-colors text-xs"
            >
              削除
            </button>
          )}
          <button
            onClick={handleLike}
            disabled={!accessToken || loading}
            className={`flex items-center gap-1 transition-colors ${
              liked ? "text-pink-500" : "text-gray-400 hover:text-pink-500"
            } disabled:cursor-default`}
          >
            <span>{liked ? "♥" : "♡"}</span>
            <span>{likeCount}</span>
          </button>
        </div>
      </div>
    </article>
  );
}
