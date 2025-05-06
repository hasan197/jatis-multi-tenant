import React from 'react';
import { 
  Card, 
  CardContent, 
  Grid, 
  Typography,
  CircularProgress,
  Box,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper
} from '@mui/material';
import { useQueueStatus } from '../../../hooks/useQueueStatus';

interface MonitoringProps {
  tenantId: string;
}

export const Monitoring: React.FC<MonitoringProps> = ({ tenantId }) => {
  const { queueStatus, loading, refetch } = useQueueStatus(tenantId);

  if (!tenantId) {
    return (
      <Box sx={{ mt: 4 }}>
        <Typography color="error">Tenant ID tidak ditemukan</Typography>
      </Box>
    );
  }

  return (
    <Box sx={{ mt: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">Monitoring</Typography>
        <Button 
          variant="outlined" 
          onClick={() => refetch()}
          disabled={loading}
        >
          Refresh
        </Button>
      </Box>

      <Grid container spacing={3}>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Status Queue
              </Typography>
              {loading ? (
                <CircularProgress size={20} />
              ) : (
                <Typography variant="h5">
                  {queueStatus?.status || 'Unknown'}
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Active Workers
              </Typography>
              {loading ? (
                <CircularProgress size={20} />
              ) : (
                <Typography variant="h5">
                  {queueStatus?.workers || 0}
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Messages in Queue
              </Typography>
              {loading ? (
                <CircularProgress size={20} />
              ) : (
                <Typography variant="h5">
                  {queueStatus?.messageCount || 0}
                </Typography>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Card sx={{ mt: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Queue Details
          </Typography>
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Metric</TableCell>
                  <TableCell>Value</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                <TableRow>
                  <TableCell>Queue Name</TableCell>
                  <TableCell>{`tenant_${tenantId}_queue`}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>Consumer Count</TableCell>
                  <TableCell>{queueStatus?.consumerCount || 0}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>Processing Rate</TableCell>
                  <TableCell>{queueStatus?.processingRate || 'N/A'}</TableCell>
                </TableRow>
                <TableRow>
                  <TableCell>Last Updated</TableCell>
                  <TableCell>{new Date().toLocaleString()}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </TableContainer>
        </CardContent>
      </Card>
    </Box>
  );
}; 