import test from 'node:test';
import assert from 'node:assert/strict';
import { Command } from '../AgentsOrchestrator/Command.js';

function createCommand(overrides = {}) {
  return new Command({
    config: { leaseSeconds: 123 },
    branchName: 'main',
    isCurrentBranch: true,
    trailers: { 'aynig-state': 'build', 'aynig-run-id': 'run-123' },
    body: 'do the thing',
    commitDate: '2025-01-01T00:00:00.000Z',
    ...overrides
  });
}

test('run writes working commit with configured lease seconds', async () => {
  const calls = [];
  const gitStub = {
    revparse: async () => 'abc123\n',
    commit: async (message) => {
      calls.push(message);
    },
    push: async () => {}
  };

  const command = createCommand({
    gitFactory: () => gitStub,
    spawnImpl: () => ({ unref: () => {} })
  });

  command.getWorkspace = async () => '/tmp/worktree';
  command.getCommandPath = async () => '/tmp/worktree/.aynig/command/build';

  await command.run();

  assert.equal(calls.length, 1);
  assert.match(calls[0], /aynig-lease-seconds: 123/);
});

test('run passes trailers, body, and commit hash to env', async () => {
  let spawnEnv;
  const gitStub = {
    revparse: async () => 'deadbeef\n',
    commit: async () => {},
    push: async () => {}
  };

  const command = createCommand({
    gitFactory: () => gitStub,
    spawnImpl: (_cmd, _args, options) => {
      spawnEnv = options.env;
      return { unref: () => {} };
    }
  });

  command.getWorkspace = async () => '/tmp/worktree';
  command.getCommandPath = async () => '/tmp/worktree/.aynig/command/build';

  await command.run();

  assert.equal(spawnEnv.AYNIG_BODY, 'do the thing');
  assert.equal(spawnEnv.AYNIG_COMMIT_HASH, 'deadbeef');
  assert.equal(spawnEnv.AYNIG_TRAILER_AYNIG_STATE, 'build');
  assert.equal(spawnEnv.AYNIG_TRAILER_AYNIG_RUN_ID, 'run-123');
});
