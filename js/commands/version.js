import { readFileSync } from 'node:fs';

const pkg = JSON.parse(readFileSync(new URL('../package.json', import.meta.url), 'utf8'));

function readBuildTimestampUnix() {
  try {
    const meta = JSON.parse(readFileSync(new URL('../build-meta.json', import.meta.url), 'utf8'));
    return String(meta.buildTimestampUnix || '0');
  } catch {
    return '0';
  }
}

export function registerVersionCommand(program) {
  program
    .command('version')
    .description('Print the current AYNIG version and build timestamp')
    .action(() => {
      process.stdout.write(`${pkg.version} ${readBuildTimestampUnix()}\n`);
    });
}
