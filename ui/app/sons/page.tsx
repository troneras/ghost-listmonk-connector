"use client";
import Link from "next/link";
import { Button } from "@/components/ui/button";

import { SonList } from "@/components/SonList";

export default function SonListPage() {

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold">Manage Sons</h1>
        <Link href="/sons/new">
          <Button>Create New Son</Button>
        </Link>
      </div>

      <SonList />
    </div>
  );
}
