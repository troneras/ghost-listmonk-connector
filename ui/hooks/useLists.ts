
import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';
import { ListmonkList } from '@/lib/types';


export function useLists() {
    const [lists, setLists] = useState<ListmonkList[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchLists = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get<{ data: ListmonkList[] }>('/lists');
            setLists(response.data.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchLists();
    }, [fetchLists]);

    return { lists, loading, error, fetchLists };
}