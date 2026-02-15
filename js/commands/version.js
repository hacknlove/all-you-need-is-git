import pkg from '../package.json' assert { type: 'json' };

export function registerVersionCommand(program) {
  program
    .command('version')
    .description('Print the current AYNIG version')
    .action(() => {
      process.stdout.write(`${pkg.version}\n`);
    });
}
