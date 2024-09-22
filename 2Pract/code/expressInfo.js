const fs = require('fs');
const path = require('path');
const express = require('express');

console.log(`Версия express: ${express.version}`);
console.log(`Путь установки express: ${require.resolve('express')}`);


const packageJsonPath = path.join(path.dirname(require.resolve('express')), 'package.json');
const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));

console.log('Зависимости express:');
console.log(JSON.stringify(packageJson.dependencies, null, 2));


const readmePath = path.join(path.dirname(require.resolve('express')), 'README.md');
if (fs.existsSync(readmePath)) {
    console.log('\nСодержимое README.md:');
    console.log(fs.readFileSync(readmePath, 'utf8'));
} else {
    console.log('\nФайл README.md не найден.');
}

