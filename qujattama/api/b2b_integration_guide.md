# B2B Integration Guide
# Руководство по B2B Интеграции
# Complete guide for enterprise integration with Shanraq.org

## Overview / Обзор

This guide provides comprehensive instructions for integrating with Shanraq.org's B2B fintech platform. Our platform offers enterprise-grade financial infrastructure through modern APIs and SDKs.

## 🚀 Quick Start / Быстрый Старт

### 1. Get API Credentials / Получение API Ключей

```bash
# Register for B2B access
curl -X POST https://api.shanraq.org/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Your Company",
    "contact_email": "contact@yourcompany.com",
    "business_type": "fintech",
    "expected_volume": "100000"
  }'
```

### 2. Install SDK / Установка SDK

```bash
# Node.js SDK
npm install @shanraq/sdk

# Python SDK
pip install shanraq-sdk

# Java SDK
<dependency>
  <groupId>org.shanraq</groupId>
  <artifactId>shanraq-sdk</artifactId>
  <version>1.0.0</version>
</dependency>
```

### 3. Initialize Client / Инициализация Клиента

```javascript
// JavaScript/Node.js
const ShanraqClient = require('@shanraq/sdk');

const client = new ShanraqClient({
  apiKey: 'your-api-key',
  secretKey: 'your-secret-key',
  environment: 'sandbox', // or 'production'
  region: 'kz' // Kazakhstan region
});
```

```python
# Python
from shanraq import ShanraqClient

client = ShanraqClient(
    api_key='your-api-key',
    secret_key='your-secret-key',
    environment='sandbox',
    region='kz'
)
```

```java
// Java
import org.shanraq.ShanraqClient;

ShanraqClient client = new ShanraqClient.Builder()
    .apiKey("your-api-key")
    .secretKey("your-secret-key")
    .environment("sandbox")
    .region("kz")
    .build();
```

## 💳 Payment Processing / Обработка Платежей

### Create Payment / Создание Платежа

```javascript
// Create a payment
const payment = await client.payments.create({
  amount: 10000, // 100.00 KZT (amount in tiyn)
  currency: 'KZT',
  description: 'Payment for services',
  customer: {
    id: 'customer_123',
    email: 'customer@example.com',
    phone: '+7 777 123 4567'
  },
  metadata: {
    order_id: 'order_123',
    product: 'Premium Service'
  }
});

console.log('Payment created:', payment.id);
```

### Process P2P Transfer / Обработка P2P Перевода

```javascript
// P2P transfer
const transfer = await client.transfers.create({
  from_account: 'account_123',
  to_account: 'account_456',
  amount: 5000, // 50.00 KZT
  currency: 'KZT',
  description: 'Money transfer',
  reference: 'transfer_ref_123'
});

console.log('Transfer processed:', transfer.id);
```

### QR Payment / QR Платеж

```javascript
// Generate QR code for payment
const qrPayment = await client.qr.create({
  amount: 2500, // 25.00 KZT
  currency: 'KZT',
  description: 'QR payment',
  expiry_minutes: 30
});

console.log('QR Code:', qrPayment.qr_code);
console.log('QR Data:', qrPayment.qr_data);
```

## 🏦 Banking Operations / Банковские Операции

### Account Management / Управление Счетами

```javascript
// Create account
const account = await client.accounts.create({
  account_type: 'checking',
  currency: 'KZT',
  owner_id: 'user_123',
  initial_balance: 0
});

// Get account balance
const balance = await client.accounts.getBalance(account.id);
console.log('Account balance:', balance.amount);

// Get account history
const history = await client.accounts.getHistory(account.id, {
  start_date: '2025-01-01',
  end_date: '2025-01-31',
  limit: 100
});
```

### Transaction Processing / Обработка Транзакций

```javascript
// Get transaction details
const transaction = await client.transactions.get('txn_123');

// List transactions
const transactions = await client.transactions.list({
  account_id: 'account_123',
  status: 'completed',
  limit: 50,
  offset: 0
});

// Search transactions
const searchResults = await client.transactions.search({
  query: 'payment',
  date_range: {
    start: '2025-01-01',
    end: '2025-01-31'
  }
});
```

## 🔒 Security & Compliance / Безопасность и Соответствие

### Authentication / Аутентификация

```javascript
// API key authentication (recommended for server-to-server)
const client = new ShanraqClient({
  apiKey: 'your-api-key',
  secretKey: 'your-secret-key'
});

// JWT authentication (for user-facing applications)
const jwtToken = await client.auth.getJWTToken({
  user_id: 'user_123',
  permissions: ['payments:create', 'accounts:read']
});
```

### Webhook Security / Безопасность Webhook

```javascript
// Verify webhook signature
app.post('/webhook', (req, res) => {
  const signature = req.headers['x-shanraq-signature'];
  const payload = req.body;
  
  const isValid = client.webhooks.verifySignature(
    payload,
    signature,
    'your-webhook-secret'
  );
  
  if (isValid) {
    // Process webhook
    console.log('Webhook verified:', payload);
  } else {
    res.status(400).send('Invalid signature');
  }
});
```

## 📊 Monitoring & Analytics / Мониторинг и Аналитика

### Real-time Monitoring / Мониторинг в Реальном Времени

```javascript
// Get platform status
const status = await client.monitoring.getStatus();
console.log('Platform status:', status.overall);

// Get performance metrics
const metrics = await client.monitoring.getMetrics({
  metric: 'response_time',
  period: '1h',
  granularity: '1m'
});

// Get error rates
const errors = await client.monitoring.getErrors({
  start_time: '2025-01-01T00:00:00Z',
  end_time: '2025-01-01T23:59:59Z'
});
```

