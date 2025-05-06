import express from 'express'
import dotenv from 'dotenv'
import { createLogger, format, transports } from 'winston'
import cors from 'cors'
import helmet from 'helmet'
import morgan from 'morgan'
import { createProxyMiddleware } from 'http-proxy-middleware'
import amqp from 'amqplib'

// Load environment variables
dotenv.config()

// Create logger
const logger = createLogger({
  format: format.combine(
    format.timestamp(),
    format.json()
  ),
  transports: [
    new transports.Console()
  ]
})

const app = express()
const port = process.env.PORT || 3000

// Middleware
app.use(helmet())
app.use(cors())
app.use(morgan('combined'))

// Parse JSON request bodies
app.use(express.json({ limit: '10mb' }))
app.use(express.urlencoded({ extended: true, limit: '10mb' }))

// RabbitMQ connection
let channel: amqp.Channel | null = null;
let connection: amqp.Connection | null = null;
let isConnecting = false;
let retryCount = 0;
const maxRetries = 10;
const retryInterval = 3000; // 3 seconds

const connectToRabbitMQ = async () => {
  if (isConnecting) return;
  isConnecting = true;
  
  try {
    logger.info('Connecting to RabbitMQ...');
    connection = await amqp.connect('amqp://guest:guest@rabbitmq:5672');
    
    // Handle connection close/error
    connection.on('error', (err: Error) => {
      logger.error('RabbitMQ connection error:', err);
      reconnect();
    });
    
    connection.on('close', () => {
      logger.warn('RabbitMQ connection closed, attempting to reconnect...');
      reconnect();
    });
    
    // Create channel
    channel = await connection.createChannel();
    
    // Handle channel close/error
    channel.on('error', (err: Error) => {
      logger.error('RabbitMQ channel error:', err);
      channel = null;
      reconnect();
    });
    
    channel.on('close', () => {
      logger.warn('RabbitMQ channel closed');
      channel = null;
      reconnect();
    });
    
    // Reset retry count on successful connection
    retryCount = 0;
    logger.info('Successfully connected to RabbitMQ');
  } catch (error) {
    logger.error('Failed to connect to RabbitMQ:', error);
    reconnect();
  } finally {
    isConnecting = false;
  }
};

const reconnect = () => {
  if (retryCount >= maxRetries) {
    logger.error(`Failed to connect to RabbitMQ after ${maxRetries} attempts. Giving up.`);
    return;
  }
  
  // Close existing connections if any
  if (channel) {
    try { channel.close(); } catch (e) { /* ignore */ }
    channel = null;
  }
  
  if (connection) {
    try { connection.close(); } catch (e) { /* ignore */ }
    connection = null;
  }
  
  retryCount++;
  const delay = retryInterval * Math.min(retryCount, 10); // Exponential backoff capped at 30 seconds
  
  logger.info(`Attempting to reconnect to RabbitMQ in ${delay/1000} seconds... (Attempt ${retryCount}/${maxRetries})`);
  setTimeout(connectToRabbitMQ, delay);
};

// Initial connection
connectToRabbitMQ();

// Endpoint yang ditangani Node.js
app.post('/api/tenants/:tenantId/publish', async (req, res) => {
  try {
    const { tenantId } = req.params;
    const payload = req.body;

    // Jika channel tidak tersedia, coba reconnect dan tunggu
    if (!channel) {
      logger.warn('RabbitMQ channel not available, attempting to reconnect...');
      await connectToRabbitMQ();
      
      // Tunggu sebentar untuk koneksi
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Periksa lagi setelah mencoba reconnect
      if (!channel) {
        throw new Error('RabbitMQ channel still not initialized after reconnect attempt');
      }
    }

    // Format queue name sesuai dengan format yang digunakan di backend Go
    const queueName = `tenant.${tenantId}`;
    
    // Langsung kirim pesan tanpa mendeklarasikan queue
    // Queue sudah dideklarasikan oleh backend Go dengan konfigurasi DLQ
    channel.sendToQueue(queueName, Buffer.from(JSON.stringify(payload)), {
      persistent: true
    });

    logger.info(`Message published to queue ${queueName}`);
    res.json({ message: 'Message published successfully' });
  } catch (error) {
    logger.error('Failed to publish message:', error);
    res.status(500).json({ error: 'Failed to publish message' });
  }
});

// Proxy ke backend Golang untuk endpoint lainnya
app.use('/api', createProxyMiddleware({
  target: 'http://backend-golang:8080',
  changeOrigin: true,
  pathRewrite: (path) => {
    // Skip proxy untuk endpoint publish
    if (path.includes('/publish')) {
      return path; // Biarkan path asli untuk endpoint publish
    }
    return path.replace('^/api', '/api'); // Proxy ke backend Golang untuk endpoint lainnya
  },
  proxyTimeout: 30000,
  timeout: 30000
}))

// Routes
app.get('/', (req, res) => {
  res.json({ message: 'Welcome to Node.js Backend' })
})

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'ok' })
})

// Start server
app.listen(port, () => {
  logger.info(`Node.js proxy server running on port ${port}`)
  logger.info(`API Proxy is forwarding requests to ${process.env.GOLANG_API_URL || 'http://backend-golang:8080'}`)
}) 