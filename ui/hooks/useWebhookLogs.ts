import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api-client';
import { WebhookLog, Pagination } from '@/lib/types';


export function useWebhookLogs() {
    const [logs, setLogs] = useState<WebhookLog[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);
    const [pagination, setPagination] = useState<Pagination | null>(null);

    const fetchLogs = async (offset = 0) => {
        try {
            setLoading(true);
            const response = await apiClient.get(`/webhook-logs?offset=${offset}&limit=10`);
            const newLogs = response.data.logs;
            setLogs(prevLogs => offset === 0 ? newLogs : [...prevLogs, ...newLogs]);
            setPagination(response.data.pagination);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchLogs();
    }, []);

    const fetchNextPage = () => {
        if (pagination && pagination.next_offset !== -1) {
            fetchLogs(pagination.next_offset);
        }
    };

    return { logs, loading, error, fetchNextPage };
}