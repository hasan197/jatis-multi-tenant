# Cursor AI Prompts Guide

## 1. Setup Awal & Database

### 1.1 Setup Project Go
```prompt
Buatkan struktur project Go dengan:
- Clean Architecture
- Dependency injection
- Environment configuration
- Logging setup
- Error handling middleware
```

### 1.2 Database Schema
```prompt
Buatkan migration SQL untuk:
- Tabel messages dengan partisi berdasarkan tenant_id
- Index untuk optimasi query
- Trigger untuk auto-partition
- Fungsi helper untuk manajemen partisi
```

### 1.3 Repository Layer
```prompt
Implementasikan repository interface untuk messages dengan:
- CRUD operations
- Cursor-based pagination
- Error handling
- Transaction support
- Unit tests
```

## 2. Backend Core

### 2.1 RabbitMQ Integration
```prompt
Buatkan service untuk RabbitMQ dengan:
- Connection management
- Channel pooling
- Error handling
- Reconnection logic
- Health check
- Unit tests
```

### 2.2 Tenant Management
```prompt
Implementasikan TenantManager service dengan:
- Tenant lifecycle management
- Queue creation/deletion
- Consumer management
- Concurrency control
- Health monitoring
- Unit tests
```

### 2.3 Message Processing
```prompt
Buatkan message processor dengan:
- Worker pool pattern
- Message validation
- Error handling
- Retry mechanism
- Dead letter queue
- Unit tests
```

### 2.4 API Endpoints
```prompt
Implementasikan REST API dengan:
- Tenant endpoints (CRUD)
- Message endpoints
- Authentication middleware
- Request validation
- Response formatting
- API documentation
- Integration tests
```

## 3. Frontend Core

### 3.1 React Setup
```prompt
Setup React project dengan:
- TypeScript
- Vite
- Tailwind CSS
- React Query
- React Router
- State management
- Component library
```

### 3.2 Layout Components
```prompt
Buatkan layout components:
- Main layout
- Sidebar
- Header
- Content area
- Responsive design
- Theme provider
- Error boundary
```

### 3.3 Tenant Management UI
```prompt
Implementasikan UI untuk tenant management:
- Tenant list
- Tenant form
- Status indicators
- Quick actions
- Confirmation dialogs
- Error handling
- Loading states
```

### 3.4 Message UI
```prompt
Buatkan UI untuk messages:
- Message list
- Message detail
- Infinite scroll
- Message actions
- Payload viewer
- Timestamp formatting
- Loading states
```

## 4. Frontend Advanced Features

### 4.1 Real-time Updates
```prompt
Implementasikan real-time features:
- WebSocket connection
- Message subscription
- Status updates
- Error handling
- Reconnection logic
- Loading states
```

### 4.2 Search & Filter
```prompt
Buatkan search & filter components:
- Search input
- Filter dropdowns
- Date range picker
- Status filters
- Query builder
- URL sync
- Loading states
```

### 4.3 Responsive Design
```prompt
Implementasikan responsive design:
- Mobile layout
- Tablet layout
- Desktop layout
- Touch interactions
- Responsive tables
- Responsive forms
- Media queries
```

## 5. Monitoring & Dashboard

### 5.1 Prometheus Setup
```prompt
Setup Prometheus monitoring:
- Metrics collection
- Custom metrics
- Alert rules
- Service discovery
- Configuration
- Documentation
```

### 5.2 Dashboard UI
```prompt
Buatkan monitoring dashboard:
- System overview
- Tenant status
- Queue metrics
- Worker status
- Error rates
- Real-time updates
- Interactive charts
```

## 6. Testing & Optimization

### 6.1 Integration Tests
```prompt
Buatkan integration tests:
- API tests
- Database tests
- RabbitMQ tests
- End-to-end tests
- Test data setup
- CI integration
```

### 6.2 Performance Testing
```prompt
Implementasikan performance tests:
- Load testing
- Stress testing
- Benchmark tests
- Memory profiling
- CPU profiling
- Bottleneck analysis
```

## 7. Documentation & Deployment

### 7.1 API Documentation
```prompt
Generate API documentation:
- OpenAPI spec
- Endpoint descriptions
- Request/response examples
- Error codes
- Authentication
- Rate limits
```

### 7.2 Deployment
```prompt
Buatkan deployment scripts:
- Docker setup
- Docker Compose
- Environment config
- Health checks
- Backup strategy
- Rollback plan
```

## Tips Penggunaan Prompt

1. **Struktur Prompt**:
   - Mulai dengan konteks
   - Spesifik dengan requirements
   - Sertakan contoh jika ada
   - Minta penjelasan untuk bagian kompleks

2. **Iterasi**:
   - Review hasil
   - Minta perbaikan jika perlu
   - Tanyakan alternatif
   - Minta optimasi

3. **Best Practices**:
   - Gunakan bahasa yang jelas
   - Spesifik dengan teknologi
   - Minta dokumentasi
   - Verifikasi security

4. **Debugging**:
   - Minta analisis error
   - Saran perbaikan
   - Logging setup
   - Error handling

5. **Optimasi**:
   - Minta code review
   - Performance tips
   - Security check
   - Best practices 