import { spawnSync } from 'node:child_process';

const INSTALL_COMMAND = 'curl -fsSL https://aynig.org/install.sh | bash';

export function registerUpdateCommand(program) {
  program
    .command('update')
    .description('Download and install the latest AYNIG release')
    .action(() => {
      const result = spawnSync('sh', ['-c', INSTALL_COMMAND], {
        stdio: 'inherit'
      });

      if (result.status !== 0) {
        process.exit(result.status ?? 1);
      }
    });
}
