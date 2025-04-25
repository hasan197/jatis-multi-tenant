import { Box, Typography, Card, CardContent, Button } from '@mui/material';
import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', p: 2, textAlign: 'center' }}>
      <Typography variant="h3" gutterBottom>
        Selamat Datang di Sample Stack
      </Typography>
      
      <Typography variant="subtitle1" sx={{ mb: 4 }}>
        Aplikasi contoh menggunakan Golang, Node.js, React, PostgreSQL, Redis, dan RabbitMQ
      </Typography>
      
      <Card sx={{ mb: 4 }}>
        <CardContent>
          <Typography variant="h5" gutterBottom>
            Fitur Demo
          </Typography>
          
          <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 2, justifyContent: 'center', mt: 2 }}>
            <Button 
              component={Link} 
              to="/hello-world" 
              variant="contained" 
              color="primary"
              size="large"
            >
              Hello World Demo
            </Button>
          </Box>
        </CardContent>
      </Card>
      
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mb: 4 }}>
        <Typography variant="h6">
          Stack Teknologi:
        </Typography>
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2, justifyContent: 'center' }}>
          <Chip label="Golang" />
          <Chip label="Node.js" />
          <Chip label="React" />
          <Chip label="PostgreSQL" />
          <Chip label="Redis" />
          <Chip label="RabbitMQ" />
        </Box>
      </Box>
    </Box>
  );
};

const Chip = ({ label }: { label: string }) => (
  <Card sx={{ minWidth: 120, bgcolor: '#f5f5f5' }}>
    <CardContent sx={{ py: 1, textAlign: 'center' }}>
      <Typography>{label}</Typography>
    </CardContent>
  </Card>
);

export default Home; 