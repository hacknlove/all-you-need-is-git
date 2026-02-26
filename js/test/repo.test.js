import { test, expect, vi } from 'vitest';
import { Repo } from '../AgentsOrchestrator/Repo.js';
import { git } from '../gitHelpers/git.js';

test('filterBranches respects only mode with empty current', () => {
  const repo = new Repo({});
  const result = repo.filterBranches(['origin/main', 'origin/dev'], '', 'only');
  expect(result).toEqual([]);
});

test('filterBranches respects skip mode in remote sets', () => {
  const repo = new Repo({});
  const result = repo.filterBranches(['origin/main', 'origin/dev'], 'origin/main', 'skip');
  expect(result).toEqual(['origin/dev']);
});

test('resolveCurrentRemoteBranch ignores upstream from other remotes', async () => {
  const warn = vi.fn();
  const repo = new Repo({ aynigRemote: 'origin', logger: { warn } });
  const rawSpy = vi.spyOn(git, 'raw').mockResolvedValue('upstream/main\n');

  try {
    const result = await repo.resolveCurrentRemoteBranch('main');
    expect(result).toBe('');
    expect(warn).toHaveBeenCalled();
  } finally {
    rawSpy.mockRestore();
  }
});

test('remoteFromTrailers reads aynig-remote trailer', () => {
  const repo = new Repo({});
  expect(repo.remoteFromTrailers({ 'aynig-remote': 'origin' })).toBe('origin');
  expect(repo.remoteFromTrailers({ 'aynig-remote': [' upstream '] })).toBe('upstream');
  expect(repo.remoteFromTrailers({})).toBe('');
});
