import type { NextPageWithExtras } from "next";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { SonList } from "@/components/SonList";

const SonListPage: NextPageWithExtras = () => {
  return (
    <div className="flex-1 space-y-4 max-w-3xl mx-auto">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Manage Sons</h1>
        <Link href="/sons/new">
          <Button>Create New Son</Button>
        </Link>
      </div>

      <SonList />
    </div>
  );
};

SonListPage.auth = true;

export default SonListPage;
