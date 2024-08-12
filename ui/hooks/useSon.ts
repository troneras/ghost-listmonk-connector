

// src/hooks/useSon.ts
import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';
import { Son } from '@/lib/types';

export function useSon(id: string) {
    const [son, setSon] = useState<Son | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<Error | null>(null);

    const fetchSon = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get<Son>(`/sons/${id}`);
            setSon(response.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    }, [id]);

    useEffect(() => {
        fetchSon();
    }, [fetchSon]);

    return { son, loading, error, fetchSon };
}