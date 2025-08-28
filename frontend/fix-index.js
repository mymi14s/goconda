import fs from 'fs';
import path from 'path';

const filesPath = path.resolve('../backend/static/admin');

// 1) Recursively process all .html and .js files under filesPath
function getAllTargetFiles(dir) {
  const results = [];
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  for (const entry of entries) {
    const full = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      results.push(...getAllTargetFiles(full));
    } else if (entry.isFile()) {
      const ext = path.extname(entry.name).toLowerCase();
      if (ext === '.html' || ext === '.js') {
        results.push(full);
      }
    }
  }
  return results;
}

function applyReplacements(text) {
  return text
    .replace(/src="\.\//g, 'src="/static/admin/')
    .replace(/href="\.\//g, 'href="/static/admin/')
    .replace('/src/', '/static/admin/');
}

// Process files
const targets = getAllTargetFiles(filesPath);
let changedCount = 0;

for (const file of targets) {
  const original = fs.readFileSync(file, 'utf8');
  const updated = applyReplacements(original);
  if (updated !== original) {
    fs.writeFileSync(file, updated, 'utf8');
    changedCount++;
  }
}

console.log(`Scanned ${targets.length} files. Updated ${changedCount} file(s).`);

// 2) Continue with the rest of your original code

const srcPath = path.resolve('../backend/static/admin/index.html');
const destPath = path.resolve('../backend/views/admin/index.html');

// Read the file
let content = fs.readFileSync(srcPath, 'utf8'); 

fs.mkdirSync(path.dirname(destPath), { recursive: true });

// Write modified file
fs.writeFileSync(destPath, content, 'utf8');


console.log(`index.html patched and copied to ${destPath}`);
