import { test, expect, vi } from 'vitest';
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

  expect(calls.length).toBe(1);
  expect(calls[0]).toMatch(/aynig-lease-seconds: 123/);
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

  expect(spawnEnv.AYNIG_BODY).toBe('do the thing');
  expect(spawnEnv.AYNIG_COMMIT_HASH).toBe('deadbeef');
  expect(spawnEnv['AYNIG_TRAILER_AYNIG_STATE']).toBe('build');
  expect(spawnEnv['AYNIG_TRAILER_AYNIG_RUN_ID']).toBe('run-123');
});

test('run skips when command path is missing', async () => {
  const gitStub = {
    revparse: vi.fn(async () => 'deadbeef\n'),
    commit: vi.fn(async () => {}),
    push: vi.fn(async () => {})
  };

  const command = createCommand({
    gitFactory: () => gitStub,
    spawnImpl: () => ({ unref: () => {} })
  });

  command.getWorkspace = async () => '/tmp/worktree';
  command.getCommandPath = async () => false;

  await command.run();

  expect(gitStub.revparse).not.toHaveBeenCalled();
  expect(gitStub.commit).not.toHaveBeenCalled();
  expect(gitStub.push).not.toHaveBeenCalled();
});
