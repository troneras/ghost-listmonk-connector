"use client";
import { useState } from "react";
import { usePathname } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  LayoutDashboard,
  ListPlus,
  ListTree,
  Settings,
  Menu,
} from "lucide-react";

const Sidebar = ({ isOpen }: { isOpen: boolean }) => {
  const pathname = usePathname();

  const menuItems = [
    { href: "/", label: "Dashboard", icon: LayoutDashboard },
    { href: "/sons/new", label: "Create Son", icon: ListPlus },
    { href: "/sons", label: "Manage Sons", icon: ListTree },
    { href: "/settings", label: "Settings", icon: Settings },
  ];

  return (
    <aside
      className={`bg-gray-800 text-white w-64 min-h-screen p-4 ${
        isOpen ? "" : "hidden"
      } md:block`}
    >
      <nav>
        <ul>
          {menuItems.map((item) => (
            <li key={item.href} className="mb-2">
              <Link href={item.href}>
                <Button
                  variant={pathname === item.href ? "secondary" : "ghost"}
                  className="w-full justify-start"
                >
                  <item.icon className="mr-2 h-4 w-4" />
                  {item.label}
                </Button>
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
};

export default function Layout({ children }: { children: React.ReactNode }) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  return (
    <div className="flex min-h-screen">
      <Sidebar isOpen={isSidebarOpen} />
      <div className="flex-1">
        <header className="bg-white shadow p-4 flex justify-between items-center">
          <Button
            variant="ghost"
            className="md:hidden"
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
          >
            <Menu className="h-6 w-6" />
          </Button>
          <h1 className="text-xl font-bold">Ghost-Listmonk Connector</h1>
        </header>
        <main className="p-4">{children}</main>
      </div>
    </div>
  );
}
