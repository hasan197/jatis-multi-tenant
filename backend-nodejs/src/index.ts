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

const connectToRabbitMQ = async () => {
  try {
    const connection = await amqp.connect('amqp://guest:guest@rabbitmq:5672');
    channel = await connection.createChannel();
    logger.info('Connected to RabbitMQ');
  } catch (error) {
    logger.error('Failed to connect to RabbitMQ:', error);
  }
};

connectToRabbitMQ();

// Endpoint yang ditangani Node.js
app.post('/api/tenants/:tenantId/publish', async (req, res) => {
  try {
    const { tenantId } = req.params;
    const payload = req.body;

    if (!channel) {
      throw new Error('RabbitMQ channel not initialized');
    }

    const queueName = `tenant_${tenantId}_queue`;
    
    // Ensure queue exists
    await channel.assertQueue(queueName, {
      durable: true
    });

    // Publish message
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