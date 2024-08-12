import { useState, useEffect, useCallback } from 'react';
import { apiClient } from '@/lib/api-client';
import { ListmonkTemplate } from '@/lib/types';



export function useTemplates() {
    const [templates, setTemplates] = useState<ListmonkTemplate[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const fetchTemplates = useCallback(async () => {
        setLoading(true);
        try {
            const response = await apiClient.get<{ data: ListmonkTemplate[] }>('/templates');
            setTemplates(response.data.data);
            setError(null);
        } catch (err) {
            setError(err instanceof Error ? err : new Error('An error occurred'));
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        fetchTemplates();
    }, [fetchTemplates]);

    return { templates, loading, error, fetchTemplates };
}