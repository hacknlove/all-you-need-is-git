#!/usr/bin/env -S node --no-warnings

import { Command } from 'commander';
import { registerIterateCommand } from './commands/iterate.js';
import { registerInitCommand } from './commands/init.js';
import { registerInstallCommand } from './commands/install.js';

const program = new Command();

// Configure CLI
program
  .name('aynig')
  .description('Git-native orchestration tool for agentic workflows')
  .version('0.0.1');

// Register commands
registerIterateCommand(program);
registerInitCommand(program);
registerInstallCommand(program);

// Parse arguments
program.parse();