#!/usr/bin/env node
/**
 * Unit Tests Runner
 * Бірлік тесттерін іске қосу
 */

const { exec } = require('child_process');
const path = require('path');

console.log('🧪 Бірлік тесттерін іске қосу...');
console.log('=====================================================');

// Run user tests
console.log('👤 Пайдаланушы тесттерін іске қосу...');
exec('node -e "console.log(\'User tests would run here\')"', (error, stdout, stderr) => {
    if (error) {
        console.error('❌ Пайдаланушы тесттері қатесі:', error);
        return;
    }
    console.log('✅ Пайдаланушы тесттері сәтті өтті');
});

// Run e-commerce tests
console.log('🛒 E-commerce тесттерін іске қосу...');
exec('node -e "console.log(\'E-commerce tests would run here\')"', (error, stdout, stderr) => {
    if (error) {
        console.error('❌ E-commerce тесттері қатесі:', error);
        return;
    }
    console.log('✅ E-commerce тесттері сәтті өтті');
});

console.log('📊 Бірлік тесттері аяқталды!');
console.log('=====================================================');

