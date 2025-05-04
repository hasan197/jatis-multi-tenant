import axios from 'axios';

// Buat instance Axios dengan konfigurasi timeout
const api = axios.create({
  // Gunakan relative URL untuk memanggil ke server Node.js yang sama
  baseURL: '/api',
  timeout: 10000, // 10 detik timeout
  headers: {
    'Content-Type': 'application/json'
  }
});

export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface CreateUserPayload {
  name: string;
  email: string;
  password: string;
}

export interface UpdateUserPayload {
  name: string;
  email: string;
}

// Get all users
export const getUsers = async (): Promise<User[]> => {
  try {
    const response = await api.get('/users');
    // Pastikan response.data.data ada dan valid
    if (response.data && response.data.data) {
      return Array.isArray(response.data.data) ? response.data.data : [];
    }
    return [];
  } catch (error) {
    console.error('Error fetching users:', error);
    throw error;
  }
};

// Get user by ID
export const getUserById = async (id: number): Promise<User> => {
  try {
    const response = await api.get(`/users/${id}`);
    return response.data.data;
  } catch (error) {
    console.error(`Error fetching user with ID ${id}:`, error);
    throw error;
  }
};

// Create new user
export const createUser = async (userData: CreateUserPayload): Promise<User> => {
  try {
    const response = await api.post('/users', userData);
    return response.data.data;
  } catch (error) {
    console.error('Error creating user:', error);
    throw error;
  }
};

// Update user
export const updateUser = async (id: number, userData: UpdateUserPayload): Promise<User> => {
  try {
    const response = await api.put(`/users/${id}`, userData);
    return response.data.data;
  } catch (error) {
    console.error(`Error updating user with ID ${id}:`, error);
    throw error;
  }
};

// Delete user
export const deleteUser = async (id: number): Promise<void> => {
  try {
    await api.delete(`/users/${id}`);
  } catch (error) {
    console.error(`Error deleting user with ID ${id}:`, error);
    throw error;
  }
}; 