import { Command } from 'commander';
import defaultConfig from '../defaultConfig.json' with { type: 'json' };
import { Repo } from '../AgentsOrchestrator/Repo.js';

/**
 * Run command - runs AYNIG for the current repository
 *
 * It:
 * 1. instantiates Repo objects for each repository defined in the config
 * 2. runs them
 * 3. waits for all to complete
 */
async function action(options) {
  const config = {
    ...defaultConfig,
    ...options
  }

  const repository = new Repo(config);
  try {
    await repository.run();
  } catch (error) {
    console.error('Error trying to set up the Repository:', error);
  }
}

export function registerRunCommand(program) {
  program
    .command('run')
    .description('Run AYNIG for the current repository')
    .option('-w, --worktree <path>', 'Specify custom worktree directory (default: .worktrees)')
    .option('--use-remote <name>', 'Specify which remote to use (default: origin)')
    .action(action);
}
