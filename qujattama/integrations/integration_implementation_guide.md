# Integration Implementation Guide

## Overview

This guide provides step-by-step instructions for implementing and configuring the Shanraq.org fintech integration system. It covers all aspects from initial setup to production deployment.

## Prerequisites

### System Requirements
- **Operating System**: Linux (Ubuntu 20.04+ recommended)
- **Memory**: 8GB RAM minimum, 16GB recommended
- **Storage**: 100GB SSD minimum
- **Network**: Stable internet connection
- **Dependencies**: Node.js 18+, Python 3.9+, Docker

### Required Services
- **Database**: PostgreSQL 13+ or MySQL 8+
- **Cache**: Redis 6+
- **Message Queue**: Kafka 2.8+ or RabbitMQ 3.8+
- **Load Balancer**: Nginx 1.18+
- **SSL Certificate**: Valid SSL certificate

## Installation Steps

### 1. Environment Setup

#### Clone Repository
```bash
git clone https://github.com/shanraq-org/shanraq.org.git
cd shanraq.org
```

#### Install Dependencies
```bash
# Install Node.js dependencies
npm install

# Install Python dependencies
pip install -r requirements.txt

# Install system dependencies
sudo apt-get update
sudo apt-get install -y postgresql redis-server nginx
```

#### Configure Environment Variables
```bash
# Copy environment template
cp .env.example .env

# Edit environment variables
nano .env
```

#### Environment Configuration
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=shanraq_fintech
DB_USER=shanraq_user
DB_PASSWORD=secure_password

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis_password

# Message Queue Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC_PREFIX=shanraq

# Security Configuration
JWT_SECRET=your_jwt_secret
API_KEY_SECRET=your_api_key_secret
ENCRYPTION_KEY=your_encryption_key

# Payment Adapters
VISA_API_KEY=your_visa_api_key
VISA_API_SECRET=your_visa_api_secret
VISA_MERCHANT_ID=your_visa_merchant_id

MASTERCARD_API_KEY=your_mastercard_api_key
MASTERCARD_API_SECRET=your_mastercard_api_secret
MASTERCARD_MERCHANT_ID=your_mastercard_merchant_id

# Webhook Configuration
WEBHOOK_SECRET=your_webhook_secret
WEBHOOK_RETRY_ATTEMPTS=5
WEBHOOK_RETRY_DELAY=1000

# KYC/AML Configuration
KYC_VERIFICATION_PROVIDER=internal
AML_MONITORING_MODE=real_time
AML_RISK_THRESHOLD=70
```

### 2. Database Setup

#### PostgreSQL Configuration
```sql
-- Create database
CREATE DATABASE shanraq_fintech;

-- Create user
CREATE USER shanraq_user WITH PASSWORD 'secure_password';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE shanraq_fintech TO shanraq_user;

-- Connect to database
\c shanraq_fintech

-- Create tables
\i sql/schema.sql
```

#### Database Schema
```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    account_number VARCHAR(20) UNIQUE NOT NULL,
    account_type VARCHAR(20) NOT NULL,
    balance DECIMAL(15,2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'KZT',
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_account_id UUID REFERENCES accounts(id),
    to_account_id UUID REFERENCES accounts(id),
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- KYC Profiles table
CREATE TABLE kyc_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    personal_info JSONB NOT NULL,
    identity_documents JSONB[],
    address_info JSONB NOT NULL,
    contact_info JSONB NOT NULL,
    risk_score INTEGER DEFAULT 0,
    verification_level VARCHAR(20) DEFAULT 'basic',
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- AML Alerts table
CREATE TABLE aml_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id),
    user_id UUID REFERENCES users(id),
    risk_score INTEGER NOT NULL,
    triggered_rules JSONB[],
    alert_type VARCHAR(20) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'open',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,
    resolution_notes TEXT
);

-- Webhooks table
CREATE TABLE webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    partner_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    webhook_url VARCHAR(500) NOT NULL,
    secret VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Webhook Deliveries table
CREATE TABLE webhook_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    webhook_id UUID REFERENCES webhooks(id),
    event_id UUID NOT NULL,
    url VARCHAR(500) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    attempt INTEGER DEFAULT 1,
    max_attempts INTEGER DEFAULT 5,
    response_code INTEGER,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    delivered_at TIMESTAMP,
    failed_at TIMESTAMP
);
```

### 3. Service Configuration

#### Redis Configuration
```bash
# Edit Redis configuration
sudo nano /etc/redis/redis.conf

# Set password
requirepass redis_password

# Restart Redis
sudo systemctl restart redis
```

#### Kafka Configuration
```bash
# Download Kafka
wget https://downloads.apache.org/kafka/2.8.0/kafka_2.13-2.8.0.tgz
tar -xzf kafka_2.13-2.8.0.tgz
cd kafka_2.13-2.8.0

# Start Zookeeper
bin/zookeeper-server-start.sh config/zookeeper.properties &

# Start Kafka
bin/kafka-server-start.sh config/server.properties &

