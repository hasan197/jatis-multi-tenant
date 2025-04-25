import { useState, useEffect } from 'react';
import { Box, Typography, CircularProgress, Button, Card, CardContent } from '@mui/material';

interface HelloResponse {
  message: string;
  timestamp: string;
}

const HelloWorld = () => {
  const [data, setData] = useState<HelloResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const fetchHelloWorld = async () => {
    setLoading(true);
    setError('');
    
    try {
      const response = await fetch('http://localhost:8080/api/hello-world');
      
      if (!response.ok) {
        throw new Error(`Error: ${response.status}`);
      }
      
      const result = await response.json();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong');
      console.error('Error fetching hello world data:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchHelloWorld();
  }, []);

  return (
    <Box sx={{ p: 4, maxWidth: 600, mx: 'auto', textAlign: 'center' }}>
      <Typography variant="h3" gutterBottom>
        Hello World Page
      </Typography>
      
      <Card sx={{ mb: 3, bgcolor: '#f5f5f5' }}>
        <CardContent>
          {loading ? (
            <CircularProgress size={30} />
          ) : error ? (
            <Typography color="error">{error}</Typography>
          ) : data ? (
            <Box>
              <Typography variant="h4" sx={{ mb: 2 }}>
                {data.message}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Server time: {data.timestamp}
              </Typography>
            </Box>
          ) : (
            <Typography>No data available</Typography>
          )}
        </CardContent>
      </Card>
      
      <Button 
        variant="contained" 
        color="primary" 
        onClick={fetchHelloWorld}
        disabled={loading}
      >
        Refresh
      </Button>
    </Box>
  );
};

export default HelloWorld; 