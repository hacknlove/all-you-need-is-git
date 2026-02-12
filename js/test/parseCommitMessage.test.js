import test from 'node:test';
import assert from 'node:assert/strict';
import { parseCommitMessage } from '../gitHelpers/parseCommitMessage.js';

test('returns empty values when commit log is missing', async () => {
  assert.deepEqual(await parseCommitMessage(), {
    firstLine: '',
    body: '',
    trailers: {}
  });
});

test('parses aynig trailers from trailing trailer block', async () => {
  const parsed = await parseCommitMessage({
    message: 'feat: schedule next run',
    body: [
      'Plan next execution step for this branch.',
      '',
      'aynig-state: queued',
      'aynig-run-id: abc-123',
      'Signed-off-by: Example User <example@example.com>'
    ].join('\n')
  });

  assert.equal(parsed.firstLine, 'feat: schedule next run');
  assert.deepEqual({...parsed.trailers}, {
    'aynig-state': 'queued',
    'aynig-run-id': 'abc-123',
    'Signed-off-by': 'Example User <example@example.com>'
  });
});

test('does not parse trailer-like lines outside trailing trailer block', async () => {
  const parsed = await parseCommitMessage({
    message: 'docs: explain state',
    body: [
      'aynig-state: queued',
      'This line is part of prompt and not a trailer block.'
    ].join('\n')
  });

  assert.equal(parsed.body, 'aynig-state: queued\nThis line is part of prompt and not a trailer block.');
  assert.deepEqual(parsed.trailers, {});
});
