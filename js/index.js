#!/usr/bin/env -S node --no-warnings

import { Command } from 'commander';
import { registerRunCommand } from './commands/run.js';
import { registerInitCommand } from './commands/init.js';
import { registerInstallCommand } from './commands/install.js';

const program = new Command();

// Configure CLI
program
  .name('aynig')
  .description('Git-native orchestration tool for agentic workflows')
  .version('0.0.1');

// Register commands
registerRunCommand(program);
registerInitCommand(program);
registerInstallCommand(program);

// Parse arguments
program.parse();