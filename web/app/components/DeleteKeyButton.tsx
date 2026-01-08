"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

type Props = {
  dbId: number;
  keyName: string;
};

export default function DeleteKeyButton({ dbId, keyName }: Props) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleDelete() {
    const ok = confirm(`Are you sure you want to delete key:\n\n${keyName}`);

    if (!ok) return;

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(
        `/api/db/${dbId}/key/${encodeURIComponent(keyName)}`,
        {
          method: "DELETE",
        }
      );

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text);
      }

      // redirect back to db page
      router.push(`/db/${dbId}`);
      router.refresh();
    } catch (err: any) {
      setError(err.message || "Failed to delete key");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="space-y-2">
      <button
        onClick={handleDelete}
        disabled={loading}
        className="
          px-4 py-2 rounded text-sm font-medium
          bg-red-600 hover:bg-red-700
          disabled:opacity-50
          transition
        "
      >
        {loading ? "Deleting..." : "Delete Key"}
      </button>

      {error && <p className="text-sm text-red-400">{error}</p>}
    </div>
  );
}
