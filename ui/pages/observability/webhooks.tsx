import { NextPageWithExtras } from "next";
import { WebhookTable } from "@/components/WebhookTable";

const WebhooksPage: NextPageWithExtras = () => {
  return (
    <div className="container mx-auto py-10">
      <h1 className="text-2xl font-bold mb-5">Webhook Logs</h1>
      <WebhookTable />
    </div>
  );
};

WebhooksPage.auth = true;

export default WebhooksPage;
