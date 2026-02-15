import { test, expect } from 'vitest';
import { spawnSync } from 'node:child_process';
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
