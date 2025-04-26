import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { 
  Box, 
  Typography, 
  TextField, 
  Button, 
  Paper,
  Grid,
  Alert,
  CircularProgress
} from '@mui/material';

import { getUserById, updateUser, UpdateUserPayload } from '../services/userService';

const UserEdit = () => {
  const [formData, setFormData] = useState<UpdateUserPayload>({
    name: '',
    email: ''
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(false);
  const [fetchLoading, setFetchLoading] = useState(true);
  const [apiError, setApiError] = useState<string | null>(null);
  const [retrying, setRetrying] = useState(false);
  const [retryCount, setRetryCount] = useState(0);
  
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  
  // Fetch user data
  useEffect(() => {
    const fetchUser = async () => {
      if (!id) return;
      
      try {
        setFetchLoading(true);
        setApiError(null);
        const userData = await getUserById(parseInt(id, 10));
        setFormData({
          name: userData.name,
          email: userData.email
        });
      } catch (err: any) {
        console.error(err);
        if (err.response?.status === 404) {
          setApiError('Pengguna tidak ditemukan');
        } else if (err.message?.includes('timeout')) {
          setApiError('Waktu permintaan habis. Server mungkin sedang sibuk atau tidak tersedia.');
        } else if (err.message?.includes('Network Error') || err.response?.data?.error?.includes('Proxy Error')) {
          setApiError('Gagal terhubung ke server. Silakan coba lagi nanti.');
        } else {
          setApiError('Gagal memuat data pengguna: ' + (err.response?.data?.error || err.message || 'Unknown error'));
        }
      } finally {
        setFetchLoading(false);
      }
    };
    
    fetchUser();
  }, [id, retryCount]);
  
  // Fungsi untuk mencoba lagi
  const handleRetry = () => {
    setRetrying(true);
    setRetryCount(prev => prev + 1);
    setTimeout(() => {
      setRetrying(false);
    }, 1000);
  };

  // Handle input change
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
    // Clear error when user types
    if (errors[name]) {
      setErrors({
        ...errors,
        [name]: ''
      });
    }
  };

  // Validate form
  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};
    
    if (!formData.name.trim()) {
      newErrors.name = 'Nama wajib diisi';
    }
    
    if (!formData.email.trim()) {
      newErrors.email = 'Email wajib diisi';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'Format email tidak valid';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  // Handle form submit
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!id) return;
    
    // Validate form
    if (!validateForm()) {
      return;
    }
    
    try {
      setLoading(true);
      setApiError(null);
      
      await updateUser(parseInt(id, 10), formData);
      navigate('/users', { state: { message: 'Pengguna berhasil diperbarui' } });
    } catch (err: any) {
      console.error(err);
      if (err.message?.includes('timeout')) {
        setApiError('Waktu permintaan habis. Server mungkin sedang sibuk atau tidak tersedia.');
      } else if (err.message?.includes('Network Error') || err.response?.data?.error?.includes('Proxy Error')) {
        setApiError('Gagal terhubung ke server. Silakan coba lagi nanti.');
      } else {
        setApiError(err.response?.data?.error || 'Terjadi kesalahan saat memperbarui pengguna');
      }
    } finally {
      setLoading(false);
    }
  };

  if (fetchLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', p: 2 }}>
      <Box sx={{ mb: 3 }}>
        <Typography variant="h4" component="h1">
          Edit Pengguna
        </Typography>
      </Box>
      
      {apiError && (
        <Alert 
          severity="error" 
          sx={{ mb: 2 }}
          action={
            <Button 
              color="inherit" 
              size="small" 
              onClick={handleRetry}
              disabled={retrying}
            >
              {retrying ? 'Mencoba Ulang...' : 'Coba Lagi'}
            </Button>
          }
        >
          {apiError}
        </Alert>
      )}
      
      <Paper sx={{ p: 3 }}>
        <form onSubmit={handleSubmit}>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Nama"
                name="name"
                value={formData.name}
                onChange={handleChange}
                error={!!errors.name}
                helperText={errors.name}
                disabled={loading}
                required
              />
            </Grid>
            
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Email"
                name="email"
                type="email"
                value={formData.email}
                onChange={handleChange}
                error={!!errors.email}
                helperText={errors.email}
                disabled={loading}
                required
              />
            </Grid>
            
            <Grid item xs={12} sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Button
                variant="outlined"
                color="secondary"
                onClick={() => navigate('/users')}
                disabled={loading}
              >
                Batal
              </Button>
              
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={loading}
                startIcon={loading ? <CircularProgress size={20} /> : null}
              >
                Simpan
              </Button>
            </Grid>
          </Grid>
        </form>
      </Paper>
    </Box>
  );
};

export default UserEdit; 