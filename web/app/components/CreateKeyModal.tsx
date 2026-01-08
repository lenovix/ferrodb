"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

type Props = {
  dbId: number;
};

export default function CreateKeyModal({ dbId }: Props) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubmit() {
    if (!key.trim()) {
      setError("Key is required");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const res = await fetch(
        `/api/db/${dbId}/key/${encodeURIComponent(key)}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ value }),
        }
      );

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text);
      }

      setOpen(false);
      setKey("");
      setValue("");

      router.refresh();
    } catch (err: any) {
      setError(err.message || "Failed to create key");
    } finally {
      setLoading(false);
    }
  }

  if (!open) {
    return (
      <button
        onClick={() => setOpen(true)}
        className="
          px-3 py-2 rounded text-sm font-medium
          bg-blue-600 hover:bg-blue-700
          transition
        "
      >
        + New Key
      </button>
    );
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60">
      <div className="w-full max-w-lg bg-zinc-900 border border-zinc-800 rounded-lg">
        {/* Header */}
        <div className="px-4 py-3 border-b border-zinc-800 flex justify-between">
          <h2 className="font-semibold">Create New Key</h2>
          <button
            onClick={() => setOpen(false)}
            className="text-zinc-400 hover:text-zinc-200"
          >
            ✕
          </button>
        </div>

        {/* Body */}
        <div className="p-4 space-y-4">
          <div>
            <label className="block text-sm text-zinc-400 mb-1">Key</label>
            <input
              value={key}
              onChange={(e) => setKey(e.target.value)}
              className="
                w-full px-3 py-2 rounded
                bg-zinc-950 border border-zinc-800
                text-sm
              "
              placeholder="my:key:name"
            />
          </div>

          <div>
            <label className="block text-sm text-zinc-400 mb-1">Value</label>
            <textarea
              value={value}
              onChange={(e) => setValue(e.target.value)}
              rows={5}
              className="
                w-full px-3 py-2 rounded
                bg-zinc-950 border border-zinc-800
                text-sm font-mono
              "
              placeholder="value..."
            />
          </div>

          {error && <p className="text-sm text-red-400">{error}</p>}
        </div>

        {/* Footer */}
        <div className="px-4 py-3 border-t border-zinc-800 flex justify-end gap-2">
          <button
            onClick={() => setOpen(false)}
            className="px-3 py-2 text-sm text-zinc-400 hover:text-zinc-200"
          >
            Cancel
          </button>
          <button
            onClick={handleSubmit}
            disabled={loading}
            className="
              px-4 py-2 rounded text-sm font-medium
              bg-blue-600 hover:bg-blue-700
              disabled:opacity-50
            "
          >
            {loading ? "Saving..." : "Create"}
          </button>
        </div>
      </div>
    </div>
  );
}
