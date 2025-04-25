import { Routes, Route } from 'react-router-dom'
import { ThemeProvider, CssBaseline, Box } from '@mui/material'
import theme from './theme'
import HelloWorld from './pages/HelloWorld'
import Home from './pages/Home'
import Navbar from './components/Navbar'
import UserList from './pages/UserList'
import UserCreate from './pages/UserCreate'
import UserEdit from './pages/UserEdit'

function App() {
  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
        <Navbar />
        <Box component="main" sx={{ flexGrow: 1, p: 2 }}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/hello-world" element={<HelloWorld />} />
            <Route path="/users" element={<UserList />} />
            <Route path="/users/create" element={<UserCreate />} />
            <Route path="/users/edit/:id" element={<UserEdit />} />
          </Routes>
        </Box>
      </Box>
    </ThemeProvider>
  )
}

export default App 