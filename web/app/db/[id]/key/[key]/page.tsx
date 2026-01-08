import Breadcrumb from "@/app/components/Breadcrumb";
import DeleteKeyButton from "@/app/components/DeleteKeyButton";
import EditKeyValue from "@/app/components/EditKeyValue";

type Props = {
  params: Promise<{
    id: string;
    key: string;
  }>;
};

export default async function KeyPage({ params }: Props) {
  const { id, key } = await params;
  const dbId = Number(id);

  if (Number.isNaN(dbId)) {
    return (
      <div className="h-full flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-red-500 font-bold mb-2">Invalid Database</h1>
          <p className="text-zinc-400 text-sm">Database ID must be a number</p>
        </div>
      </div>
    );
  }

  const res = await fetch(
    `http://localhost:8080/api/db/${dbId}/key/${encodeURIComponent(key)}`,
    { cache: "no-store" }
  );

  if (!res.ok) {
    const text = await res.text();
    return (
      <pre className="bg-zinc-900 border border-zinc-800 text-red-400 p-4 rounded">
        {text}
      </pre>
    );
  }

  const data = await res.json();

  return (
    <section className="space-y-6">
      {/* Breadcrumb */}
      <Breadcrumb
        items={[
          { label: "Home", href: "/" },
          { label: `DB ${dbId}`, href: `/db/${dbId}` },
          { label: data.key },
        ]}
      />

      {/* Header */}
      <div>
        <h1 className="text-xl font-semibold mt-2">
          Key <span className="text-blue-500 break-all">{data.key}</span>
        </h1>
        <p className="text-sm text-zinc-400">
          Database <span className="text-zinc-300">{dbId}</span>
        </p>
      </div>

      {/* Value Card */}
      <div className="border border-zinc-800 rounded-lg overflow-hidden">
        <div className="px-4 py-2 text-sm text-zinc-400 border-b border-zinc-800">
          Value
        </div>
        <EditKeyValue dbId={dbId} keyName={data.key} value={data.value} />
      </div>

      {/* Meta */}
      <div className="flex items-center gap-6 text-sm">
        <div className="text-zinc-400">
          TTL: <span className="text-zinc-200 font-medium">{data.ttl}</span>
        </div>
      </div>

      {/* Actions */}
      <div className="pt-4 border-t border-zinc-800">
        <DeleteKeyButton dbId={dbId} keyName={data.key} />
      </div>
    </section>
  );
}
