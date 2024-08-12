
// src/hooks/useSons.ts
import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';
import { Son } from '@/lib/types';


export function useSons() {
    const [sons, setSons] = useState<Son[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchSons = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get<Son[]>('/sons');
            setSons(response.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    }, []);

    const createSon = useCallback(async (sonData: Omit<Son, 'id' | 'created_at' | 'updated_at'>) => {
        setLoading(true);
        try {
            const response = await apiClient.post<Son>('/sons', sonData);
            setSons(prevSons => [...prevSons, response.data]);
            setError(null);
            return response.data;
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
            throw err;
        } finally {
            setLoading(false);
        }
    }, []);

    const updateSon = useCallback(async (id: string, sonData: Partial<Son>) => {
        setLoading(true);
        try {
            const response = await apiClient.put<Son>(`/sons/${id}`, sonData);
            setSons(prevSons => prevSons.map(son => son.id === id ? response.data : son));
            setError(null);
            return response.data;
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
            throw err;
        } finally {
            setLoading(false);
        }
    }, []);

    const deleteSon = useCallback(async (id: string) => {
        setLoading(true);
        try {
            await apiClient.delete(`/sons/${id}`);
            setSons(prevSons => prevSons.filter(son => son.id !== id));
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
            throw err;
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchSons();
    }, [fetchSons]);

    return { sons, loading, error, fetchSons, createSon, updateSon, deleteSon };
}