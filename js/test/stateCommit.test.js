import { test, expect } from 'vitest';
import { buildStateCommitMessage } from '../AgentsOrchestrator/stateCommit.js';

test('buildStateCommitMessage renders prompt and contiguous trailer block', () => {
  const message = buildStateCommitMessage('chore: working', 'keep lease alive', [
    { key: 'aynig-state', value: 'working' },
    { key: 'aynig-run-id', value: 'abc' }
  ]);

  expect(message).toBe([
    'chore: working',
    '',
    'keep lease alive',
    '',
    'aynig-state: working',
    'aynig-run-id: abc',
    ''
  ].join('\n'));
});
