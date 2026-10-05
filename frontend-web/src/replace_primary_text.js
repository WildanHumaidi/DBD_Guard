const fs = require('fs');
const path = require('path');

function walk(dir) {
    let results = [];
    const list = fs.readdirSync(dir);
    list.forEach((file) => {
        file = path.join(dir, file);
        const stat = fs.statSync(file);
        if (stat && stat.isDirectory()) {
            results = results.concat(walk(file));
        } else {
            if (file.endsWith('.tsx') || file.endsWith('.ts')) {
                results.push(file);
            }
        }
    });
    return results;
}

const files = walk('e:/Projects/DBD_Guard/frontend-web/src');

files.forEach(file => {
    let content = fs.readFileSync(file, 'utf8');
    let changed = false;

    // Find all class strings
    // E.g. className="..." or className={`...`}
    // A simple regex to replace 'text-white' with 'text-primary-foreground' 
    // ONLY IF the same string contains 'bg-primary' or 'bg-primary-deep'.
    
    // We can do a simpler replace: split by lines, if line contains 'bg-primary' and 'text-white', replace 'text-white' with 'text-primary-foreground'
    let lines = content.split('\n');
    for(let i=0; i<lines.length; i++) {
        let line = lines[i];
        if ((line.includes('bg-primary') || line.includes('bg-primary-deep') || line.includes('bg-primary-tint') || line.includes('from-teal-500') || line.includes('bg-emerald-600') || line.includes('bg-teal-600') || line.includes('bg-emerald-500')) && line.includes('text-white')) {
            // Replace text-white with text-primary-foreground
            lines[i] = line.replace(/text-white/g, 'text-primary-foreground');
            changed = true;
        }
    }

    if (changed) {
        fs.writeFileSync(file, lines.join('\n'), 'utf8');
        console.log(`Updated ${file}`);
    }
});
