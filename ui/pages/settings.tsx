import { NextPageWithExtras } from "next";
import { WebhookInfo } from "@/components/WebhookInfo";

const SettingsPage: NextPageWithExtras = () => {
  return (
    <div className="flex-1 space-y-4 max-w-3xl mx-auto">
      <h1 className="text-2xl font-bold">Settings</h1>
      <WebhookInfo />
    </div>
  );
};

SettingsPage.auth = true;

export default SettingsPage;
