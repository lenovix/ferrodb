import Breadcrumb from "@/app/components/Breadcrumb";
import CreateKeyModal from "@/app/components/CreateKeyModal";

type Props = {
  params: Promise<{ id: string }>;
};

export default async function DBPage({ params }: Props) {
  const { id } = await params;
  const dbId = Number(id);

  if (Number.isNaN(dbId)) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-red-500 font-bold mb-2">Invalid Database ID</h1>
          <p className="text-zinc-400 text-sm">Database ID must be a number</p>
        </div>
      </div>
    );
  }

  const res = await fetch(`http://localhost:8080/api/db/${dbId}/keys`, {
    cache: "no-store",
  });

  if (!res.ok) {
    const text = await res.text();
    return <pre className="text-red-400 bg-zinc-900 p-4 rounded">{text}</pre>;
  }

  const data = await res.json();

  return (
    <section className="space-y-6">
      <Breadcrumb
        items={[{ label: "Home", href: "/" }, { label: `DB ${dbId}` }]}
      />
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold">
            Database <span className="text-blue-500">{dbId}</span>
          </h1>
          <p className="text-sm text-zinc-400">Keys stored in this database</p>
        </div>

        <CreateKeyModal dbId={dbId} />
      </div>

      {/* Keys */}
      <div className="border border-zinc-800 rounded-lg overflow-hidden">
        {data.keys.length === 0 ? (
          <div className="p-6 text-center text-zinc-500 text-sm">
            No keys found in this database
          </div>
        ) : (
          <ul className="divide-y divide-zinc-800">
            {data.keys.map((key: string) => (
              <li key={key}>
                <a
                  href={`/db/${dbId}/key/${encodeURIComponent(key)}`}
                  className="
                    block px-4 py-3 text-sm
                    text-zinc-200
                    hover:bg-zinc-800
                    transition
                  "
                >
                  {key}
                </a>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  );
}
