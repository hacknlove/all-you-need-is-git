import { test, expect, vi } from 'vitest';
import { Command } from '../AgentsOrchestrator/Command.js';

function createCommand(overrides = {}) {
  return new Command({
    config: {},
    branchName: 'main',
    isCurrentBranch: true,
    trailers: {
      'aynig-state': 'working',
      'aynig-run-id': 'run-789',
      'aynig-lease-seconds': '10'
    },
    body: '',
    commitDate: '2025-01-01T00:00:00.000Z',
    ...overrides
  });
}

test('checkWorking skips when lease not expired', async () => {
  const commitCalls = [];
  const gitStub = {
    commit: async (message) => commitCalls.push(message),
    push: async () => {},
    raw: async () => commitCalls[commitCalls.length - 1] || ''
  };

  const command = createCommand({
    gitFactory: () => gitStub
  });

  command.getWorkspace = async () => '/tmp/worktree';

  const originalNow = Date.now;
  Date.now = () => new Date('2025-01-01T00:00:05.000Z').getTime();
  try {
    await command.checkWorking();
  } finally {
    Date.now = originalNow;
  }
  expect(commitCalls.length).toBe(0);
});

test('checkWorking emits stalled commit when lease expired', async () => {
  const commitCalls = [];
  const gitStub = {
    commit: async (message) => commitCalls.push(message),
    push: async () => {},
    raw: async () => commitCalls[commitCalls.length - 1] || ''
  };

  const command = createCommand({
    gitFactory: () => gitStub
  });

  command.getWorkspace = async () => '/tmp/worktree';

  const originalNow = Date.now;
  Date.now = () => new Date('2025-01-01T00:00:20.000Z').getTime();
  try {
    await command.checkWorking();
  } finally {
    Date.now = originalNow;
  }

  expect(commitCalls.length).toBe(1);
  expect(commitCalls[0]).toMatch(/aynig-state: stalled/);
  expect(commitCalls[0]).toMatch(/aynig-stalled-run: run-789/);
});

test('checkWorking skips when commit date is invalid', async () => {
  const gitStub = {
    commit: vi.fn(async () => {}),
    push: vi.fn(async () => {}),
    raw: vi.fn(async () => '')
  };

  const command = createCommand({
    commitDate: 'not-a-date',
    gitFactory: () => gitStub
  });

  command.getWorkspace = async () => '/tmp/worktree';

  await command.checkWorking();

  expect(gitStub.commit).not.toHaveBeenCalled();
  expect(gitStub.push).not.toHaveBeenCalled();
});
