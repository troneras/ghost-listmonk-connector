import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';
import { Son, Webhook, ListmonkList, ListmonkTemplate } from '@/lib/types';

interface InitialData {
    sons: Son[];
    webhook: Webhook;
    lists: ListmonkList[];
    templates: ListmonkTemplate[];
}

export function useInitialData() {
    const [data, setData] = useState<InitialData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchInitialData = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get<InitialData>('/initial-data');
            setData(response.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchInitialData();
    }, [fetchInitialData]);

    return { data, loading, error, fetchInitialData };
}