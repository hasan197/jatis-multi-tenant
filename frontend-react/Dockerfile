# Build stage
FROM node:18-alpine AS builder
RUN apk add --no-cache npm

WORKDIR /app

# Copy package files first (lebih baik untuk caching Docker)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code setelah install dependencies
COPY . .

# Build the application
RUN npm run build

# Development stage
FROM node:18-alpine AS development
RUN apk add --no-cache npm

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies terlebih dahulu
RUN npm install

# Install vite sebagai dependency global setelah install local dependencies
RUN npm install -g vite@5.0.12

# Install icons material yang sering menyebabkan masalah
RUN npm install @mui/icons-material

# Copy source code setelah install dependencies
# Tetapi exclude node_modules dengan .dockerignore
COPY . .

# Expose port
EXPOSE 5173

# Start the application in development mode
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# Production stage
FROM nginx:alpine AS production

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]