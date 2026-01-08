export default function HomePage() {
  return (
    <div className="h-full flex items-center justify-center">
      <div className="text-center max-w-md">
        <h2 className="text-2xl font-bold mb-2">Welcome to FerroDB</h2>
        <p className="text-zinc-400 mb-4">
          Select a database from the sidebar to view or manage keys.
        </p>

        <div className="border border-zinc-800 rounded p-4 text-sm text-zinc-500">
          No database selected
        </div>
      </div>
    </div>
  );
}
