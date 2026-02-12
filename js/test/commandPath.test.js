import test from 'node:test';
import assert from 'node:assert/strict';
import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';
import { Command } from '../AgentsOrchestrator/Command.js';

async function createExecutable(filePath) {
  await fs.mkdir(path.dirname(filePath), { recursive: true });
  await fs.writeFile(filePath, '#!/bin/sh\nexit 0\n');
  await fs.chmod(filePath, 0o755);
}

function createCommand(state) {
  return new Command({
    config: {},
    branchName: 'main',
    isCurrentBranch: true,
    trailers: { 'aynig-state': state },
    body: ''
  });
}

test('getCommandPath resolves nested commands under .aynig/command', async (t) => {
  const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), 'aynig-'));
  t.after(async () => {
    await fs.rm(tempDir, { recursive: true, force: true });
  });

  const commandPath = path.join(tempDir, '.aynig', 'command', 'foo', 'bar', 'baz');
  await createExecutable(commandPath);

  const command = createCommand('foo/bar/baz');
  const resolved = await command.getCommandPath(tempDir);
  assert.equal(resolved, commandPath);
});

test('getCommandPath rejects traversal outside .aynig/command', async (t) => {
  const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), 'aynig-'));
  t.after(async () => {
    await fs.rm(tempDir, { recursive: true, force: true });
  });

  const cases = [
    '../evil',
    'foo/../../evil',
    '/tmp/evil',
    'foo\\bar'
  ];

  for (const state of cases) {
    const command = createCommand(state);
    const resolved = await command.getCommandPath(tempDir);
    assert.equal(resolved, false, `expected false for ${state}`);
  }
});
