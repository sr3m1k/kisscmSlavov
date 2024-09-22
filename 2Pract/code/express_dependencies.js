const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

function getDependencies(packageName, visited = new Set()) {
    if (visited.has(packageName)) return;
    visited.add(packageName);

    const packageJsonPath = path.join('node_modules', packageName, 'package.json');

    if (!fs.existsSync(packageJsonPath)) return;

    const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf8'));
    const dependencies = packageJson.dependencies || {};

    for (const dep in dependencies) {
        console.log(`  "${packageName}" -> "${dep}";`);
        getDependencies(dep, visited);
    }
}

console.log('digraph ExpressDependencies {');
console.log('  rankdir="LR";');
getDependencies('express');
console.log('}');


const output = execSync('node express_dependencies.js').toString();
fs.writeFileSync('express_dependencies.dot', output);


