import { test, expect } from 'vitest';
import { spawnSync } from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import { fileURLToPath } from 'node:url';
import path from 'node:path';

const testDir = path.dirname(fileURLToPath(import.meta.url));
const projectDir = path.resolve(testDir, '..');

test('CLI help renders successfully', () => {
  const result = spawnSync(process.execPath, ['--no-warnings', 'index.js', '--help'], {
    cwd: projectDir,
    encoding: 'utf8'
  });

  expect(result.status).toBe(0);
  expect(result.stdout).toMatch(/Usage: aynig/);
  expect(result.stdout).toMatch(/run/);
  expect(result.stdout).toMatch(/init/);
  expect(result.stdout).toMatch(/install/);
  expect(result.stderr).toBe('');
});

test('init creates command directory and clean command', () => {
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'aynig-init-'));
  const initRepo = spawnSync('git', ['init'], {
    cwd: tempDir,
    encoding: 'utf8'
  });
  expect(initRepo.status).toBe(0);

  const result = spawnSync(process.execPath, ['--no-warnings', path.join(projectDir, 'index.js'), 'init'], {
    cwd: tempDir,
    encoding: 'utf8'
  });

  expect(result.status).toBe(0);
  expect(result.stderr).toBe('');

  const commandDir = path.join(tempDir, '.aynig', 'command');
  const cleanCommand = path.join(commandDir, 'clean');
  const commandsDoc = path.join(tempDir, '.aynig', 'COMMANDS.md');

  expect(fs.existsSync(commandDir)).toBe(true);
  expect(fs.existsSync(cleanCommand)).toBe(true);
  expect(fs.existsSync(commandsDoc)).toBe(true);
});
