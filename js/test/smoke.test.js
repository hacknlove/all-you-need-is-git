import test from 'node:test';
import assert from 'node:assert/strict';
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

  assert.equal(result.status, 0);
  assert.match(result.stdout, /Usage: aynig/);
  assert.match(result.stdout, /run/);
  assert.match(result.stdout, /init/);
  assert.match(result.stdout, /install/);
  assert.equal(result.stderr, '');
});
