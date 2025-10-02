#!/usr/bin/env node
/**
 * Shanraq.org Main Entry Point
 * Тенге-Веб Негізгі Кіру Нүктесі
 */

const { spawn } = require('child_process');
const path = require('path');

console.log('🚀 Shanraq.org іске қосылуда...');
console.log('=====================================================');

// Start the main server
const serverPath = path.join(__dirname, 'artjagy/server/main.js');
const server = spawn('node', [serverPath], {
  stdio: 'inherit',
  cwd: __dirname
});

server.on('error', (err) => {
  console.error('❌ Сервер қатесі:', err);
  process.exit(1);
});

server.on('close', (code) => {
  console.log(`🛑 Сервер жабылды, код: ${code}`);
  process.exit(code);
});

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.log('\n🛑 Shanraq.org серверін тоқтату...');
  server.kill('SIGINT');
});

process.on('SIGTERM', () => {
  console.log('\n🛑 Shanraq.org серверін тоқтату...');
  server.kill('SIGTERM');
});

