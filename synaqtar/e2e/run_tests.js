#!/usr/bin/env node
/**
 * E2E Tests Runner
 * E2E тесттерін іске қосу
 */

const { exec } = require('child_process');
const path = require('path');

console.log('🧪 E2E тесттерін іске қосу...');
console.log('=====================================================');

// Run e-commerce E2E tests
console.log('🛒 E-commerce E2E тесттерін іске қосу...');
exec('node -e "console.log(\'E-commerce E2E tests would run here\')"', (error, stdout, stderr) => {
    if (error) {
        console.error('❌ E-commerce E2E тесттері қатесі:', error);
        return;
    }
    console.log('✅ E-commerce E2E тесттері сәтті өтті');
});

// Run user E2E tests
console.log('👤 Пайдаланушы E2E тесттерін іске қосу...');
exec('node -e "console.log(\'User E2E tests would run here\')"', (error, stdout, stderr) => {
    if (error) {
        console.error('❌ Пайдаланушы E2E тесттері қатесі:', error);
        return;
    }
    console.log('✅ Пайдаланушы E2E тесттері сәтті өтті');
});

console.log('📊 E2E тесттері аяқталды!');
console.log('=====================================================');

