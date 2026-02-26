import { useState } from "react";
import { countMora, HAIKU_MORA, type KuKey } from "@/lib/mora";

export function useHaikuValidation() {
  const [ku1, setKu1] = useState("");
  const [ku2, setKu2] = useState("");
  const [ku3, setKu3] = useState("");

  const counts: Record<KuKey, number> = {
    ku1: countMora(ku1),
    ku2: countMora(ku2),
    ku3: countMora(ku3),
  };

  const isValid =
    counts.ku1 === HAIKU_MORA.ku1 &&
    counts.ku2 === HAIKU_MORA.ku2 &&
    counts.ku3 === HAIKU_MORA.ku3;

  return { ku1, ku2, ku3, setKu1, setKu2, setKu3, counts, isValid };
}
