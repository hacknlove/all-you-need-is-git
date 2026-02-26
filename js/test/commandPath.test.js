import { beforeEach, test, expect, vi } from 'vitest';
import fs from 'node:fs/promises';
import os from 'node:os';
import path from 'node:path';

const commitMock = vi.fn(async () => {});
const pushMock = vi.fn(async () => {});
const rawMock = vi.fn(async () => {
  const lastCall = commitMock.mock.calls[commitMock.mock.calls.length - 1];
  return lastCall ? lastCall[0] : '';
});

vi.mock('simple-git', () => ({
  default: () => ({
    commit: commitMock,
    push: pushMock,
    raw: rawMock
  })
}));

const { Command } = await import('../AgentsOrchestrator/Command.js');

beforeEach(() => {
  commitMock.mockClear();
  pushMock.mockClear();
  rawMock.mockClear();
});

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

test('getCommandPath resolves nested commands under .aynig/command', async () => {
  const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), 'aynig-'));
  try {
    const commandPath = path.join(tempDir, '.aynig', 'command', 'foo', 'bar', 'baz');
    await createExecutable(commandPath);

    const command = createCommand('foo/bar/baz');
    const resolved = await command.getCommandPath(tempDir);
    expect(resolved).toBe(commandPath);
  } finally {
    await fs.rm(tempDir, { recursive: true, force: true });
  }
});

test('getCommandPath rejects traversal outside .aynig/command', async () => {
  const tempDir = await fs.mkdtemp(path.join(os.tmpdir(), 'aynig-'));
  try {
    const cases = [
      '../evil',
      'foo/../../evil',
      '/tmp/evil',
      'foo\\bar'
    ];

    for (const state of cases) {
      const command = createCommand(state);
      const resolved = await command.getCommandPath(tempDir);
      expect(resolved).toBe(false);
    }
  } finally {
    await fs.rm(tempDir, { recursive: true, force: true });
  }
});

test('checkWorking creates a stalled commit when lease expires', async () => {
  const command = new Command({
    config: {},
    branchName: 'main',
    isCurrentBranch: true,
    trailers: {
      'aynig-state': 'working',
      'aynig-lease-seconds': '1',
      'aynig-run-id': 'run-123'
    },
    body: '',
    commitDate: new Date(Date.now() - 2000).toISOString()
  });

  await command.checkWorking();

  expect(commitMock).toHaveBeenCalledTimes(1);
  const [message] = commitMock.mock.calls[0];
  expect(message).toMatch(/aynig-state: stalled/);
  expect(message).toMatch(/aynig-stalled-run: run-123/);
  expect(pushMock).not.toHaveBeenCalled();
});

test('checkWorking does nothing when lease is still valid', async () => {
  const command = new Command({
    config: {},
    branchName: 'main',
    isCurrentBranch: true,
    trailers: {
      'aynig-state': 'working',
      'aynig-lease-seconds': '300',
      'aynig-run-id': 'run-123'
    },
    body: '',
    commitDate: new Date().toISOString()
  });

  await command.checkWorking();

  expect(commitMock).not.toHaveBeenCalled();
  expect(pushMock).not.toHaveBeenCalled();
});

test('checkWorking ignores invalid lease values', async () => {
  const command = new Command({
    config: {},
    branchName: 'main',
    isCurrentBranch: true,
    trailers: {
      'aynig-state': 'working',
      'aynig-lease-seconds': 'not-a-number',
      'aynig-run-id': 'run-123'
    },
    body: '',
    commitDate: new Date(Date.now() - 5000).toISOString()
  });

  await command.checkWorking();

  expect(commitMock).not.toHaveBeenCalled();
  expect(pushMock).not.toHaveBeenCalled();
});
