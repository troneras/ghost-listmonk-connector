import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";
import ContextWrapper from "@/components/ContextWrapper";
import { Sidebar } from "@/components/Sidebar";
import { AuthButton } from "./AuthButton";

export default function Layout({ children }: { children: React.ReactNode }) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  return (
    <div className="flex">
      <Sidebar isOpen={isSidebarOpen} />
      <div className="flex-1 flex flex-col">
        <header className="bg-white shadow p-4 flex justify-between items-center">
          <Button
            variant="ghost"
            className="md:hidden"
            onClick={() => setIsSidebarOpen(!isSidebarOpen)}
          >
            <Menu className="h-6 w-6" />
          </Button>
          <h1 className="text-xl font-bold">Ghost-Listmonk Connector</h1>
          <AuthButton />
        </header>
        <main className="flex-1 bg-gray-100 p-4 overflow-y-auto">
          <ContextWrapper>{children}</ContextWrapper>
        </main>
      </div>
    </div>
  );
}
