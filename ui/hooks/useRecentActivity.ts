// ui/hooks/useRecentActivity.ts
import { useState, useEffect } from 'react';
import { apiClient } from '@/lib/api-client';
import { RecentActivity } from '@/lib/types';

export function useRecentActivity() {
    const [activities, setActivities] = useState<RecentActivity[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchActivities = async () => {
        try {
            setLoading(true);
            const response = await apiClient.get<RecentActivity[]>('/recent-activity');
            setActivities(response.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchActivities();
    }, []);

    return { activities, loading, error, fetchActivities };
}