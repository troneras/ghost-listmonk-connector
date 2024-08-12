import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';
import { Webhook } from '@/lib/types';


export function useWebhook() {
    const [webhook, setWebhook] = useState<Webhook | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchWebhook = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get<{ data: Webhook }>('/webhook-info');
            setWebhook(response.data.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchWebhook();
    }, [fetchWebhook]);

    return { webhook, loading, error, fetchWebhook };
}