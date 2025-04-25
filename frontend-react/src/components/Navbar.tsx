import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { Link } from 'react-router-dom';

const Navbar = () => {
  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          Sample Stack App
        </Typography>
        <Box>
          <Button color="inherit" component={Link} to="/">Home</Button>
          <Button color="inherit" component={Link} to="/hello-world">Hello World</Button>
          <Button color="inherit" component={Link} to="/users">Users</Button>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Navbar; 