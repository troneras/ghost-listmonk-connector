import React, { createContext, useContext } from "react";
import { ListmonkList } from "@/lib/types";
import { useLists } from "@/hooks/useLists";

interface ListContextProps {
  lists: ListmonkList[];
  loading: boolean;
  error: Error | null;
  fetchLists: () => void;
}

const ListContext = createContext<ListContextProps | undefined>(undefined);

export const ListProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { lists, loading, error, fetchLists } = useLists();

  return (
    <ListContext.Provider
      value={{
        lists,
        loading,
        error,
        fetchLists,
      }}
    >
      {children}
    </ListContext.Provider>
  );
};

export const useListContext = () => {
  const context = useContext(ListContext);
  if (!context) {
    throw new Error("useListContext must be used within a ListProvider");
  }
  return context;
};
