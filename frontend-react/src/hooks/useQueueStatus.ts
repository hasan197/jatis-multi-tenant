import { useState, useEffect } from 'react';

interface QueueStatus {
  status: string;
  workers: number;
  messageCount: number;
}

export const useQueueStatus = (tenantId?: string) => {
  const [queueStatus, setQueueStatus] = useState<QueueStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!tenantId) {
      setQueueStatus(null);
      setLoading(false);
      return;
    }

    const fetchQueueStatus = async () => {
      try {
        const response = await fetch(`/api/tenants/${tenantId}/status`);
        if (!response.ok) {
          throw new Error('Failed to fetch queue status');
        }
        const data = await response.json();
        setQueueStatus(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'An error occurred');
      } finally {
        setLoading(false);
      }
    };

    fetchQueueStatus();
  }, [tenantId]);

  return { queueStatus, loading, error };
}; 