#!/usr/bin/env -S node --no-warnings

import { Command } from 'commander';
import { registerRunCommand } from './commands/run.js';
import { registerInitCommand } from './commands/init.js';
import { registerInstallCommand } from './commands/install.js';
import { registerUpdateCommand } from './commands/update.js';
import { registerVersionCommand } from './commands/version.js';
import pkg from './package.json' assert { type: 'json' };

const program = new Command();

// Configure CLI
program
  .name('aynig')
  .description('Git-native orchestration tool for agentic workflows')
  .version(pkg.version);

// Register commands
registerRunCommand(program);
registerInitCommand(program);
registerInstallCommand(program);
registerUpdateCommand(program);
registerVersionCommand(program);

// Parse arguments
program.parse();
