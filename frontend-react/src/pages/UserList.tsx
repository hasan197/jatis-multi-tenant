import { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { 
  Box, 
  Typography, 
  Button, 
  Paper, 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Alert,
  CircularProgress,
  Card,
  CardContent
} from '@mui/material';
import { format } from 'date-fns';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import AddIcon from '@mui/icons-material/Add';

import { getUsers, deleteUser, User } from '../services/userService';

const UserList = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [userToDelete, setUserToDelete] = useState<User | null>(null);
  const [retryCount, setRetryCount] = useState(0);
  
  const navigate = useNavigate();

  // Load users data
  useEffect(() => {
    const fetchUsers = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await getUsers();
        if (data && Array.isArray(data)) {
          setUsers(data);
        } else {
          setError('Format data tidak valid');
          setUsers([]);
        }
      } catch (err: any) {
        console.error('Error fetching users:', err);
        let errorMessage = 'Gagal memuat data pengguna';
        
        if (err.message?.includes('timeout')) {
          errorMessage = 'Waktu permintaan habis. Server mungkin sedang sibuk atau tidak tersedia.';
        } else if (err.message?.includes('Network Error') || err.response?.data?.error?.includes('Proxy Error')) {
          errorMessage = 'Gagal terhubung ke server. Silakan coba lagi nanti.';
        } else if (err.response?.data?.error) {
          errorMessage = err.response.data.error;
        }
        
        setError(errorMessage);
        setUsers([]);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, [retryCount]);

  // Handle retry
  const handleRetry = () => {
    setRetryCount(prev => prev + 1);
  };

  // Handle edit user
  const handleEditUser = (userId: number) => {
    navigate(`/users/edit/${userId}`);
  };

  // Handle open delete dialog
  const handleOpenDeleteDialog = (user: User) => {
    setUserToDelete(user);
    setDeleteDialogOpen(true);
  };

  // Handle close delete dialog
  const handleCloseDeleteDialog = () => {
    setDeleteDialogOpen(false);
    setUserToDelete(null);
  };

  // Handle confirm delete
  const handleConfirmDelete = async () => {
    if (!userToDelete) return;
    
    try {
      await deleteUser(userToDelete.id);
      setUsers(users.filter(user => user.id !== userToDelete.id));
      handleCloseDeleteDialog();
    } catch (err: any) {
      console.error('Error deleting user:', err);
      let errorMessage = 'Gagal menghapus pengguna';
      
      if (err.message?.includes('timeout')) {
        errorMessage = 'Waktu permintaan habis. Server mungkin sedang sibuk atau tidak tersedia.';
      } else if (err.message?.includes('Network Error') || err.response?.data?.error?.includes('Proxy Error')) {
        errorMessage = 'Gagal terhubung ke server. Silakan coba lagi nanti.';
      } else if (err.response?.data?.error) {
        errorMessage = err.response.data.error;
      }
      
      setError(errorMessage);
    }
  };

  // Format date
  const formatDate = (dateString: string) => {
    try {
      return format(new Date(dateString), 'dd/MM/yyyy HH:mm');
    } catch (err) {
      return dateString;
    }
  };

  // Render empty state
  const renderEmptyState = () => (
    <Card sx={{ mt: 2 }}>
      <CardContent sx={{ textAlign: 'center', py: 4 }}>
        <Typography variant="h6" color="text.secondary" gutterBottom>
          Belum ada data pengguna
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          Silakan tambahkan pengguna baru untuk memulai
        </Typography>
        <Button
          component={Link}
          to="/users/create"
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
        >
          Tambah Pengguna Baru
        </Button>
      </CardContent>
    </Card>
  );

  return (
    <Box sx={{ maxWidth: 1200, mx: 'auto', p: 2 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" component="h1">
          Daftar Pengguna
        </Typography>
        <Button 
          component={Link} 
          to="/users/create" 
          variant="contained" 
          color="primary"
          startIcon={<AddIcon />}
        >
          Tambah Pengguna
        </Button>
      </Box>

      {error && (
        <Alert 
          severity="error" 
          sx={{ mb: 2 }}
          action={
            <Button 
              color="inherit" 
              size="small" 
              onClick={handleRetry}
              disabled={loading}
            >
              {loading ? 'Memuat...' : 'Coba Lagi'}
            </Button>
          }
        >
          {error}
        </Alert>
      )}
      
      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      ) : !users || users.length === 0 ? (
        renderEmptyState()
      ) : (
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Nama</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Dibuat Pada</TableCell>
                <TableCell>Diperbarui Pada</TableCell>
                <TableCell align="center">Aksi</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.id}</TableCell>
                  <TableCell>{user.name}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{formatDate(user.created_at)}</TableCell>
                  <TableCell>{formatDate(user.updated_at)}</TableCell>
                  <TableCell align="center">
                    <IconButton 
                      color="primary" 
                      onClick={() => handleEditUser(user.id)}
                      aria-label="edit"
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton 
                      color="error" 
                      onClick={() => handleOpenDeleteDialog(user)}
                      aria-label="delete"
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      {/* Dialog konfirmasi hapus */}
      <Dialog
        open={deleteDialogOpen}
        onClose={handleCloseDeleteDialog}
      >
        <DialogTitle>Konfirmasi Hapus</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Apakah Anda yakin ingin menghapus pengguna "{userToDelete?.name}"?
            Tindakan ini tidak dapat dibatalkan.
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteDialog} color="primary">
            Batal
          </Button>
          <Button onClick={handleConfirmDelete} color="error" autoFocus>
            Hapus
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default UserList; 