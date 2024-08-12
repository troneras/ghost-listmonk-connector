import { useState, useEffect } from "react";
import { useSonContext } from "@/contexts/SonContext";
import { Son } from "@/lib/types";
import { apiClient } from "@/lib/api-client";

export function useSon(id: string) {
    const { sons, setSons } = useSonContext();  // Assuming you expose setSons in the context
    const [son, setSon] = useState<Son | null>(null);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        // Check if the Son is already in the context
        const existingSon = sons.find((s) => s.id === id);
        if (existingSon) {
            setSon(existingSon);
            setLoading(false);
        } else {
            // If not, fetch from the server
            const fetchSon = async () => {
                setLoading(true);
                try {
                    const response = await apiClient.get<Son>(`/sons/${id}`);
                    setSon(response.data);
                    setError(null);

                    // Optimistically update the context with the new Son
                    setSons((prevSons) => [...prevSons, response.data]);
                } catch (err) {
                    setError(err instanceof Error ? err : new Error("An error occurred"));
                } finally {
                    setLoading(false);
                }
            };

            fetchSon();
        }
    }, [id, sons, setSons]);

    return { son, loading, error };
}
