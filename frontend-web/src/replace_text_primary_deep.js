const fs = require('fs');
const path = require('path');

const walkSync = (dir, filelist = []) => {
  fs.readdirSync(dir).forEach(file => {
    const filePath = path.join(dir, file);
    if (fs.statSync(filePath).isDirectory()) {
      if (file !== 'node_modules' && file !== '.next') {
        filelist = walkSync(filePath, filelist);
      }
    } else {
      if (filePath.endsWith('.tsx') || filePath.endsWith('.ts')) {
        filelist.push(filePath);
      }
    }
  });
  return filelist;
};

const files = walkSync(path.join(__dirname, 'app'));

let changedFiles = 0;
files.forEach(file => {
  let content = fs.readFileSync(file, 'utf8');
  let newContent = content;
  
  // Replace text-primary-deep with text-primary (but carefully avoiding border-primary-deep or bg-primary-deep)
  // We can just regex replace `text-primary-deep` with `text-primary`.
  newContent = newContent.replace(/text-primary-deep/g, 'text-primary');

  if (content !== newContent) {
    fs.writeFileSync(file, newContent, 'utf8');
    console.log(`Updated ${file}`);
    changedFiles++;
  }
});

console.log(`Total files changed: ${changedFiles}`);
