import { usePathname } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  LayoutDashboard,
  ListPlus,
  ListTree,
  Settings,
  Activity,
} from "lucide-react";

export const Sidebar = ({ isOpen }: { isOpen: boolean }) => {
  const pathname = usePathname();

  const mainMenuItems = [
    { href: "/", label: "Dashboard", icon: LayoutDashboard },
    { href: "/sons/new", label: "Create Son", icon: ListPlus },
    { href: "/sons", label: "Manage Sons", icon: ListTree },
    { href: "/settings", label: "Settings", icon: Settings },
  ];

  const observabilityMenuItems = [
    { href: "/observability/webhooks", label: "Webhook Logs", icon: Activity },
    { href: "/observability/son-logs", label: "Son Logs", icon: Activity },
  ];

  return (
    <aside
      className={`bg-primary text-primary-foreground w-64 min-h-screen p-4 ${
        isOpen ? "" : "hidden"
      } md:block`}
    >
      <nav>
        <ul>
          {mainMenuItems.map((item) => (
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
          <li className="mt-4 mb-4">
            <span className="text-sm font-semibold text-primary-foreground-50">
              Observability
            </span>
          </li>
          {observabilityMenuItems.map((item) => (
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
