import { test, expect } from 'vitest';
import { buildStateCommitMessage } from '../AgentsOrchestrator/stateCommit.js';

test('buildStateCommitMessage renders prompt and contiguous trailer block', () => {
  const message = buildStateCommitMessage('chore: working', 'keep lease alive', [
    { key: 'dwp-state', value: 'working' },
    { key: 'dwp-run-id', value: 'abc' }
  ]);

  expect(message).toBe([
    'chore: working',
    '',
    'keep lease alive',
    '',
    'dwp-state: working',
    'dwp-run-id: abc',
    ''
  ].join('\n'));
});
