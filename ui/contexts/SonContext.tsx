// ui/contexts/SonContext.tsx
import React, { createContext, useContext } from "react";
import { Son } from "@/lib/types";
import { useSons } from "@/hooks/useSons";
import { Dispatch, SetStateAction } from "react";
import { useWebhook } from "@/hooks/useWebhook";
import { Webhook } from "@/lib/types";

interface SonContextProps {
  sons: Son[];
  setSons: Dispatch<SetStateAction<Son[]>>;
  loading: boolean;
  error: Error | null;
  fetchSons: () => void;
  createSon: (
    sonData: Omit<Son, "id" | "created_at" | "updated_at">
  ) => Promise<Son>;
  updateSon: (id: string, sonData: Partial<Son>) => Promise<Son>;
  deleteSon: (id: string) => Promise<void>;
  webhook: Webhook | null;
  webhookLoading: boolean;
  webhookError: Error | null;
  fetchWebhook: () => Promise<void>;
}

const SonContext = createContext<SonContextProps | undefined>(undefined);

export const SonProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const {
    sons,
    setSons,
    loading,
    error,
    fetchSons,
    createSon,
    updateSon,
    deleteSon,
  } = useSons();
  const {
    webhook,
    loading: webhookLoading,
    error: webhookError,
    fetchWebhook,
  } = useWebhook();

  return (
    <SonContext.Provider
      value={{
        sons,
        setSons,
        loading,
        error,
        fetchSons,
        createSon,
        updateSon,
        deleteSon,
        webhook,
        webhookLoading,
        webhookError,
        fetchWebhook,
      }}
    >
      {children}
    </SonContext.Provider>
  );
};

export const useSonContext = () => {
  const context = useContext(SonContext);
  if (!context) {
    throw new Error("useSonContext must be used within a SonProvider");
  }
  return context;
};
