# ArvFinder Deployment Guide

This guide covers deploying ArvFinder to production environments.

## Prerequisites

- Docker and Docker Compose
- PostgreSQL database
- Stripe account with API keys
- Domain name and SSL certificate

## Environment Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Update the following variables in `.env`:
   ```bash
   # Production Database
   DATABASE_URL=postgres://username:password@host:5432/arvfinder?sslmode=require
   
   # Secure JWT Secret (generate with: openssl rand -base64 32)
   JWT_SECRET=your-secure-jwt-secret-here
   
   # Production Stripe Keys
   STRIPE_SECRET_KEY=sk_live_your_live_secret_key
   STRIPE_PUBLISHABLE_KEY=pk_live_your_live_publishable_key
   STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
   
   # Production Settings
   GIN_MODE=release
   PORT=8080
   
   # Frontend Configuration
   VITE_API_URL=https://your-domain.com/api/v1
   VITE_STRIPE_PUBLISHABLE_KEY=pk_live_your_live_publishable_key
   ```

## Stripe Setup

### 1. Create Products and Prices

Run the setup endpoint to create Stripe products and prices:

```bash
curl -X POST https://your-domain.com/api/v1/payments/setup-prices
```

This creates:
- **Professional Plan**: $29/month
- **Enterprise Plan**: $59/month

### 2. Configure Webhooks

1. Go to Stripe Dashboard → Webhooks
2. Add endpoint: `https://your-domain.com/api/v1/payments/webhook`
3. Select events:
   - `payment_intent.succeeded`
   - `invoice.payment_succeeded`
   - `customer.subscription.deleted`
   - `customer.subscription.updated`
4. Copy the webhook secret to your `.env` file

## Database Setup

### 1. Create Database

```sql
CREATE DATABASE arvfinder;
CREATE USER arvfinder WITH PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE arvfinder TO arvfinder;
```

### 2. Run Schema

```bash
psql -h your-host -U arvfinder -d arvfinder -f backend/database/schema.sql
```

## Docker Deployment

### 1. Production Docker Compose

Create `docker-compose.prod.yml`:

```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=${DATABASE_URL}
      - JWT_SECRET=${JWT_SECRET}
      - STRIPE_SECRET_KEY=${STRIPE_SECRET_KEY}
      - GIN_MODE=release
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend/frontend
      dockerfile: Dockerfile.prod
    ports:
      - "80:80"
      - "443:443"
    environment:
      - VITE_API_URL=${VITE_API_URL}
      - VITE_STRIPE_PUBLISHABLE_KEY=${VITE_STRIPE_PUBLISHABLE_KEY}
    restart: unless-stopped
    volumes:
      - ./ssl:/etc/ssl/certs
```

### 2. Create Production Frontend Dockerfile

`frontend/frontend/Dockerfile.prod`:

```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]
```

### 3. Deploy

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## Cloud Deployment (AWS/GCP/Azure)

### AWS ECS with Fargate

1. **Create Task Definition**:
   - Backend: 512 CPU, 1024 MB memory
   - Frontend: 256 CPU, 512 MB memory

2. **Set up Load Balancer**:
   - Application Load Balancer
   - SSL termination with ACM certificate
   - Health checks on `/health` endpoint

3. **Configure RDS**:
   - PostgreSQL instance
   - Multi-AZ for high availability
   - Automated backups

### Environment Variables in ECS

```json
{
  "name": "arvfinder-backend",
  "environment": [
    {"name": "DATABASE_URL", "value": "postgres://..."},
    {"name": "JWT_SECRET", "value": "..."},
    {"name": "STRIPE_SECRET_KEY", "value": "sk_live_..."},
    {"name": "GIN_MODE", "value": "release"}
  ]
}
```

## Monitoring and Logging

### 1. Health Checks

- Backend: `GET /health`
- Database connectivity check
- Stripe API connectivity check

### 2. Logging

Configure structured logging:

```go
// In main.go
if gin.Mode() == gin.ReleaseMode {
    gin.DefaultWriter = log.Writer()
    log.SetFormatter(&log.JSONFormatter{})
}
```

### 3. Metrics

Integrate with monitoring services:
- AWS CloudWatch
- Google Cloud Monitoring
- Datadog
- New Relic

## Security Considerations

### 1. Environment Variables

- Never commit `.env` files
- Use secrets management services in production
- Rotate JWT secrets regularly

### 2. Database Security

- Use SSL connections
- Implement connection pooling
- Regular security updates

### 3. API Security

- Rate limiting on API endpoints
- Input validation and sanitization
- CORS configuration for production domains

### 4. Stripe Security

- Validate webhook signatures
- Use HTTPS for all Stripe communications
- Monitor for unusual payment activities

## Backup Strategy

### 1. Database Backups

- Daily automated backups
- Point-in-time recovery
- Cross-region backup storage

### 2. Application Backups

- Container image versioning
- Configuration backup
- Code repository tags

## Scaling Considerations

### 1. Horizontal Scaling

- Load balancing across multiple backend instances
- CDN for frontend static assets
- Database read replicas

### 2. Performance Optimization

- Redis for session storage
- Database query optimization
- API response caching

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check DATABASE_URL format
   - Verify network connectivity
   - Check database credentials

2. **Stripe Webhook Failures**
   - Verify webhook secret
   - Check endpoint accessibility
   - Review Stripe dashboard logs

3. **CORS Issues**
   - Update allowed origins
   - Check protocol (HTTP vs HTTPS)
   - Verify headers configuration

### Logs Locations

- Backend logs: `/var/log/arvfinder/backend.log`
- Frontend logs: Browser console + server logs
- Database logs: PostgreSQL logs
- Stripe logs: Stripe Dashboard

## Support

For deployment assistance:
- Check GitHub Issues
- Review documentation
- Contact support team