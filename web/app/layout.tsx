import "./globals.css";
import Sidebar from "./components/Sidebar";

export const metadata = {
  title: "FerroDB Admin",
  description: "FerroDB Admin Console",
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // fetch jumlah DB sekali di root
  const res = await fetch("http://localhost:8080/api/db", {
    cache: "no-store",
  });

  const data = await res.json();
  const dbCount = data.db_count ?? 1;

  return (
    <html lang="en" className="dark">
      <body className="bg-zinc-900 text-zinc-100">
        <div className="flex min-h-screen">
          <Sidebar dbCount={dbCount} />
          <main className="flex-1 p-6 overflow-auto">{children}</main>
        </div>
      </body>
    </html>
  );
}
