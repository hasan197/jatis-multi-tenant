import React, { useState } from 'react';
import { 
  Card, 
  CardContent, 
  Button, 
  Box,
  Alert,
  Snackbar
} from '@mui/material';
import MonacoEditor from '@monaco-editor/react';
import { usePublishMessage } from '../../hooks/usePublishMessage';

interface PublishFormProps {
  tenantId: string;
}

export const PublishForm: React.FC<PublishFormProps> = ({ tenantId }) => {
  const [jsonPayload, setJsonPayload] = useState<string>('{\n  "message": "Hello World",\n  "priority": "high",\n  "metadata": {\n    "source": "web-ui"\n  }\n}');
  const { publishMessage, loading } = usePublishMessage();
  const [snackbar, setSnackbar] = useState<{
    open: boolean;
    message: string;
    severity: 'success' | 'error';
  }>({
    open: false,
    message: '',
    severity: 'success'
  });

  const handleValidate = () => {
    try {
      JSON.parse(jsonPayload);
      setSnackbar({
        open: true,
        message: 'JSON valid',
        severity: 'success'
      });
    } catch (error) {
      setSnackbar({
        open: true,
        message: 'Invalid JSON format',
        severity: 'error'
      });
    }
  };

  const handlePublish = async () => {
    try {
      const payload = JSON.parse(jsonPayload);
      await publishMessage(tenantId, payload);
      setSnackbar({
        open: true,
        message: 'Message published successfully',
        severity: 'success'
      });
    } catch (error) {
      setSnackbar({
        open: true,
        message: 'Failed to publish message',
        severity: 'error'
      });
    }
  };

  return (
    <Card sx={{ mt: 4 }}>
      <CardContent>
        <Box sx={{ mb: 2 }}>
          <MonacoEditor
            height="300px"
            language="json"
            theme="vs-dark"
            value={jsonPayload}
            onChange={(value) => setJsonPayload(value || '')}
            options={{
              minimap: { enabled: false },
              formatOnPaste: true,
              formatOnType: true,
            }}
          />
        </Box>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button 
            variant="outlined" 
            onClick={handleValidate}
          >
            Validate JSON
          </Button>
          <Button 
            variant="contained" 
            onClick={handlePublish}
            disabled={loading}
          >
            Publish Message
          </Button>
        </Box>
      </CardContent>
      <Snackbar 
        open={snackbar.open} 
        autoHideDuration={6000} 
        onClose={() => setSnackbar({ ...snackbar, open: false })}
      >
        <Alert 
          onClose={() => setSnackbar({ ...snackbar, open: false })} 
          severity={snackbar.severity}
        >
          {snackbar.message}
        </Alert>
      </Snackbar>
    </Card>
  );
}; 