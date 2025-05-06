import React, { useState } from 'react';
import { Card, CardContent, Typography, Button, Box } from '@mui/material';
import Editor from '@monaco-editor/react';
import { usePublishMessage } from '../../../hooks/usePublishMessage';

interface PublishFormProps {
  tenantId: string;
}

export const PublishForm: React.FC<PublishFormProps> = ({ tenantId }) => {
  const [message, setMessage] = useState('{\n  "message": "Hello World",\n  "priority": "high",\n  "metadata": {\n    "source": "web-ui"\n  }\n}');
  const { publish, loading, error } = usePublishMessage();

  const handlePublish = async () => {
    try {
      const jsonMessage = JSON.parse(message);
      await publish(tenantId, jsonMessage);
      setMessage('{\n  "message": "Hello World",\n  "priority": "high",\n  "metadata": {\n    "source": "web-ui"\n  }\n}');
    } catch (err) {
      console.error('Invalid JSON:', err);
    }
  };

  return (
    <Card sx={{ mt: 4 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Publish Message
        </Typography>
        <Box sx={{ height: 300, mb: 2 }}>
          <Editor
            height="100%"
            defaultLanguage="json"
            value={message}
            onChange={(value: string | undefined) => setMessage(value || '')}
            options={{
              minimap: { enabled: false },
              formatOnPaste: true,
              formatOnType: true
            }}
          />
        </Box>
        {error && (
          <Typography color="error" sx={{ mb: 2 }}>
            {error}
          </Typography>
        )}
        <Button
          variant="contained"
          onClick={handlePublish}
          disabled={loading}
        >
          {loading ? 'Publishing...' : 'Publish Message'}
        </Button>
      </CardContent>
    </Card>
  );
}; 