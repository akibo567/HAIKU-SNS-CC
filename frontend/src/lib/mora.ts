/**
 * フロントエンド補助用モーラカウンター。
 * バックエンドのカウントが正式。これはUX補助のみ。
 */
export function countMora(text: string): number {
  let count = 0;
  for (const char of text) {
    if (char === " " || char === "　" || char === "\t" || char === "\n") {
      continue;
    }
    count++;
  }
  return count;
}

export const HAIKU_MORA = { ku1: 5, ku2: 7, ku3: 5 } as const;
export type KuKey = keyof typeof HAIKU_MORA;
