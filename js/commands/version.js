import { readFileSync } from 'node:fs';

const pkg = JSON.parse(readFileSync(new URL('../package.json', import.meta.url), 'utf8'));

export function registerVersionCommand(program) {
  program
    .command('version')
    .description('Print the current AYNIG version')
    .action(() => {
      process.stdout.write(`${pkg.version}\n`);
    });
}