# Create topics
bin/kafka-topics.sh --create --topic shanraq-transactions --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic shanraq-payments --bootstrap-server localhost:9092
bin/kafka-topics.sh --create --topic shanraq-webhooks --bootstrap-server localhost:9092
```

#### Nginx Configuration
```nginx
# /etc/nginx/sites-available/shanraq
server {
    listen 80;
    server_name api.shanraq.org;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.shanraq.org;

    ssl_certificate /etc/ssl/certs/shanraq.crt;
    ssl_certificate_key /etc/ssl/private/shanraq.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### 4. Application Setup

#### Initialize Services
```bash
# Start database migrations
npm run migrate

# Start Redis
redis-server

# Start Kafka
kafka-server-start.sh config/server.properties

# Start application
npm start
```

#### Verify Installation
```bash
# Check database connection
npm run db:test

# Check Redis connection
npm run redis:test

# Check Kafka connection
npm run kafka:test

# Run health check
curl https://api.shanraq.org/health
```

## Configuration

### 1. Payment Adapters

#### Visa Configuration
```javascript
const visaConfig = {
  merchantId: process.env.VISA_MERCHANT_ID,
  terminalId: process.env.VISA_TERMINAL_ID,
  apiKey: process.env.VISA_API_KEY,
  apiSecret: process.env.VISA_API_SECRET,
  apiEndpoint: process.env.VISA_API_ENDPOINT,
  environment: process.env.NODE_ENV === 'production' ? 'production' : 'sandbox'
};
```

#### Mastercard Configuration
```javascript
const mastercardConfig = {
  merchantId: process.env.MASTERCARD_MERCHANT_ID,
  terminalId: process.env.MASTERCARD_TERMINAL_ID,
  apiKey: process.env.MASTERCARD_API_KEY,
  apiSecret: process.env.MASTERCARD_API_SECRET,
  apiEndpoint: process.env.MASTERCARD_API_ENDPOINT,
  environment: process.env.NODE_ENV === 'production' ? 'production' : 'sandbox'
};
```

#### QR Code Configuration
```javascript
const qrConfig = {
  provider: process.env.QR_PROVIDER,
  merchantId: process.env.QR_MERCHANT_ID,
  apiKey: process.env.QR_API_KEY,
  apiSecret: process.env.QR_API_SECRET,
  apiEndpoint: process.env.QR_API_ENDPOINT
};
```

#### KaspiPay Configuration
```javascript
const kaspiPayConfig = {
  merchantId: process.env.KASPIPAY_MERCHANT_ID,
  terminalId: process.env.KASPIPAY_TERMINAL_ID,
  apiKey: process.env.KASPIPAY_API_KEY,
  apiSecret: process.env.KASPIPAY_API_SECRET,
  apiEndpoint: process.env.KASPIPAY_API_ENDPOINT
};
```

### 2. Webhook Configuration

#### Webhook Endpoints
```javascript
const webhookEndpoints = {
  'transaction.created': 'https://partner.example.com/webhooks/transaction-created',
  'transaction.executed': 'https://partner.example.com/webhooks/transaction-executed',
  'payment.completed': 'https://partner.example.com/webhooks/payment-completed',
  'p2p_transfer.completed': 'https://partner.example.com/webhooks/p2p-transfer-completed'
};
```

#### Webhook Security
```javascript
const webhookSecurity = {
  secret: process.env.WEBHOOK_SECRET,
  signatureHeader: 'X-Webhook-Signature',
  retryAttempts: 5,
  retryDelay: 1000,
  timeout: 30000
};
```

### 3. KYC/AML Configuration

#### KYC Settings
```javascript
const kycSettings = {
  verificationProvider: 'internal',
  documentStorage: 'encrypted',
  biometricVerification: true,
  livenessDetection: true,
  riskThresholds: {
    low: 30,
    medium: 60,
    high: 80,
    critical: 90
  }
};
```

#### AML Settings
```javascript
const amlSettings = {
  monitoringMode: 'real_time',
  ruleEngine: 'advanced',
  machineLearning: true,
  patternDetection: true,
  riskThresholds: {
    low: 30,
    medium: 60,
    high: 80,
    critical: 90
  }
};
```

## Testing

### 1. Unit Testing
```bash
# Run unit tests
npm test

# Run with coverage
npm run test:coverage

# Run specific test file
npm test -- --grep "payment adapters"
```

### 2. Integration Testing
```bash
# Run integration tests
npm run test:integration

# Run with database
npm run test:integration:db

# Run with external services
npm run test:integration:external
```

### 3. End-to-End Testing
```bash
# Run E2E tests
npm run test:e2e

# Run with browser
npm run test:e2e:browser

# Run with mobile
npm run test:e2e:mobile
```

### 4. Load Testing
```bash
# Run load tests
npm run test:load

# Run stress tests
npm run test:stress

# Run performance tests
npm run test:performance
```

## Deployment

### 1. Production Deployment

#### Docker Configuration
```dockerfile
# Dockerfile
FROM node:18-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --only=production

COPY . .

EXPOSE 3000

CMD ["npm", "start"]
```

#### Docker Compose
```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - DB_HOST=postgres
      - REDIS_HOST=redis
      - KAFKA_BROKERS=kafka:9092
    depends_on:
      - postgres
      - redis
      - kafka

  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=shanraq_fintech
      - POSTGRES_USER=shanraq_user
      - POSTGRES_PASSWORD=secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:6-alpine
    command: redis-server --requirepass redis_password

  kafka:
    image: confluentinc/cp-kafka:latest
    environment:
      - KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
      - KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092
      - KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1

  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      - ZOOKEEPER_CLIENT_PORT=2181
      - ZOOKEEPER_TICK_TIME=2000

volumes:
  postgres_data:
```

#### Kubernetes Configuration
```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: shanraq-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: shanraq-app
  template:
    metadata:
      labels:
        app: shanraq-app
    spec:
      containers:
      - name: shanraq-app
        image: shanraq/app:latest
        ports:
        - containerPort: 3000
        env:
        - name: NODE_ENV
          value: "production"
        - name: DB_HOST
          value: "postgres-service"
        - name: REDIS_HOST
          value: "redis-service"
        - name: KAFKA_BROKERS
          value: "kafka-service:9092"
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "500m"
```

### 2. Monitoring Setup

#### Prometheus Configuration
```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'shanraq-app'
    static_configs:
      - targets: ['localhost:3000']
    metrics_path: '/metrics'
    scrape_interval: 5s
```

#### Grafana Dashboard
```json
{
  "dashboard": {
    "title": "Shanraq Fintech Dashboard",
    "panels": [
      {
        "title": "Transaction Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(transactions_total[5m])",
            "legendFormat": "Transactions/sec"
          }
        ]
      },
      {
        "title": "Error Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(errors_total[5m])",
            "legendFormat": "Errors/sec"
          }
        ]
      }
    ]
  }
}
```

## Maintenance

### 1. Regular Maintenance

#### Database Maintenance
```sql
-- Analyze tables
ANALYZE;

-- Vacuum tables
VACUUM ANALYZE;

-- Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

#### Log Rotation
```bash
# Configure logrotate
sudo nano /etc/logrotate.d/shanraq

# Log rotation configuration
/var/log/shanraq/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 shanraq shanraq
    postrotate
        systemctl reload shanraq
    endscript
}
```

### 2. Backup Strategy

#### Database Backup
```bash
# Create backup script
#!/bin/bash
BACKUP_DIR="/backups/shanraq"
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -h localhost -U shanraq_user -d shanraq_fintech > $BACKUP_DIR/backup_$DATE.sql
gzip $BACKUP_DIR/backup_$DATE.sql
```

#### Configuration Backup
```bash
# Backup configuration files
tar -czf /backups/shanraq/config_$(date +%Y%m%d_%H%M%S).tar.gz \
  /etc/nginx/sites-available/shanraq \
  /etc/redis/redis.conf \
  /etc/kafka/server.properties \
  /home/shanraq/.env
