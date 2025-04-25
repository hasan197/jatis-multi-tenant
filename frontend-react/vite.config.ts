import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: true,
    port: 5173,
    proxy: {
      // Proxy semua permintaan /api ke server Node.js
      '/api': {
        target: 'http://backend-nodejs:3000',
        changeOrigin: true,
        secure: false
      }
    }
  },
}) 