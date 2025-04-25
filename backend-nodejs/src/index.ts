import express, { Request, Response } from 'express'
import dotenv from 'dotenv'
import { createLogger, format, transports } from 'winston'
import cors from 'cors'
import helmet from 'helmet'
import morgan from 'morgan'
import proxyRouter from './routes/proxy'

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

// Parse JSON request bodies - harus ada sebelum proxyRouter
app.use(express.json({ limit: '10mb' }))
app.use(express.urlencoded({ extended: true, limit: '10mb' }))

// Proxy routes to Golang backend
app.use(proxyRouter)

// Routes
app.get('/', (req: Request, res: Response) => {
  res.json({ message: 'Welcome to Node.js Backend' })
})

// Health check endpoint
app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'ok' })
})

// Start server
app.listen(port, () => {
  logger.info(`Server is running on port ${port}`)
  logger.info(`API Proxy is forwarding requests to ${process.env.GOLANG_API_URL || 'http://backend-golang:8080'}`)
}) 