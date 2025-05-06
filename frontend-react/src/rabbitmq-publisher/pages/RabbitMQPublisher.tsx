import React, { useState } from 'react';
import { Box, AppBar, Toolbar, Typography, Container } from '@mui/material';
import { Dashboard } from '../components/RabbitMQPublisher/Dashboard';
import { PublishForm } from '../components/RabbitMQPublisher/PublishForm';

export const RabbitMQPublisher: React.FC = () => {
  const [selectedTenant, setSelectedTenant] = useState<string>('');

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBar position="static" color="default" elevation={1}>
        <Toolbar>
          <Typography variant="h6" component="div">
            RabbitMQ Publisher
          </Typography>
        </Toolbar>
      </AppBar>
      <Container maxWidth="lg" sx={{ mt: 4, mb: 4 }}>
        <Dashboard onTenantSelect={setSelectedTenant} />
        {selectedTenant && <PublishForm tenantId={selectedTenant} />}
      </Container>
    </Box>
  );
}; 