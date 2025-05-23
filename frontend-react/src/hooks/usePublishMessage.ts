import { useState } from 'react';

export const usePublishMessage = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const publish = async (tenantId: string, payload: any) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`/api/tenants/${tenantId}/publish`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });
      if (!response.ok) {
        throw new Error('Failed to publish message');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { publish, loading, error };
}; 