// ui/hooks/useSonStats.ts

import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api-client';
import { SonStats } from '@/lib/types';

export function useSonStats(timeframe: string) {
    const [stats, setStats] = useState<SonStats[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchStats = async () => {
        try {
            setLoading(true);
            const response = await apiClient.get<SonStats[]>(`/son-stats?timeframe=${timeframe}`);
            setStats(response.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchStats();
    }, [timeframe]);

    return { stats, loading, error, fetchStats };
}