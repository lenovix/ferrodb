"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

type Props = {
  dbId: number;
  keyName: string;
  value: string;
};

export default function EditKeyValue({ dbId, keyName, value }: Props) {
  const router = useRouter();
  const [editing, setEditing] = useState(false);
  const [newValue, setNewValue] = useState(value);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSave() {
    setLoading(true);
    setError(null);

    try {
      const res = await fetch(
        `/api/db/${dbId}/key/${encodeURIComponent(keyName)}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ value: newValue }),
        }
      );

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text);
      }

      setEditing(false);
      router.refresh();
    } catch (err: any) {
      setError(err.message || "Failed to update value");
    } finally {
      setLoading(false);
    }
  }

  if (!editing) {
    return (
      <div className="space-y-2">
        <pre className="bg-zinc-950 border border-zinc-800 p-4 rounded text-sm overflow-auto">
          {value}
        </pre>

        <button
          onClick={() => setEditing(true)}
          className="
            px-3 py-1.5 rounded text-sm
            bg-zinc-800 hover:bg-zinc-700
          "
        >
          ✏️ Edit Value
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <textarea
        value={newValue}
        onChange={(e) => setNewValue(e.target.value)}
        rows={8}
        className="
          w-full bg-zinc-950 border border-zinc-800
          rounded p-3 text-sm font-mono
        "
      />

      {error && <p className="text-sm text-red-400">{error}</p>}

      <div className="flex gap-2">
        <button
          onClick={() => {
            setNewValue(value);
            setEditing(false);
          }}
          className="px-3 py-1.5 text-sm text-zinc-400 hover:text-zinc-200"
        >
          Cancel
        </button>

        <button
          onClick={handleSave}
          disabled={loading}
          className="
            px-4 py-1.5 rounded text-sm font-medium
            bg-blue-600 hover:bg-blue-700
            disabled:opacity-50
          "
        >
          {loading ? "Saving..." : "Save"}
        </button>
      </div>
    </div>
  );
}
