import test from 'node:test';
import assert from 'node:assert/strict';
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

test('checkWorking skips when lease not expired', async (t) => {
  const commitCalls = [];
  const gitStub = {
    commit: async (message) => commitCalls.push(message),
    push: async () => {}
  };

  const command = createCommand({
    gitFactory: () => gitStub
  });

  command.getWorkspace = async () => '/tmp/worktree';

  const originalNow = Date.now;
  Date.now = () => new Date('2025-01-01T00:00:05.000Z').getTime();
  t.after(() => {
    Date.now = originalNow;
  });

  await command.checkWorking();

  assert.equal(commitCalls.length, 0);
});

test('checkWorking emits stalled commit when lease expired', async (t) => {
  const commitCalls = [];
  const gitStub = {
    commit: async (message) => commitCalls.push(message),
    push: async () => {}
  };

  const command = createCommand({
    gitFactory: () => gitStub
  });

  command.getWorkspace = async () => '/tmp/worktree';

  const originalNow = Date.now;
  Date.now = () => new Date('2025-01-01T00:00:20.000Z').getTime();
  t.after(() => {
    Date.now = originalNow;
  });

  await command.checkWorking();

  assert.equal(commitCalls.length, 1);
  assert.match(commitCalls[0], /aynig-state: stalled/);
  assert.match(commitCalls[0], /aynig-stalled-run: run-789/);
});
