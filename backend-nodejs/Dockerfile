# Build stage
FROM node:18-alpine AS builder
RUN apk add --no-cache npm

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Development stage
FROM node:18-alpine AS development
RUN apk add --no-cache npm

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies including devDependencies
RUN npm install

# Copy source code
COPY . .

# Ensure TypeScript definitions are installed
RUN npm i --save-dev @types/express @types/cors @types/morgan

# Expose port
EXPOSE 3000

# Start the application in development mode
CMD ["npm", "run", "dev"]

# Production stage
FROM node:18-alpine AS production
RUN apk add --no-cache npm

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install production dependencies only
RUN npm install --production

# Copy built assets from builder stage
COPY --from=builder /app/dist ./dist

# Expose port
EXPOSE 3000

# Start the application
CMD ["node", "dist/index.js"] 