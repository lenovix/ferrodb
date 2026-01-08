import Link from "next/link";

type Props = {
  dbCount: number;
};

export default function Sidebar({ dbCount }: Props) {
  return (
    <aside className="w-64 border-r border-zinc-800 bg-zinc-950 p-4">
      {/* Logo */}
      <div className="mb-6">
        <h1 className="text-lg font-bold text-blue-500">FerroDB</h1>
        <p className="text-xs text-zinc-500">Admin Console</p>
      </div>

      {/* DB List */}
      <nav className="space-y-1">
        {Array.from({ length: dbCount }).map((_, i) => (
          <Link
            key={i}
            href={`/db/${i}`}
            className="block rounded px-3 py-2 text-sm
              text-zinc-300 hover:bg-zinc-800 hover:text-white"
          >
            DB {i}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
