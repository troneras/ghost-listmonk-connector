import React, { createContext, useContext } from "react";
import { ListmonkTemplate } from "@/lib/types";
import { useTemplates } from "@/hooks/useTemplates";

interface TemplateContextProps {
  templates: ListmonkTemplate[];
  loading: boolean;
  error: Error | null;
  fetchTemplates: () => void;
}

const TemplateContext = createContext<TemplateContextProps | undefined>(
  undefined
);

export const TemplateProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { templates, loading, error, fetchTemplates } = useTemplates();

  return (
    <TemplateContext.Provider
      value={{
        templates,
        loading,
        error,
        fetchTemplates,
      }}
    >
      {children}
    </TemplateContext.Provider>
  );
};

export const useTemplateContext = () => {
  const context = useContext(TemplateContext);
  if (!context) {
    throw new Error(
      "useTemplateContext must be used within a TemplateProvider"
    );
  }
  return context;
};
