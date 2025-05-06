import { useState, useEffect } from 'react';

interface QueueStatus {
  status: string;
  workers: number;
  messageCount: number;
  consumerCount: number;
  processingRate: string;
}

export const useQueueStatus = (tenantId?: string) => {
  const [queueStatus, setQueueStatus] = useState<QueueStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchQueueStatus = async () => {
    if (!tenantId) return;
    
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/tenants/${tenantId}/queue-status`);
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

  useEffect(() => {
    fetchQueueStatus();
    const interval = setInterval(fetchQueueStatus, 5000);
    return () => clearInterval(interval);
  }, [tenantId]);

  return { queueStatus, loading, error, refetch: fetchQueueStatus };
}; 