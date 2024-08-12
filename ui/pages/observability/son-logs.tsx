// pages/observability/son-logs.tsx
import { NextPageWithExtras } from "next";
import SonLogsPage from "@/components/SonLogsPage";

const SonLogsPageWrapper: NextPageWithExtras = () => {
  return <SonLogsPage />;
};

SonLogsPageWrapper.auth = true;

export default SonLogsPageWrapper;
