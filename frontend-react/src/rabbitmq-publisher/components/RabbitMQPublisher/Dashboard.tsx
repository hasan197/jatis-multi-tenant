import React from 'react';
import { 
  Card, 
  CardContent, 
  FormControl, 
  InputLabel, 
  Select, 
  MenuItem, 
  Grid, 
  Typography,
  CircularProgress
} from '@mui/material';
import { useTenants } from '../../hooks/useTenants';
import { useQueueStatus } from '../../hooks/useQueueStatus';

interface DashboardProps {
  onTenantSelect: (tenantId: string) => void;
}

export const Dashboard: React.FC<DashboardProps> = ({ onTenantSelect }) => {
  const { tenants, loading: tenantsLoading } = useTenants();
  const { queueStatus, loading: statusLoading } = useQueueStatus();

  return (
    <>
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <FormControl fullWidth>
            <InputLabel>Pilih Tenant</InputLabel>
            <Select
              label="Pilih Tenant"
              onChange={(e) => onTenantSelect(e.target.value)}
              disabled={tenantsLoading}
            >
              {tenants.map((tenant) => (
                <MenuItem key={tenant.id} value={tenant.id}>
                  {tenant.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </CardContent>
      </Card>

      {queueStatus && (
        <Grid container spacing={3}>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="textSecondary" gutterBottom>
                  Status Tenant
                </Typography>
                {statusLoading ? (
                  <CircularProgress size={20} />
                ) : (
                  <Typography variant="h5">
                    {queueStatus.status}
                  </Typography>
                )}
              </CardContent>
            </Card>
          </Grid>
          <Grid item xs={12} md={4}>
            <Card>
              <CardContent>
                <Typography color="textSecondary" gutterBottom>
                  Workers
                </Typography>
                {statusLoading ? (
                  <CircularProgress size={20} />
                ) : (
                  <Typography variant="h5">
                    {queueStatus.workers}
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
                {statusLoading ? (
                  <CircularProgress size={20} />
                ) : (
                  <Typography variant="h5">
                    {queueStatus.messageCount}
                  </Typography>
                )}
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}
    </>
  );
}; 