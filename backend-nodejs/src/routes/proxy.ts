import express, { Request, Response, NextFunction } from 'express';
import { createProxyMiddleware, responseInterceptor, Options } from 'http-proxy-middleware';
import http from 'http';

const router = express.Router();
// Gunakan backend-golang untuk nama host dalam jaringan Docker
const golangApiUrl = process.env.GOLANG_API_URL || 'http://backend-golang:8080';

// Konfigurasi agen HTTP dengan timeout lebih agresif
const httpAgent = new http.Agent({
  keepAlive: true,
  maxSockets: 25,
  timeout: 10000 // 10 detik timeout pada level socket
});

// Proxy semua permintaan ke /api/* ke backend Golang
router.use('/api', createProxyMiddleware({
  target: golangApiUrl,
  changeOrigin: true,
  pathRewrite: {
    '^/api': '/api', // tidak perlu mengubah path karena sudah sama
  },
  // Konfigurasi timeout yang lebih pendek
  timeout: 15000,
  // Pastikan body request dan respons ditangani dengan benar
  onProxyReq: (proxyReq, req, res) => {
    // Log permintaan untuk debugging
    console.log(`Proxying ${req.method} request to ${golangApiUrl}${proxyReq.path}`);

    // Jika ada body dan content-type adalah application/json
    if (req.body && Object.keys(req.body).length > 0) {
      const bodyData = JSON.stringify(req.body);
      // Update header Content-Length dengan ukuran body yang benar
      proxyReq.setHeader('Content-Length', Buffer.byteLength(bodyData));
      // Tulis body ke proxyReq
      proxyReq.write(bodyData);
      proxyReq.end();
    }
  },
  // Tangani error
  onError: (err, req, res) => {
    console.error(`Proxy error: ${err.message}`);
    // Pastikan respons belum dikirim
    if (!res.headersSent) {
      res.writeHead(500, {'Content-Type': 'application/json'});
      res.end(JSON.stringify({
        error: 'Proxy Error',
        message: 'Tidak dapat terhubung ke backend service',
        details: err.message
      }));
    }
  }
} as Options));

export default router; 