### Analytics Dashboard / Аналитическая Панель

```javascript
// Get transaction analytics
const analytics = await client.analytics.getTransactions({
  account_id: 'account_123',
  period: '30d',
  group_by: 'day'
});

// Get revenue analytics
const revenue = await client.analytics.getRevenue({
  start_date: '2025-01-01',
  end_date: '2025-01-31',
  currency: 'KZT'
});
```

## 🔄 Event Handling / Обработка Событий

### Webhook Configuration / Конфигурация Webhook

```javascript
// Create webhook endpoint
const webhook = await client.webhooks.create({
  url: 'https://your-app.com/webhook',
  events: [
    'payment.completed',
    'payment.failed',
    'transfer.completed',
    'account.created'
  ],
  secret: 'your-webhook-secret'
});

// List webhooks
const webhooks = await client.webhooks.list();

// Update webhook
await client.webhooks.update(webhook.id, {
  events: ['payment.completed', 'transfer.completed']
});
```

### Event Processing / Обработка Событий

```javascript
// Process webhook events
app.post('/webhook', async (req, res) => {
  const event = req.body;
  
  switch (event.type) {
    case 'payment.completed':
      await handlePaymentCompleted(event.data);
      break;
    case 'transfer.completed':
      await handleTransferCompleted(event.data);
      break;
    case 'account.created':
      await handleAccountCreated(event.data);
      break;
  }
  
  res.status(200).send('OK');
});
```

## 🧪 Testing & Development / Тестирование и Разработка

### Sandbox Environment / Песочница

```javascript
// Use sandbox for testing
const sandboxClient = new ShanraqClient({
  apiKey: 'sandbox-api-key',
  environment: 'sandbox'
});

// Create test payment
const testPayment = await sandboxClient.payments.create({
  amount: 1000, // 10.00 KZT
  currency: 'KZT',
  description: 'Test payment',
  test_mode: true
});
```

### Test Data / Тестовые Данные

```javascript
// Create test accounts
const testAccount1 = await client.accounts.create({
  account_type: 'test',
  currency: 'KZT',
  initial_balance: 100000 // 1000.00 KZT
});

const testAccount2 = await client.accounts.create({
  account_type: 'test',
  currency: 'KZT',
  initial_balance: 0
});

// Test transfer between accounts
const testTransfer = await client.transfers.create({
  from_account: testAccount1.id,
  to_account: testAccount2.id,
  amount: 5000, // 50.00 KZT
  currency: 'KZT',
  description: 'Test transfer'
});
```

## 📚 API Reference / Справочник API

### Authentication / Аутентификация

All API requests require authentication using API keys or JWT tokens.

```bash
# API Key Authentication
curl -H "Authorization: Bearer your-api-key" \
     -H "X-API-Secret: your-secret-key" \
     https://api.shanraq.org/v1/payments

# JWT Authentication
curl -H "Authorization: Bearer jwt-token" \
     https://api.shanraq.org/v1/accounts
```

### Rate Limits / Лимиты Скорости

- **Standard Plan**: 1,000 requests per minute
- **Professional Plan**: 10,000 requests per minute
- **Enterprise Plan**: 100,000 requests per minute

### Error Handling / Обработка Ошибок

```javascript
try {
  const payment = await client.payments.create(paymentData);
} catch (error) {
  if (error.code === 'INSUFFICIENT_FUNDS') {
    console.log('Insufficient funds in account');
  } else if (error.code === 'INVALID_CURRENCY') {
    console.log('Invalid currency specified');
  } else {
    console.log('Payment failed:', error.message);
  }
}
```

## 🚀 Production Deployment / Продакшн Развертывание

### Environment Configuration / Конфигурация Окружения

```javascript
// Production configuration
const productionClient = new ShanraqClient({
  apiKey: process.env.SHANRAQ_API_KEY,
  secretKey: process.env.SHANRAQ_SECRET_KEY,
  environment: 'production',
  region: 'kz',
  timeout: 30000, // 30 seconds
  retries: 3
});
```

### Security Best Practices / Лучшие Практики Безопасности

1. **Store API keys securely** - Use environment variables or secure key management
2. **Use HTTPS** - Always use HTTPS for API requests
3. **Validate webhooks** - Always verify webhook signatures
4. **Implement rate limiting** - Respect API rate limits
5. **Monitor usage** - Track API usage and costs

## 📞 Support & Resources / Поддержка и Ресурсы

### Documentation / Документация
- **API Reference**: https://docs.shanraq.org/api
- **SDK Documentation**: https://docs.shanraq.org/sdk
- **Integration Examples**: https://docs.shanraq.org/examples

### Support Channels / Каналы Поддержки
- **Email**: support@shanraq.org
- **Slack**: https://shanraq.slack.com
- **Phone**: +7 727 123 4567
- **Office Hours**: 9:00 AM - 6:00 PM (Almaty time)

### Community / Сообщество
- **GitHub**: https://github.com/shanraq
- **Discord**: https://discord.gg/shanraq
- **Telegram**: https://t.me/shanraq_community

---

**Ready to integrate?** Start with our sandbox environment and build your financial services on Shanraq.org's robust infrastructure.

**Готовы к интеграции?** Начните с нашей песочницы и создавайте финансовые сервисы на надежной инфраструктуре Shanraq.org.
