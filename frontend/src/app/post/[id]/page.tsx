import { fetchPost } from "@/lib/api/haiku";
import { HaikuCard } from "@/components/haiku/HaikuCard";
import { ReplySection } from "@/components/haiku/ReplySection";
import { notFound } from "next/navigation";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function PostDetailPage({ params }: Props) {
  const { id } = await params;

  let post;
  try {
    const res = await fetchPost(id);
    post = res.data;
  } catch {
    notFound();
  }

  return (
    <div className="max-w-md mx-auto">
      <h1 className="text-lg font-semibold text-gray-500 mb-4">俳句の詳細</h1>
      <HaikuCard post={post} />
      <ReplySection postId={id} />
    </div>
  );
}