```

### 3. Security Updates

#### System Updates
```bash
# Update system packages
sudo apt update
sudo apt upgrade -y

# Update Node.js
sudo npm install -g n
sudo n stable

# Update Python
sudo apt install python3.9 python3.9-pip
```

#### Application Updates
```bash
# Update application dependencies
npm update

# Update security patches
npm audit fix

# Update Docker images
docker pull shanraq/app:latest
```

## Troubleshooting

### 1. Common Issues

#### Database Connection Issues
```bash
# Check database status
sudo systemctl status postgresql

# Check connection
psql -h localhost -U shanraq_user -d shanraq_fintech

# Check logs
sudo tail -f /var/log/postgresql/postgresql-13-main.log
```

#### Redis Connection Issues
```bash
# Check Redis status
sudo systemctl status redis

# Check connection
redis-cli ping

# Check logs
sudo tail -f /var/log/redis/redis-server.log
```

#### Kafka Connection Issues
```bash
# Check Kafka status
sudo systemctl status kafka

# Check topics
kafka-topics.sh --list --bootstrap-server localhost:9092

# Check logs
tail -f /opt/kafka/logs/server.log
```

### 2. Performance Issues

#### Database Performance
```sql
-- Check slow queries
SELECT query, mean_time, calls
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

#### Application Performance
```bash
# Check memory usage
free -h

# Check CPU usage
top

# Check disk usage
df -h

# Check network connections
netstat -tulpn
```

### 3. Security Issues

#### SSL Certificate Issues
```bash
# Check certificate validity
openssl x509 -in /etc/ssl/certs/shanraq.crt -text -noout

# Check certificate expiration
openssl x509 -in /etc/ssl/certs/shanraq.crt -dates -noout

# Renew certificate
certbot renew
```

#### Authentication Issues
```bash
# Check JWT tokens
jwt decode <token>

# Check API keys
curl -H "Authorization: Bearer <token>" https://api.shanraq.org/health

# Check webhook signatures
openssl dgst -sha256 -hmac <secret> <payload>
```

## Conclusion

This implementation guide provides comprehensive instructions for setting up and maintaining the Shanraq.org fintech integration system. Follow these steps carefully to ensure a successful deployment and ongoing operation.

For additional support, refer to the documentation or contact the development team.
