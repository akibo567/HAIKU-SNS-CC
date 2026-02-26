import { HAIKU_MORA, type KuKey } from "@/lib/mora";

type Props = {
  kuKey: KuKey;
  count: number;
};

export function HaikuCounter({ kuKey, count }: Props) {
  const target = HAIKU_MORA[kuKey];

  const color =
    count === target
      ? "text-green-600"
      : count === 0
        ? "text-gray-400"
        : "text-red-500";

  return (
    <span className={`text-xs font-mono tabular-nums ${color}`}>
      {count}/{target}音
    </span>
  );
}
