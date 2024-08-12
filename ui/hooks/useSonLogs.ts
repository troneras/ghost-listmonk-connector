// hooks/useSonLogs.ts
import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';
import { SonExecutionLog, ActionExecutionLog, Son } from '@/lib/types';

interface Pagination {
    total: number;
    limit: number;
    offset: number;
}

interface EnhancedSonExecutionLog extends SonExecutionLog {
    sonName: string;
}

export function useSonLogs() {
    const [logs, setLogs] = useState<EnhancedSonExecutionLog[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);
    const [pagination, setPagination] = useState<Pagination>({ total: 0, limit: 10, offset: 0 });

    const fetchLogs = useCallback(async (offset = 0) => {
        try {
            setLoading(true);
            const response = await apiClient.get<{ logs: SonExecutionLog[], pagination: Pagination }>(`/son-execution-logs?offset=${offset}&limit=${pagination.limit}`);

            const enhancedLogs = await Promise.all(response.data.logs.map(async (log) => {
                try {
                    const sonResponse = await apiClient.get<Son>(`/sons/${log.son_id}`);
                    return { ...log, sonName: sonResponse.data.name };
                } catch (error) {
                    console.error(`Failed to fetch Son name for ID ${log.son_id}:`, error);
                    return { ...log, sonName: 'Unknown' };
                }
            }));

            setLogs(enhancedLogs);
            setPagination(response.data.pagination);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    }, [pagination.limit]);

    useEffect(() => {
        fetchLogs();
    }, [fetchLogs]);

    const fetchActionLogs = async (executionId: string): Promise<ActionExecutionLog[]> => {
        try {
            const response = await apiClient.get<{ logs: ActionExecutionLog[] }>(`/son-executions/${executionId}/action-logs`);
            return response.data.logs;
        } catch (err) {
            throw err instanceof Error ? err : new Error('An error occurred');
        }
    };

    const fetchNextPage = () => {
        if (pagination.offset + pagination.limit < pagination.total) {
            fetchLogs(pagination.offset + pagination.limit);
        }
    };

    return { logs, loading, error, pagination, fetchLogs, fetchActionLogs, fetchNextPage };
}