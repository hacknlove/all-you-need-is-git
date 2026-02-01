#!/usr/bin/env node

import { Command } from 'commander';
import defaultConfig from './defaultConfig.json' with { type: 'json' };
import { Repo } from './AgentsOrchestrator/Repo.js';

const program = new Command();

/**
 * Iterate command - runs AYNIG iteration for the current repository
 *
 * It:
 * 1. instantiates Repo objects for each repository defined in the config
 * 2. runs them
 * 3. waits for all to complete
 */
async function iterate(options) {
  const config = {
    ...defaultConfig,
    ...options
  }

  const repository = new Repo(repoConfig);
  try {
    await repository.run();
  } catch (error) {
    console.error('Error trying to set up the Repository:', error);
  }
}

/**
 * Init command - initializes AYNIG in the current repository
 */
async function init() {
  console.log('Hello from init command!');
  // TODO: Implement init logic
  // - create .aynig/ directory
  // - create .worktrees/ directory
  // - add .worktrees/ to .gitignore
}

/**
 * Install command - installs AYNIG workflows from another repository
 */
async function install(repo, ref, subfolder) {
  console.log('Hello from install command!');
  console.log(`Repo: ${repo}`);
  console.log(`Ref: ${ref}`);
  console.log(`Subfolder: ${subfolder || '(none)'}`);
  // TODO: Implement install logic
  // - clone/fetch the specified repository
  // - copy .aynig/ scripts from source to current repo
}

// Configure CLI
program
  .name('aynig')
  .description('Git-native orchestration tool for agentic workflows')
  .version('0.0.1');

// Iterate command
program
  .command('iterate')
  .description('Run AYNIG iteration for the current repository')
  .option('-w, --worktree <path>', 'Specify custom worktree directory (default: .worktrees)')
  .option('--use-remote <name>', 'Specify which remote to use (default: origin)')
  .action(iterate);

// Init command
program
  .command('init')
  .description('Initialize AYNIG in the current repository')
  .action(init);

// Install command
program
  .command('install <repo> <ref> [subfolder]')
  .description('Install AYNIG workflows from another repository')
  .action(install);

// Parse arguments
program.parse();