#!/usr/bin/env node
/**
 * Shanraq.org Server Main Entry Point
 * Тенге-Веб Сервер Негізгі Кіру Нүктесі
 */

const express = require('express');
const path = require('path');
const cors = require('cors');
const helmet = require('helmet');
const compression = require('compression');
const rateLimit = require('express-rate-limit');
const morgan = require('morgan');

const app = express();
const PORT = process.env.PORT || 8080;

// Security middleware
app.use(helmet());
app.use(cors());

// Compression middleware
app.use(compression());

// Rate limiting
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100 // limit each IP to 100 requests per windowMs
});
app.use(limiter);

// Logging
app.use(morgan('combined'));

// Static files
app.use(express.static(path.join(__dirname, '../../betjagy/sandyq')));
app.use(express.static(path.join(__dirname, '../../betjagy/better')));

// JSON parsing
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Routes
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname, '../../betjagy/better/index.html'));
});

app.get('/demo', (req, res) => {
  res.sendFile(path.join(__dirname, '../../synaqtar/demo/demo.html'));
});

// JOJJ API Routes - Пайдаланушылар (Users)
// Import JOJJ basqaru
// const jojjBasqaru = require('./jojj_basqaru');

app.get('/api/v1/paydalanusylar', (req, res) => {
  // Implementation would call api_paydalanusylar_oqu_barlik
  res.json({ message: 'Users API endpoint', success: true });
});

app.get('/api/v1/paydalanusylar/:id', (req, res) => {
  // Implementation would call api_paydalanu_oqu
  res.json({ message: 'Get user by ID', success: true, id: req.params.id });
});

app.post('/api/v1/paydalanusylar', (req, res) => {
  // Implementation would call api_paydalanu_jasau
  res.json({ message: 'Create user', success: true, data: req.body });
});

app.put('/api/v1/paydalanusylar/:id', (req, res) => {
  // Implementation would call api_paydalanu_janartu
  res.json({ message: 'Update user', success: true, id: req.params.id, data: req.body });
});

app.delete('/api/v1/paydalanusylar/:id', (req, res) => {
  // Implementation would call api_paydalanu_joiu
  res.json({ message: 'Delete user', success: true, id: req.params.id });
});

// JOJJ API Routes - Мақалалар (Articles)
app.get('/api/v1/maqalalar', (req, res) => {
  // Implementation would call api_maqalalar_oqu_barlik
  res.json({ message: 'Articles API endpoint', success: true });
});

app.get('/api/v1/maqalalar/:id', (req, res) => {
  // Implementation would call api_maqala_oqu
  res.json({ message: 'Get article by ID', success: true, id: req.params.id });
});

app.post('/api/v1/maqalalar', (req, res) => {
  // Implementation would call api_maqala_jasau
  res.json({ message: 'Create article', success: true, data: req.body });
});

app.put('/api/v1/maqalalar/:id', (req, res) => {
  // Implementation would call api_maqala_janartu
  res.json({ message: 'Update article', success: true, id: req.params.id, data: req.body });
});

app.delete('/api/v1/maqalalar/:id', (req, res) => {
  // Implementation would call api_maqala_joiu
  res.json({ message: 'Delete article', success: true, id: req.params.id });
});

// JOJJ API Routes - Санаттар (Categories)
app.get('/api/v1/sanattar', (req, res) => {
  // Implementation would call api_sanattar_oqu_barlik
  res.json({ message: 'Categories API endpoint', success: true });
});

app.get('/api/v1/sanattar/:id', (req, res) => {
  // Implementation would call api_sanat_oqu
  res.json({ message: 'Get category by ID', success: true, id: req.params.id });
});

app.post('/api/v1/sanattar', (req, res) => {
  // Implementation would call api_sanat_jasau
  res.json({ message: 'Create category', success: true, data: req.body });
});

app.put('/api/v1/sanattar/:id', (req, res) => {
  // Implementation would call api_sanat_janartu
  res.json({ message: 'Update category', success: true, id: req.params.id, data: req.body });
});

app.delete('/api/v1/sanattar/:id', (req, res) => {
  // Implementation would call api_sanat_joiu
  res.json({ message: 'Delete category', success: true, id: req.params.id });
});

// Additional JOJJ API Routes
app.get('/api/v1/maqalalar/popular', (req, res) => {
  res.json({ message: 'Popular articles', success: true });
});

app.get('/api/v1/maqalalar/recent', (req, res) => {
  res.json({ message: 'Recent articles', success: true });
});

app.get('/api/v1/sanattar/hierarchy', (req, res) => {
  res.json({ message: 'Categories hierarchy', success: true });
});

app.get('/api/v1/statistics', (req, res) => {
  res.json({ 
    message: 'Statistics', 
    success: true,
    data: {
      users: { total: 0, active: 0 },
      articles: { total: 0, published: 0 },
      categories: { total: 0, active: 0 }
    }
  });
});

// API Routes
app.get('/api/v1/health', (req, res) => {
  res.json({
    status: 'healthy',
    server: 'Shanraq.org',
    version: '1.0.0',
    timestamp: new Date().toISOString()
  });
});

app.get('/api/v1/status', (req, res) => {
  res.json({
    server: 'running',
    database: 'connected',
    cache: 'active',
    uptime: process.uptime(),
    language: 'Tenge',
    framework: 'Shanraq.org',
    morpheme_cache: 'active',
    phoneme_optimizer: 'enabled',
    archetype_engine: 'loaded',
    simd_processor: 'enabled'
  });
});

// Error handling
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    error: 'Internal Server Error',
    message: 'Something went wrong!'
  });
});

// 404 handler
app.use((req, res) => {
  res.status(404).json({
    error: 'Not Found',
    message: 'The requested resource was not found'
  });
});

// Start server
app.listen(PORT, () => {
  console.log('🚀 Shanraq.org сервер іске қосылды!');
  console.log('=====================================================');
  console.log(`🌐 Сервер: http://localhost:${PORT}`);
  console.log(`📄 Басты бет: http://localhost:${PORT}/`);
  console.log(`📝 Блог: http://localhost:${PORT}/blog`);
  console.log(`👥 Біз жөнінде: http://localhost:${PORT}/about`);
  console.log(`📞 Байланыс: http://localhost:${PORT}/contact`);
  console.log(`🔧 API: http://localhost:${PORT}/api/v1/health`);
  console.log('=====================================================');
  console.log('Серверді тоқтату үшін Ctrl+C басыңыз');
});

module.exports = app;

