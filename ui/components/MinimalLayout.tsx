import ContextWrapper from "@/components/ContextWrapper";

export default function MinimalLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <main className="flex min-h-screen p-4">{children}</main>;
}
