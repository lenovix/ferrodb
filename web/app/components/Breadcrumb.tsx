import Link from "next/link";

type Crumb = {
  label: string;
  href?: string;
};

type Props = {
  items: Crumb[];
};

export default function Breadcrumb({ items }: Props) {
  return (
    <nav className="text-sm text-zinc-400">
      <ol className="flex items-center gap-2 flex-wrap">
        {items.map((item, i) => (
          <li key={i} className="flex items-center gap-2">
            {item.href ? (
              <Link href={item.href} className="hover:text-zinc-200 transition">
                {item.label}
              </Link>
            ) : (
              <span className="text-zinc-200 font-medium break-all">
                {item.label}
              </span>
            )}

            {i < items.length - 1 && <span className="text-zinc-600">/</span>}
          </li>
        ))}
      </ol>
    </nav>
  );
}
